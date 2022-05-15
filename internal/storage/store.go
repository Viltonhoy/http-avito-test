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

func NewStore(logger *zap.Logger) (*Storage, error) {
	if logger == nil {
		return nil, errors.New("no logger provided")
	}

	conf := NewServ()
	//строка подключения
	var connString = fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", conf.User, conf.Password, conf.Database)
	config, _ := pgxpool.ParseConfig(connString)

	config.ConnConfig.Logger = zapadapter.NewLogger(logger)
	config.ConnConfig.LogLevel = pgx.LogLevelError

	ctx := context.Background()
	// 	//создать пул соединений
	pool, _ := pgxpool.ConnectConfig(ctx, config)

	err := pool.Ping(ctx)
	if err != nil {
		logger.Sugar().Fatalf("connection is lost", err)
	}

	return &Storage{
		logger: logger,
		db:     pool,
	}, nil
}

func (s *Storage) Close() {
	s.logger.Info("closing Storage connection")
	s.db.Close()
}

func (s *Storage) ReadClient(user_id int64, ctx context.Context) (u User, err error) {
	logger := s.logger
	logger.Sugar().Debug(`reading the balance of %d user`, user_id)

	tsql := `SELECT SUM(amount) FROM posting WHERE account_id = $1;`

	//выполнить запрос
	err = s.db.QueryRow(ctx, tsql, user_id).Scan(&u.Balance)
	if err != nil {
		logger.Sugar().Error("cannot return user with specified ID: %d", user_id)
		return User{}, err
	}

	u.ID = user_id

	return User{
		u.ID,
		u.Balance,
	}, nil
}

func (s *Storage) Deposit(user_id int64, amount decimal.Decimal, ctx context.Context) (err error) {
	logger := s.logger
	logger.Sugar().Debugf(`Updating users account information: ID: %d, Balance: %d`, user_id, amount)

	var tp = "deposit"
	var t = time.Now()

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}

	ftsql := `INSERT INTO posting (account_id, cb_journal, accounting_period, amount, date, addressee)
			VALUES ($1, $4, $5, $2, $3, '');`

	_, err = tx.Exec(
		ctx,
		ftsql,
		user_id,
		amount,
		t.Format("2006-01-02 15:04:05"),
		tp,
		t.Format("2006"),
	)
	if err != nil {
		tx.Rollback(ctx)
		logger.Sugar().Error("failed to insert record: %v", err)
		return err
	}

	stsql := `INSERT INTO posting (account_id, cb_journal, accounting_period, amount, date, addressee)
			VALUES (0, $3, 2022, -1 * $1, $2, '');`

	_, err = tx.Exec(
		ctx,
		stsql,
		amount,
		t.Format("2006-01-02 15:04:05"),
		tp,
	)
	if err != nil {
		tx.Rollback(ctx)
		logger.Sugar().Error("failed to insert record: %v", err)
		return err
	}
	err = tx.Commit(ctx)
	return err
}

