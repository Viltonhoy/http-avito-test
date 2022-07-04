package storage

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/zapadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

// Storage defines fields used in interaction processes of database
type Storage struct {
	Logger *zap.Logger
	DB     *pgxpool.Pool
}

const cacheBookAccountID = int64(0)

const updateRollUpTable = `
	with var1 as (
	select id from posting where account_id = $1 order by id desc limit 1
	), var2 as(
	select coalesce(sum(amount),0) from posting where account_id = $1 and id > (select coalesce((select last_tx_id from balances where account_id = $1),0))
	) insert into balances (
	balance,
	account_id,
	last_tx_id
	) values (
	(select * from var2),
	$1,
	(select * from var1)
	) on conflict (account_id) do update
	set last_tx_id = (select * from var1),
	balance = (select * from var2) + (select balance from balances where account_id = $1) returning balance`

var ErrUserAvailability = errors.New("sender does not exist")

var (
	ErrWithdrawal = errors.New("not enough money to withdraw")
	ErrTransfer   = errors.New("not enough money to transfer")
)

// NewStore constructs Store instance with configured logger
func NewStorage(ctx context.Context, logger *zap.Logger) (*Storage, error) {
	if logger == nil {
		return nil, errors.New("no logger provided")
	}

	// taking connect info from environment variables
	config, _ := pgxpool.ParseConfig("")

	config.ConnConfig.Logger = zapadapter.NewLogger(logger)
	config.ConnConfig.LogLevel = pgx.LogLevelError

	// create a pool connection
	pool, err := pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		logger.Error("error database connection", zap.Error(err))
		return nil, err
	}

	err = pool.Ping(ctx)
	if err != nil {
		logger.Error("connection was not established", zap.Error(err))
		return &Storage{}, err
	}

	return &Storage{
		Logger: logger,
		DB:     pool,
	}, err
}

// Close closes all database connections in pool
func (s *Storage) Close() {
	s.Logger.Info("closing Storage connection")
	s.DB.Close()
}

//ReadUser reads user's balance and returns it's id and balance
func (s *Storage) ReadUserByID(ctx context.Context, userID int64) (u User, err error) {
	logger := s.Logger.With(zap.Int64("user_ID", userID))
	logger.Debug("reading the user balance")

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return User{}, err
	}

	defer func() {
		if err != nil {
			if errRollback := tx.Rollback(ctx); errRollback != nil {
				logger.Error("error rolls back the transaction", zap.Error(err))
			}
		}
	}()

	//query execution
	err = tx.QueryRow(ctx, updateRollUpTable, userID).Scan(&u.Balance)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.NotNullViolation {
			logger.Error("error returning user balance with specified id: user does not exist", zap.Error(err))
			return User{}, ErrUserAvailability
		}
		logger.Error("error returning user balance with specified id", zap.Error(err))
		return User{}, err
	}

	u.AccountID = userID

	err = tx.Commit(ctx)
	return User{
		u.AccountID,
		u.Balance,
	}, err
}

func (s *Storage) Deposit(ctx context.Context, userID int64, amount decimal.Decimal) (err error) {
	logger := s.Logger.With(zap.Int64(`user_ID`, userID))
	logger.Debug("money deposit")

	var now = time.Now()

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			if errRollback := tx.Rollback(ctx); errRollback != nil {
				logger.Error("error rolls back the transaction", zap.Error(err))
			}
		}
	}()

	firstInsertExec := `INSERT INTO posting (account_id, cb_journal, accounting_period, amount, date)
			VALUES ($1, $4, $5, $2, $3);`

	_, err = tx.Exec(
		ctx,
		firstInsertExec,
		userID,
		amount,
		now.Format(time.RFC3339),
		OperationTypeDeposit,
		now,
	)
	if err != nil {
		logger.Error("failed to insert record", zap.Error(err))
		return err
	}

	secondInsertExec := `INSERT INTO posting (account_id, cb_journal, accounting_period, amount, date)
			VALUES ($5, $3, $4, -1 * $1, $2);`

	_, err = tx.Exec(
		ctx,
		secondInsertExec,
		amount,
		now.Format(time.RFC3339),
		OperationTypeDeposit,
		now,
		cacheBookAccountID,
	)
	if err != nil {
		logger.Error("failed to insert record", zap.Error(err))
		return err
	}
	err = tx.Commit(ctx)
	return err
}

