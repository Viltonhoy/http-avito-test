package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/zapadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type Storage struct {
	logger *zap.Logger
	db     *pgxpool.Pool
}

const cacheBookAccountID = int64(0)

func NewStore(ctx context.Context, logger *zap.Logger) (*Storage, error) {
	if logger == nil {
		return nil, errors.New("no logger provided")
	}

	//строка подключения
	config, _ := pgxpool.ParseConfig("")

	config.ConnConfig.Logger = zapadapter.NewLogger(logger)
	config.ConnConfig.LogLevel = pgx.LogLevelError

	// 	//создать пул соединений
	pool, err := pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		logger.Error("cannot connect using config", zap.Error(err))
		return nil, err
	}

	err = pool.Ping(ctx)
	if err != nil {
		logger.Error("connection is lost", zap.Error(err))
		return &Storage{}, err
	}

	return &Storage{
		logger: logger,
		db:     pool,
	}, err
}

func (s *Storage) Close() {
	s.logger.Info("closing Storage connection")
	s.db.Close()
}

func (s *Storage) ReadUser(ctx context.Context, user_id int64) (u UserBalance, err error) {
	logger := s.logger
	logger.Info("reading the user balance", zap.Int64("userID", user_id))

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return UserBalance{}, err
	}

	refreshSql := `REFRESH MATERIALIZED VIEW account_balances;`

	_, err = tx.Exec(ctx, refreshSql)
	if err != nil {
		tx.Rollback(ctx)
		logger.Error("failed to refresh materialized view account_balances", zap.Error(err))
		return UserBalance{}, err
	}

	selectSql := `SELECT balance FROM account_balances WHERE user_id = $1;`

	//выполнить запрос
	err = tx.QueryRow(ctx, selectSql, user_id).Scan(&u.Balance)
	if err != nil {
		tx.Rollback(ctx)
		logger.Error("cannot return user with specified ID")
		if errors.Is(err, pgx.ErrNoRows) {
			return u, errors.New("user does not exist")
		}
		return UserBalance{}, err
	}

	u.AccountID = user_id

	err = tx.Commit(ctx)
	return UserBalance{
		u.AccountID,
		u.Balance,
	}, err
}

func (s *Storage) Deposit(ctx context.Context, user_id int64, amount decimal.Decimal) (err error) {
	logger := s.logger
	logger.Info("money deposit", zap.Int64(`userID`, user_id), zap.String(`amount`, amount.String()))

	var now = time.Now()

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}

	firstInsertSql := `INSERT INTO posting (account_id, cb_journal, accounting_period, amount, date)
			VALUES ($1, $4, $5, $2, $3);`

	_, err = tx.Exec(
		ctx,
		firstInsertSql,
		user_id,
		amount,
		now.Format("2006-01-02 15:04:05"),
		operationTypeDeposit,
		fmt.Sprintf(`Period: %d`, now.Year()),
	)
	if err != nil {
		tx.Rollback(ctx)
		logger.Error("failed to insert record", zap.Error(err))
		return err
	}

	secondInsertSql := `INSERT INTO posting (account_id, cb_journal, accounting_period, amount, date)
			VALUES ($5, $3, $4, -1 * $1, $2);`

	_, err = tx.Exec(
		ctx,
		secondInsertSql,
		amount,
		now.Format("2006-01-02 15:04:05"),
		operationTypeDeposit,
		fmt.Sprintf(`Period: %d`, now.Year()),
		cacheBookAccountID,
	)
	if err != nil {
		tx.Rollback(ctx)
		logger.Error("failed to insert record: %v", zap.Error(err))
		return err
	}
	err = tx.Commit(ctx)
	return err
}

