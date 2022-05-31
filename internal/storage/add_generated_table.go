package storage

import (
	"context"
	"fmt"
	generatetable "http-avito-test/internal/generateTable"
	"time"

	"github.com/jackc/pgx/v4"
)

func AddGeneratedTable(s *Storage, userCount, totalRecordCount int) {
	logger := s.logger
	logger.Sugar().Debug(`add %d new rows for %d users to database`, totalRecordCount, userCount)

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
		logger.Error("")
		return
	}
	duration := time.Since(start)
	s.db.Close()
	fmt.Printf(`The number of copied rows: %d;  request execution time: %fsec`, num, duration.Seconds())
}
