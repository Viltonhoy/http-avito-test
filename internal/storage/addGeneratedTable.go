package storage

import (
	"context"
	"fmt"
	generatetable "http-avito-test/internal/generateTable"
	"time"

	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"
)

func AddGeneratedTableData(s *Storage, userCount, totalRecordCount int) {
	logger := s.logger
	logger.Info(`add new rows for users to database`, zap.Int("totalRecordCount", totalRecordCount), zap.Int("userCount", userCount))

	columnName := []string{"account_id", "cb_journal", "accounting_period", "amount", "date", "addressee"}

	var rows = generatetable.GenerateTableData(userCount, totalRecordCount)
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
		logger.Error("cannot add new rows", zap.Error(err))
		return
	}
	duration := time.Since(start)
	s.db.Close()
	fmt.Printf(`The number of copied rows: %d;  request execution time: %fsec`, num, duration.Seconds())
}
