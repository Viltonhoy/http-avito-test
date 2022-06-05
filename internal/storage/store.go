package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/zapadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
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

	conf := NewDbSreverConfig()
	//строка подключения
	var connString = fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", conf.User, conf.Password, conf.Database)
	config, _ := pgxpool.ParseConfig(connString)

	config.ConnConfig.Logger = zapadapter.NewLogger(logger)
	config.ConnConfig.LogLevel = pgx.LogLevelError

	// 	//создать пул соединений
	pool, err := pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("cannot connect using config %+v: %w", config, err)
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
	logger.With(zap.Int64(`reading the balance of %d user`, user_id))

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
	logger.Info(``, zap.Int64(`Updating users account information ID: %d`, user_id), zap.String(`Balance: %s`, amount.String()))

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

func (s *Storage) Withdrawal(ctx context.Context, user_id int64, amount decimal.Decimal) (err error) {
	logger := s.logger
	logger.With(zap.Int64(`Updating users account information: %d`, user_id), zap.String(`Balance: %s`, amount.String()))

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

	firstInsertSql := `INSERT INTO posting (account_id, cb_journal, accounting_period, amount, date)
			VALUES ($1, $4, $5, -1 * $2, $3);`

	_, err = tx.Exec(
		ctx,
		firstInsertSql,
		user_id,
		amount,
		now.Format("2006-01-02 15:04:05"),
		operationTypeWithdrawal,
		fmt.Sprintf(`Period: %d`, now.Year()),
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

func (s *Storage) Transfer(ctx context.Context, user_id1, user_id2 int64, amount decimal.Decimal) error {
	logger := s.logger
	logger.Sugar().Debugf(`money transfer from %d user to %d`, user_id1, user_id2)

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

	firstInsertSql := `INSERT INTO posting (account_id, cb_journal, accounting_period, amount, date, addressee) 
			VALUES ($1, $6, $4, -1 * $2, $3, $5);`

	_, err = tx.Exec(
		ctx,
		firstInsertSql,
		user_id1,
		amount,
		now.Format("2006-01-02 15:04:05"),
		fmt.Sprintf(`Period: %d`, now.Year()),
		user_id2,
		operationTypeTransfer,
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
	logger.Info("reading a history list")

	var sql string

	amountSql := `SELECT account_id, cb_journal, amount, date, addressee FROM posting 
		WHERE account_id = $1 ORDER BY amount LIMIT $2 OFFSET $3;`

	dateSql := `SELECT account_id, cb_journal, amount, date, addressee FROM posting 
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
		err := rows.Scan(&tt.AcountID, &tt.CBjournal, &tt.Amount, &tt.Date, &tt.Addressee)
		if err != nil {
			s.logger.Error("scanning row", zap.Error(err))
		}
		tt.Amount = decimal.New(tt.Amount.IntPart(), -2)
		list = append(list, tt)
	}
	return list, nil
}
