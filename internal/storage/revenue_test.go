package storage

import (
	"context"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRevenue(t *testing.T) {
	s := bootstrap(t)

	var expectedPostingTable = []testPosting{
		{
			Posting: ReadUserHistoryResult{
				AccountID: int64(2),
				CashBook:  OperationTypeDeposit,
				Amount:    decimal.NewFromInt(100000),
			},
		},
		{
			Posting: ReadUserHistoryResult{
				AccountID: cacheBookAccountID,
				CashBook:  OperationTypeDeposit,
				Amount:    decimal.NewFromInt(-100000),
			},
		},
		{
			Posting: ReadUserHistoryResult{
				AccountID: int64(2),
				CashBook:  OperationTypeTransfer,
				Amount:    decimal.NewFromInt(-10000),
			},
		},
		{
			Posting: ReadUserHistoryResult{
				AccountID: reserveAccountID,
				CashBook:  OperationTypeTransfer,
				Amount:    decimal.NewFromInt(10000),
			},
		},
		{
			Posting: ReadUserHistoryResult{
				AccountID: reserveAccountID,
				CashBook:  OperationTypeTransfer,
				Amount:    decimal.NewFromInt(-10000),
			},
		},
		{
			Posting: ReadUserHistoryResult{
				AccountID: cacheBookAccountID,
				CashBook:  OperationTypeTransfer,
				Amount:    decimal.NewFromInt(10000),
			},
		},
	}
	err := s.Deposit(context.Background(), 2, decimal.NewFromInt(100000))
	require.NoError(t, err)

	description := "test"
	err = s.Reservation(context.Background(), 2, 2, 2, decimal.NewFromInt(10000), &description)
	require.NoError(t, err)

	err = s.Revenue(context.Background(), 2, 2, 2, decimal.NewFromInt(10000), &description)
	require.NoError(t, err)

	sql := "select * from posting"
	rows, err := s.DB.Query(context.Background(), sql)
	require.NoError(t, err)

	var pp []testPosting
	for rows.Next() {
		var p testPosting
		err := rows.Scan(
			&p.Id,
			&p.Posting.AccountID,
			&p.Posting.CashBook,
			&p.AccountingPeriod,
			&p.Posting.Amount,
			&p.Posting.Date,
			&p.Posting.Addressee,
			&p.Posting.Description)
		require.NoError(t, err)
		pp = append(pp, p)

	}

	assert.Len(t, pp, len(expectedPostingTable))
	dateTime := time.Now()

	for i, expectedRecord := range expectedPostingTable {
		assert.Equal(t, expectedRecord.Posting.AccountID, pp[i].Posting.AccountID)
		assert.Equal(t, expectedRecord.Posting.CashBook, pp[i].Posting.CashBook)
		assert.Equal(t, expectedRecord.Posting.Amount, pp[i].Posting.Amount)
		assert.LessOrEqual(t, dateTime.Format(time.RFC3339), pp[i].Posting.Date.Format(time.RFC3339))
		assert.LessOrEqual(t, dateTime.Format("2006-01-02"), pp[i].AccountingPeriod.Format("2006-01-02"))
	}
}

func TestRevenueOrderAlreadyExists(t *testing.T) {
	s := bootstrap(t)

	err := s.Deposit(context.Background(), 2, decimal.NewFromInt(100000))
	require.NoError(t, err)

	description := "test"
	err = s.Reservation(context.Background(), 2, 2, 2, decimal.NewFromInt(10000), &description)
	require.NoError(t, err)

	err = s.Reservation(context.Background(), 2, 2, 3, decimal.NewFromInt(10000), &description)
	require.NoError(t, err)

	err = s.Revenue(context.Background(), 2, 2, 2, decimal.NewFromInt(10000), &description)
	require.NoError(t, err)

	err = s.Revenue(context.Background(), 2, 2, 2, decimal.NewFromInt(10000), &description)
	require.ErrorIs(t, ErrRecordExist, err)
}

func TestUnreservationOrderExists(t *testing.T) {
	s := bootstrap(t)

	err := s.Deposit(context.Background(), 2, decimal.NewFromInt(100000))
	require.NoError(t, err)

	description := "test"
	err = s.Reservation(context.Background(), 2, 2, 2, decimal.NewFromInt(10000), &description)
	require.NoError(t, err)

	err = s.Reservation(context.Background(), 2, 2, 3, decimal.NewFromInt(10000), &description)
	require.NoError(t, err)

	err = s.Unreservation(context.Background(), 2, 2, 2, &description)
	require.NoError(t, err)

	err = s.Revenue(context.Background(), 2, 2, 2, decimal.NewFromInt(10000), &description)
	require.ErrorIs(t, ErrRecordExist, err)
}