func (s *Storage) Withdrawal(user_id int64, amount decimal.Decimal, ctx context.Context) (err error) {
	logger := s.logger
	logger.Sugar().Debugf(`Updating users account information: ID: %d, Balance: %d`, user_id, amount)

	var tp = "withdrawal"
	var t = time.Now()

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}

	check_sql := `REFRESH MATERIALIZED VIEW account_balances;`

	_, err = tx.Exec(ctx, check_sql)
	if err != nil {
		tx.Rollback(ctx)
		logger.Sugar().Error("failed to refresh materialized view account_balances", err)
		return err
	}

	select_sql := `SELECT * FROM account_balances WHERE user_id = $1;`

	var balance Balance
	err = tx.QueryRow(ctx, select_sql, user_id).Scan(&balance.ID, &balance.Sum)
	if err != nil {
		tx.Rollback(ctx)
		logger.Sugar().Error("failed ")
	}

	if amount.GreaterThan(balance.Sum) {
		tx.Rollback(ctx)
		logger.Error("insufficient funds on the user's account")
		return err
	}

	ftsql := `INSERT INTO posting (account_id, cb_journal, accounting_period, amount, date, addressee)
			VALUES ($1, $4, $5, -1 * $2, $3, '');`

	_, err = tx.Exec(
		ctx,
		ftsql,
		user_id,
		amount,
		t.Format("2006-01-02 15:04:05"),
		tp,
		t.Format("2006"),
	)
	if err != nil {
		tx.Rollback(ctx)
		logger.Sugar().Error("failed to insert record: %v", err)
		return err
	}

	stsql := `INSERT INTO posting (account_id, cb_journal, accounting_period, amount, date, addressee)
			VALUES (0, $3, 2022, $1, $2, '');`

	_, err = tx.Exec(
		ctx,
		stsql,
		amount,
		t.Format("2006-01-02 15:04:05"),
		tp,
	)
	if err != nil {
		tx.Rollback(ctx)
		logger.Sugar().Error("failed to insert record: %v", err)
		return err
	}
	err = tx.Commit(ctx)
	return err
}

func (s *Storage) Transfer(user_id1, user_id2 int64, amount decimal.Decimal, ctx context.Context) error {
	logger := s.logger
	logger.Sugar().Debugf(`money transfer from %d user to %d`, user_id1, user_id2)

	var tp = "transfer"
	var t = time.Now()

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}

	check_sql := `REFRESH MATERIALIZED VIEW account_balances;`

	_, err = tx.Exec(ctx, check_sql)
	if err != nil {
		tx.Rollback(ctx)
		logger.Sugar().Error("failed to refresh materialized view account_balances", err)
		return err
	}

	select_sql := `SELECT * FROM account_balances WHERE user_id = $1;`

	var balance Balance
	err = tx.QueryRow(ctx, select_sql, user_id1).Scan(&balance.ID, &balance.Sum)
	if err != nil {
		tx.Rollback(ctx)
		logger.Sugar().Error("failed ")
	}

	if amount.GreaterThan(balance.Sum) {
		tx.Rollback(ctx)
		logger.Error("insufficient funds on the user's account")
		return err
	}

	ftsql := `INSERT INTO posting (account_id, cb_journal, accounting_period, amount, date, addressee) VALUES ($1, $6, $4, -1 * $2, $3, $5);`

	_, err = tx.Exec(
		ctx,
		ftsql,
		user_id1,
		amount,
		t.Format("2006-01-02 15:04:05"),
		t.Format("2006"),
		fmt.Sprintf(`account_id: %d`, user_id2),
		tp,
	)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	stsql := `INSERT INTO posting (account_id, cb_journal, accounting_period, amount, date, addressee) VALUES ($1, $6, $4, $2, $3, $5);`

	_, err = tx.Exec(
		ctx,
		stsql,
		user_id2,
		amount,
		t.Format("2006-01-02 15:04:05"),
		t.Format("2006"),
		fmt.Sprintf(`account_id: %d`, user_id1),
		tp,
	)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	err = tx.Commit(ctx)
	return err
}

func (s *Storage) ReadUserHistoryList(user_id int64, sort string, ctx context.Context) (l []Transf, err error) {
	logger := s.logger
	logger.Debug("reading a history list")

	tsql := "SELECT account_id, cb_journal, amount, date, addressee FROM posting WHERE account_id = $1 ORDER BY $2;"

	rows, err := s.db.Query(ctx, tsql, user_id, sort)
	if err != nil {
		logger.Sugar().Errorf("cannot return user list with specified ID: %d", user_id)
		return TransfList{}, err
	}

	var list []Transf
	for rows.Next() {
		var l Transf
		err := rows.Scan(&l.ID, &l.Type, &l.Sum, &l.Date, &l.Addressee)
		if err != nil {
			s.logger.Error("scanning row", zap.Error(err))
		}
		list = append(list, l)
	}
	return list, nil
}
