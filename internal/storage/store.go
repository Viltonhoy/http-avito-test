package storage

import (
	"context"
	"errors"
	"fmt"

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

func (s *Storage) ReadClient(user_id int64) (u User, err error) {
	logger := s.logger
	logger.Sugar().Debug(`reading the balance of %d user`, user_id)

	ctx := context.Background()
	tsql := "SELECT * FROM bankacc WHERE user_id = $1"

	//выполнить запрос
	err = s.db.QueryRow(ctx, tsql, user_id).Scan(&u.ID, &u.Balance)
	if err != nil {
		logger.Sugar().Error("cannot return user with specified ID: %d", user_id)
		return User{}, err
	}

	return User{
		u.ID,
		u.Balance,
	}, nil
}

func (s *Storage) UpdateClient(user_id int64, balance decimal.Decimal) error {
	logger := s.logger
	logger.Sugar().Debugf(`Updating users account information: ID: %d, Balance: %d`, user_id, balance)

	ctx := context.Background()
	tsql := `INSERT INTO bankacc (user_id, balance)
			VALUES ($1, $2) ON CONFLICT (user_id)
 			DO UPDATE SET balance = (SELECT balance + $2 FROM bankacc WHERE user_id = $1)
 			WHERE bankacc.user_id = $1;`

	_, err := s.db.Exec(
		ctx,
		tsql,
		user_id,
		balance,
	)
	if err != nil {
		logger.Sugar().Error("failed to update client balance: %v", err)
		return err
	}

	return err
}

func (s *Storage) MoneyTransfer(user_id1, user_id2 int64, balance decimal.Decimal) error {
	logger := s.logger
	logger.Sugar().Debugf(`money transfer from %d user to %d`, user_id1, user_id2)

	ctx := context.Background()
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}

	ftsql := `UPDATE bankacc SET balance = (SELECT balance - $2 FROM bankacc WHERE user_id = $1);`

	_, err = tx.Exec(
		ctx,
		ftsql,
		user_id1,
		balance,
	)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	stsql := `INSERT INTO bankacc (user_id, balance)
	 	VALUES ($1, $2) ON CONFLICT (user_id)
	 	DO UPDATE SET balance = (SELECT balance + $2 FROM bankacc WHERE user_id = $1)
	 	WHERE bankacc.user_id = $1;`

	_, err = tx.Exec(
		ctx,
		stsql,
		user_id2,
		balance,
	)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	err = tx.Commit(ctx)
	return err
}

// func NewStore(logger *zap.SugaredLogger) (*Storage, error) {
// 	if logger == nil {
// 		return nil, errors.New("no logger provided")
// 	}

// 	conf := NewServ()
// 	//строка подключения
// 	var connString = fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", conf.User, conf.Password, conf.Database)

// 	//создать пул соединений
// 	var db, err = sql.Open("postgres", connString)
// 	if err != nil {
// 		logger.Errorf("cannot connect using connection connString %s: %w", connString, err)
// 		return nil, err
// 	}
// 	ctx := context.Background()
// 	err = db.PingContext(ctx)
// 	if err != nil {
// 		logger.Fatal("connection is lost", err)
// 	}
// 	logger.Info("connected!\n")
// 	return &Storage{db: db, logger: logger}, nil
// }

// func (s *Storage) Close() {
// 	s.logger.Info("closing Storage connection")
// 	s.db.Close()
// }

// func (s *Storage) ReadClient(user_id int64) (u User, err error) {
// 	logger := s.logger
// 	logger.Debugf(`reading the balance of %d user`, user_id)

// 	ctx := context.Background()
// 	tsql := "SELECT * FROM bankacc WHERE user_id = $1"

// 	//выполнить запрос
// 	rows, err := s.db.QueryContext(ctx, tsql, user_id)
// 	if err != nil {
// 		logger.Errorf("cannot return user with specified ID: %d", user_id)
// 		return User{}, err
// 	}

// 	defer rows.Close()

// 	//итерация наборов результатов
// 	for rows.Next() {
// 		//получить значение из строки
// 		err := rows.Scan(&u.ID, &u.Balance)
// 		if err != nil {
// 			logger.Error(`cannot copy the columns in the current row into the values`)
// 			return User{}, err
// 		}
// 	}

// 	return User{u.ID, u.Balance}, nil
// }

// func (s *Storage) UpdateClient(user_id int64, balance decimal.Decimal) error {
// 	logger := s.logger
// 	logger.Debugf(`Updating users account information: ID: %d, Balance: %d`, user_id, balance)

// 	ctx := context.Background()
// 	tsql := `INSERT INTO bankacc (user_id, balance)
// 			VALUES ($1, $2) ON CONFLICT (user_id)
// 			DO UPDATE SET balance = (SELECT balance + $2 FROM bankacc WHERE user_id = $1)
// 			WHERE bankacc.user_id = $1;`

// 	_, err := s.db.ExecContext(
// 		ctx,
// 		tsql,
// 		user_id,
// 		balance,
// 	)
// 	if err != nil {
// 		logger.Error("failed to update client balance: %v", err)
// 		return err
// 	}

// 	// htsql := `INSERT INTO transfer_history () VALUES ();`

// 	// _, err = s.db.ExecContext(
// 	// 	ctx,
// 	// 	htsql,
// 	// )

// 	return err
// }

// func (s *Storage) MoneyTransfer(user_id1, user_id2 int64, balance decimal.Decimal) error {
// 	logger := s.logger
// 	logger.Debugf(`money transfer from %d user to %d`, user_id1, user_id2)

// 	ctx := context.Background()
// 	tx, err := s.db.Begin()
// 	if err != nil {
// 		return err
// 	}

// 	ftsql := `UPDATE bankacc SET balance = (SELECT balance - $2 FROM bankacc WHERE user_id = $1);`

// 	_, err = tx.ExecContext(
// 		ctx,
// 		ftsql,
// 		user_id1,
// 		balance)
// 	if err != nil {
// 		tx.Rollback()
// 		return err
// 	}

// 	stsql := `INSERT INTO bankacc (user_id, balance)
// 	VALUES ($1, $2) ON CONFLICT (user_id)
// 	DO UPDATE SET balance = (SELECT balance + $2 FROM bankacc WHERE user_id = $1)
// 	WHERE bankacc.user_id = $1;`

// 	_, err = tx.ExecContext(
// 		ctx,
// 		stsql,
// 		user_id2,
// 		balance)
// 	if err != nil {
// 		tx.Rollback()
// 		return err
// 	}

// 	err = tx.Commit()
// 	return err
// }

// func (s *Storage) ReadTransfHistoryList(user_id int64) (l TransfList, err error) {
// 	logger := s.logger
// 	logger.Debug("reading a history list")

// 	ctx := context.Background()
// 	tsql := "SELECT * FROM transfer_history WHERE user_id = $1;"

// 	rows, err := s.db.QueryContext(ctx, tsql, user_id)
// 	if err != nil {
// 		logger.Errorf("cannot return user list with specified ID: %d", user_id)
// 		return TransfList{}, err
// 	}

// 	defer rows.Close()

// 	for rows.Next() {
// 		err := scan.Rows(&l, rows)
// 		if err != nil {
// 			logger.Error(`cannot copy the columns in the current row into the values`)
// 			return TransfList{}, err
// 		}
// 	}
// 	return l, nil
// }
