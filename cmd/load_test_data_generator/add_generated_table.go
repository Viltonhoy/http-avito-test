package main

import (
	"context"
	"fmt"
	"http-avito-test/internal/storage"
	"time"

	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"
)

// AddGeneratedTableData creates and populates a slice of values ​​to populate the posting table
// using the input values ​​for the number of users and the number of records
func AddGeneratedTableData(s *storage.Storage, userCount, totalRecordCount int) {
	s.Logger.Debug(`add new rows for users to database`, zap.Int("totalRecordCount", totalRecordCount), zap.Int("userCount", userCount))

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
	num, err := s.DB.CopyFrom(context.Background(), pgx.Identifier{"posting"}, columnName, pgx.CopyFromRows(newSlice))
	if err != nil {
		s.Logger.Error("cannot add new rows", zap.Error(err))
		return
	}
	duration := time.Since(start)
	s.DB.Close()
	fmt.Printf(`The number of copied rows: %d;  request execution time: %f sec`, num, duration.Seconds())
}