func (s *Storage) Withdrawal(ctx context.Context, user_id int64, amount decimal.Decimal, description string) (err error) {
	logger := s.logger
	logger.Info("money withdrawal", zap.Int64(`userID`, user_id), zap.String(`amount`, amount.String()))

	var now = time.Now()

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}

	refreshSql := `REFRESH MATERIALIZED VIEW account_balances;`

	_, err = tx.Exec(ctx, refreshSql)
	if err != nil {
		tx.Rollback(ctx)
		logger.Error("failed to refresh materialized view account_balances", zap.Error(err))
		return err
	}

	selectSql := `SELECT balance FROM account_balances WHERE user_id = $1;`

	var balance UserBalance
	err = tx.QueryRow(ctx, selectSql, user_id).Scan(&balance.Balance)
	if err != nil {
		tx.Rollback(ctx)
		logger.Error("cannot return balance with specified ID")
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("user does not exist")
		}
		return err
	}

	if amount.GreaterThan(balance.Balance) {
		tx.Rollback(ctx)
		logger.Error("insufficient funds on the user's account")
		return err
	}

	firstInsertSql := `INSERT INTO posting (account_id, cb_journal, accounting_period, amount, date, description)
			VALUES ($1, $4, $5, -1 * $2, $3, $6);`

	_, err = tx.Exec(
		ctx,
		firstInsertSql,
		user_id,
		amount,
		now.Format("2006-01-02 15:04:05"),
		operationTypeWithdrawal,
		fmt.Sprintf(`Period: %d`, now.Year()),
		description,
	)
	if err != nil {
		tx.Rollback(ctx)
		logger.Error("failed to insert record: %v", zap.Error(err))
		return err
	}

	secondInsertSql := `INSERT INTO posting (account_id, cb_journal, accounting_period, amount, date)
			VALUES ($5, $3, $4, $1, $2);`

	_, err = tx.Exec(
		ctx,
		secondInsertSql,
		amount,
		now.Format("2006-01-02 15:04:05"),
		operationTypeWithdrawal,
		fmt.Sprintf(`Period: %d`, now.Year()),
		cacheBookAccountID,
	)
	if err != nil {
		tx.Rollback(ctx)
		logger.Error("failed to insert record: %v", zap.Error(err))
		return err
	}
	err = tx.Commit(ctx)
	return err
}

func (s *Storage) Transfer(ctx context.Context, user_id1, user_id2 int64, amount decimal.Decimal, description string) error {
	logger := s.logger
	logger.Info("money transfer", zap.Int64("senderID", user_id1), zap.Int64("recipientID", user_id2), zap.String("amount", amount.String()))

	var now = time.Now()

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}

	refreshSql := `REFRESH MATERIALIZED VIEW account_balances;`

	_, err = tx.Exec(ctx, refreshSql)
	if err != nil {
		tx.Rollback(ctx)
		logger.Error("failed to refresh materialized view account_balances", zap.Error(err))
		return err
	}

	selectSql := `SELECT balance FROM account_balances WHERE user_id = $1;`

	var balance UserBalance
	err = tx.QueryRow(ctx, selectSql, user_id1).Scan(&balance.Balance)
	if err != nil {
		tx.Rollback(ctx)
		logger.Error("cannot return balance with specified ID")
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("user does not exist")
		}
		return err
	}

	if amount.GreaterThan(balance.Balance) {
		tx.Rollback(ctx)
		logger.Error("insufficient funds on the user's account")
		return err
	}

	firstInsertSql := `INSERT INTO posting (account_id, cb_journal, accounting_period, amount, date, addressee, description) 
			VALUES ($1, $6, $4, -1 * $2, $3, $5, $7);`

	_, err = tx.Exec(
		ctx,
		firstInsertSql,
		user_id1,
		amount,
		now.Format("2006-01-02 15:04:05"),
		fmt.Sprintf(`Period: %d`, now.Year()),
		user_id2,
		operationTypeTransfer,
		description,
	)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	secondInsertSql := `INSERT INTO posting (account_id, cb_journal, accounting_period, amount, date, addressee) 
			VALUES ($1, $6, $4, $2, $3, $5);`

	_, err = tx.Exec(
		ctx,
		secondInsertSql,
		user_id2,
		amount,
		now.Format("2006-01-02 15:04:05"),
		fmt.Sprintf(`Period: %d`, now.Year()),
		user_id1,
		operationTypeTransfer,
	)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	err = tx.Commit(ctx)
	return err
}

func (s *Storage) ReadUserHistoryList(ctx context.Context, user_id int64, order string, limit, offset int64) ([]Transfer, error) {
	logger := s.logger
	logger.Info("reading the user history list", zap.Int64("userID", user_id), zap.String("order", order), zap.Int64("limit", limit), zap.Int64("offset", offset))

	var sql string

	amountSql := `SELECT account_id, cb_journal, amount, date, addressee, description FROM posting 
		WHERE account_id = $1 ORDER BY amount LIMIT $2 OFFSET $3;`

	dateSql := `SELECT account_id, cb_journal, amount, date, addressee, description FROM posting 
	WHERE account_id = $1 ORDER BY date LIMIT $2 OFFSET $3;`

	switch order {
	case "amount":
		sql = amountSql
	case "date":
		sql = dateSql
	}

	rows, err := s.db.Query(
		ctx,
		sql,
		user_id,
		limit,
		offset,
	)

	if err != nil {
		logger.Error("cannot return user list with specified ID")
		return []Transfer{}, err
	}

	var list []Transfer
	for rows.Next() {
		var tt Transfer
		err := rows.Scan(&tt.AcountID, &tt.CBjournal, &tt.Amount, &tt.Date, &tt.Addressee, &tt.Description)
		if err != nil {
			s.logger.Error("scanning row", zap.Error(err))
		}
		tt.Amount = decimal.New(tt.Amount.IntPart(), -2)
		list = append(list, tt)
	}
	return list, nil
}
