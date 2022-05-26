package generatetable

import (
	"context"
	"errors"
	"fmt"
	"http-avito-test/internal/storage"
	"http-avito-test/internal/zapadapter"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
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

	conf := storage.NewServ()
	var connString = fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", conf.User, conf.Password, conf.Database)
	config, _ := pgxpool.ParseConfig(connString)

	config.ConnConfig.Logger = zapadapter.NewLogger(logger)
	config.ConnConfig.LogLevel = pgx.LogLevelError

	ctx := context.Background()
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

func AddGeneratedTable(s *Storage, userCount, totalRecordCount int) {
	logger := s.logger
	logger.Sugar().Debug(`add %d new rows for %d users to database`, totalRecordCount, userCount)

	columnName := []string{"account_id", "cb_journal", "accounting_period", "amount", "date", "addressee"}

	var rows = GenerateTableData(userCount, totalRecordCount)
	newSlice := make([][]interface{}, 0, len(rows))

	for _, row := range rows {
		newSlice = append(
			newSlice,
			[]interface{}{
				row.AccountID,
				row.CBjournal,
				row.AccountingPeriod,
				row.Amount,
				row.Date,
				row.Addressee,
			},
		)
	}
	start := time.Now()
	num, err := s.db.CopyFrom(context.Background(), pgx.Identifier{"posting"}, columnName, pgx.CopyFromRows(newSlice))
	if err != nil {
		logger.Error("")
		return
	}
	duration := time.Since(start)
	s.db.Close()
	fmt.Printf(`The number of copied rows: %d;  request execution time: %fsec`, num, duration.Seconds())
}
