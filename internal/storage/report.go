package storage

import (
	"context"
	"fmt"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type report struct {
	serviceId int64
	sum       decimal.Decimal
}

func (s *Storage) MonthlyReport(ctx context.Context, year int64, month int64) ([][]string, error) {
	logger := s.Logger.With(zap.Int64("Year", year), zap.Int64("Month", month))
	logger.Debug("reading the consolidated report")

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			if errRollback := tx.Rollback(ctx); errRollback != nil {
				logger.Error("error rolls back the transaction")
			}
		}
	}()

	selectQuery := `SELECT service_id, sum(sum) FROM consolidated_report cr
						INNER JOIN posting ON cr.tx_id = posting.id AND (SELECT EXTRACT (MONTH FROM (SELECT date FROM posting WHERE id = cr.tx_id))) = $1 
						AND (SELECT EXTRACT (YEAR FROM (SELECT date FROM posting WHERE id = cr.tx_id))) = $2 
						GROUP BY service_id`

	rows, err := tx.Query(
		ctx,
		selectQuery,
		month,
		year,
	)
	if err != nil {
		if rows == nil {
			logger.Error("error returning consolidated report records: records do not exist", zap.Error(ErrNoRecords))
			return nil, ErrNoRecords
		}
		logger.Error("Query error", zap.Error(err))
		return nil, err
	}
	var ss = make([][]string, 0)
	ss = append(ss, []string{"service_id", "total_revenue"})
	for rows.Next() {
		var r report
		var s = make([]string, 0)
		err := rows.Scan(&r.serviceId, &r.sum)
		if err != nil {
			logger.Error("scanning row error", zap.Error(err))
			return nil, err
		}
		s = append(s, fmt.Sprintf(`%d`, r.serviceId), decimal.New(r.sum.IntPart(), -2).String())
		ss = append(ss, s)
	}
	err = tx.Commit(ctx)
	return ss, nil
}