func (s *Storage) Withdrawal(ctx context.Context, userID int64, amount decimal.Decimal, description *string) (err error) {
	logger := s.Logger.With(zap.Int64("userID", userID))
	logger.Debug("money withdrawal")

	var now = time.Now()

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			if errRollback := tx.Rollback(ctx); errRollback != nil {
				logger.Error("error rolls back the transaction", zap.Error(err))
			}
		}
	}()

	var balance User
	err = tx.QueryRow(ctx, updateRollUpTable, userID).Scan(&balance.Balance)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.NotNullViolation {
			logger.Error("error returning user balance with specified id: user does not exist", zap.Error(err))
			return ErrUserAvailability
		}
		logger.Error("error returning user balance with specified id", zap.Error(err))
		return err
	}

	if amount.GreaterThan(balance.Balance) {
		logger.Error("insufficient funds on the user's account", zap.Error(ErrWithdrawal))
		return ErrWithdrawal
	}

	firstInsertExec := `INSERT INTO posting (account_id, cb_journal, accounting_period, amount, date, description)
			VALUES ($1, $4, $5, -1 * $2, $3, $6);`

	_, err = tx.Exec(
		ctx,
		firstInsertExec,
		userID,
		amount,
		now.Format(time.RFC3339),
		OperationTypeWithdrawal,
		now,
		description,
	)
	if err != nil {
		logger.Error("failed to insert record", zap.Error(err))
		return err
	}

	secondInsertExec := `INSERT INTO posting (account_id, cb_journal, accounting_period, amount, date)
			VALUES ($5, $3, $4, $1, $2);`

	_, err = tx.Exec(
		ctx,
		secondInsertExec,
		amount,
		now.Format(time.RFC3339),
		OperationTypeWithdrawal,
		now,
		cacheBookAccountID,
	)
	if err != nil {
		logger.Error("failed to insert record", zap.Error(err))
		return err
	}
	err = tx.Commit(ctx)
	return err
}

func (s *Storage) Transfer(ctx context.Context, sender, recipient int64, amount decimal.Decimal, description *string) error {
	logger := s.Logger.With(zap.Int64("senderID", sender), zap.Int64("recipientID", recipient))
	logger.Debug("money transfer")

	var now = time.Now()

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			if errRollback := tx.Rollback(ctx); errRollback != nil {
				logger.Error("error rolls back the transaction")
			}
		}
	}()

	var balance User
	err = tx.QueryRow(ctx, updateRollUpTable, sender).Scan(&balance.Balance)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.NotNullViolation {
			logger.Error("error returning user balance with specified id: user does not exist", zap.Error(err))
			return ErrUserAvailability
		}
		logger.Error("error returning user balance with specified id", zap.Error(err))
		return err
	}

	if amount.GreaterThan(balance.Balance) {
		logger.Error("insufficient funds on the sender's account", zap.Error(ErrTransfer))
		return ErrTransfer
	}

	firstInsertExec := `INSERT INTO posting (account_id, cb_journal, accounting_period, amount, date, addressee, description) 
			VALUES ($1, $6, $4, -1 * $2, $3, $5, $7);`

	_, err = tx.Exec(
		ctx,
		firstInsertExec,
		sender,
		amount,
		now.Format(time.RFC3339),
		now,
		recipient,
		OperationTypeTransfer,
		description,
	)
	if err != nil {
		logger.Error("failed to insert record", zap.Error(err))
		return err
	}

	secondInsertExec := `INSERT INTO posting (account_id, cb_journal, accounting_period, amount, date, addressee) 
			VALUES ($1, $6, $4, $2, $3, $5);`

	_, err = tx.Exec(
		ctx,
		secondInsertExec,
		recipient,
		amount,
		now.Format(time.RFC3339),
		now,
		sender,
		OperationTypeTransfer,
	)
	if err != nil {
		logger.Error("failed to insert record", zap.Error(err))
		return err
	}

	err = tx.Commit(ctx)
	return err
}

func (s *Storage) ReadUserHistoryList(
	ctx context.Context,
	userID int64,
	order OrdBy,
	limit, offset int64) ([]ReadUserHistoryResult, error) {
	logger := s.Logger.With(zap.Int64("user_ID", userID))
	logger.Debug("reading the user history list", zap.String("order", string(order)), zap.Int64("limit", limit), zap.Int64("offset", offset))

	var sql string

	amountQuery := `SELECT account_id, cb_journal, amount, date, addressee, description FROM posting 
		WHERE account_id = $1 ORDER BY amount LIMIT $2 OFFSET $3;`

	dateQuery := `SELECT account_id, cb_journal, amount, date, addressee, description FROM posting 
		WHERE account_id = $1 ORDER BY date LIMIT $2 OFFSET $3;`

	switch order {
	case OrderByAmount:
		sql = amountQuery
	case OrderByDate:
		sql = dateQuery
	}

	rows, err := s.DB.Query(
		ctx,
		sql,
		userID,
		limit,
		offset,
	)

	if err != nil {
		logger.Error("error returning user balance with specified id: user does not exist", zap.Error(err))
		return nil, err
	}

	var rr []ReadUserHistoryResult
	for rows.Next() {
		var r ReadUserHistoryResult
		err := rows.Scan(&r.AccountID, &r.CashBook, &r.Amount, &r.Date, &r.Addressee, &r.Description)
		if err != nil {
			logger.Error("scanning row error", zap.Error(err))
			return nil, err
		}
		r.Amount = decimal.New(r.Amount.IntPart(), -2)
		rr = append(rr, r)
	}
	return rr, nil
}
