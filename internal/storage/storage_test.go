package storage

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type testPosting struct {
	Posting          ReadUserHistoryResult
	AccountingPeriod time.Time
	Id               int64
}

func bootstrap(t *testing.T) *Storage {
	err := godotenv.Load("../../.env")
	require.NoError(t, err)

	logger, err := zap.NewDevelopment()
	require.NoError(t, err)
	s, err := NewStorage(context.Background(), logger)
	require.NoError(t, err)

	truncate := `TRUNCATE posting, balances RESTART IDENTITY;`

	_, err = s.DB.Exec(context.Background(), truncate)
	require.NoError(t, err)
	return s
}

func TestDeposit(t *testing.T) {
	s := bootstrap(t)

	var expectedPostingTable = []testPosting{
		{
			Posting: ReadUserHistoryResult{
				AccountID: int64(1),
				CashBook:  OperationTypeDeposit,
				Amount:    decimal.NewFromInt(10000),
			},
		},
		{
			Posting: ReadUserHistoryResult{
				AccountID: cacheBookAccountID,
				CashBook:  OperationTypeDeposit,
				Amount:    decimal.NewFromInt(-10000),
			},
		},
	}

	err := s.Deposit(context.Background(), 1, decimal.NewFromInt(10000))
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
func TestWithdrawal(t *testing.T) {
	s := bootstrap(t)

	var expectedPostingTable = []testPosting{
		{
			Posting: ReadUserHistoryResult{
				AccountID: int64(1),
				CashBook:  OperationTypeDeposit,
				Amount:    decimal.NewFromInt(10000),
			},
		},
		{
			Posting: ReadUserHistoryResult{
				AccountID: cacheBookAccountID,
				CashBook:  OperationTypeDeposit,
				Amount:    decimal.NewFromInt(-10000),
			},
		},
		{
			Posting: ReadUserHistoryResult{
				AccountID: int64(1),
				CashBook:  OperationTypeWithdrawal,
				Amount:    decimal.NewFromInt(-10000),
			},
		},
		{
			Posting: ReadUserHistoryResult{
				AccountID: cacheBookAccountID,
				CashBook:  OperationTypeWithdrawal,
				Amount:    decimal.NewFromInt(10000),
			},
		},
	}
	err := s.Deposit(context.Background(), 1, decimal.NewFromInt(10000))
	require.NoError(t, err)

	description := "test"
	err = s.Withdrawal(context.Background(), 1, decimal.NewFromInt(10000), &description)
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

func TestTransfer(t *testing.T) {
	s := bootstrap(t)

	var expectedPostingTable = []testPosting{
		{
			Posting: ReadUserHistoryResult{
				AccountID: int64(1),
				CashBook:  OperationTypeDeposit,
				Amount:    decimal.NewFromInt(10000),
			},
		},
		{
			Posting: ReadUserHistoryResult{
				AccountID: cacheBookAccountID,
				CashBook:  OperationTypeDeposit,
				Amount:    decimal.NewFromInt(-10000),
			},
		},
		{
			Posting: ReadUserHistoryResult{
				AccountID: int64(1),
				CashBook:  OperationTypeTransfer,
				Amount:    decimal.NewFromInt(-10000),
			},
		},
		{
			Posting: ReadUserHistoryResult{
				AccountID: int64(2),
				CashBook:  OperationTypeTransfer,
				Amount:    decimal.NewFromInt(10000),
			},
		},
	}
	err := s.Deposit(context.Background(), 1, decimal.NewFromInt(10000))
	require.NoError(t, err)

	description := "test"
	err = s.Transfer(context.Background(), 1, 2, decimal.NewFromInt(10000), &description)
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

func TestReadUserById(t *testing.T) {
	s := bootstrap(t)
	expectBalance := decimal.NewFromInt(25000)

	err := s.Deposit(context.Background(), 1, decimal.NewFromInt(10000))
	require.NoError(t, err)

	err = s.Deposit(context.Background(), 1, decimal.NewFromInt(20000))
	require.NoError(t, err)

	err = s.Deposit(context.Background(), 1, decimal.NewFromInt(15000))
	require.NoError(t, err)

	description := "test"
	err = s.Withdrawal(context.Background(), 1, decimal.NewFromInt(10000), &description)
	require.NoError(t, err)

	err = s.Transfer(context.Background(), 1, 2, decimal.NewFromInt(10000), &description)
	require.NoError(t, err)

	user, err := s.ReadUserByID(context.Background(), 1)
	require.NoError(t, err)

	assert.Equal(t, expectBalance, user.Balance)
}

func TestReadUserHistory(t *testing.T) {
	s := bootstrap(t)

	var expectedPostingTable = []testPosting{
		{
			Posting: ReadUserHistoryResult{
				AccountID: int64(1),
				CashBook:  OperationTypeWithdrawal,
				Amount:    decimal.New(decimal.NewFromInt(-10000).IntPart(), -2),
			},
		},
		{
			Posting: ReadUserHistoryResult{
				AccountID: int64(1),
				CashBook:  OperationTypeTransfer,
				Amount:    decimal.New(decimal.NewFromInt(-10000).IntPart(), -2),
				Addressee: sql.NullInt64{
					Int64: 2,
					Valid: true,
				},
			},
		},
		{
			Posting: ReadUserHistoryResult{
				AccountID: int64(1),
				CashBook:  OperationTypeDeposit,
				Amount:    decimal.New(decimal.NewFromInt(10000).IntPart(), -2),
			},
		},
		{
			Posting: ReadUserHistoryResult{
				AccountID: int64(1),
				CashBook:  OperationTypeDeposit,
				Amount:    decimal.New(decimal.NewFromInt(15000).IntPart(), -2),
			},
		},
		{
			Posting: ReadUserHistoryResult{
				AccountID: int64(1),
				CashBook:  OperationTypeDeposit,
				Amount:    decimal.New(decimal.NewFromInt(20000).IntPart(), -2),
			},
		},
	}

	err := s.Deposit(context.Background(), 1, decimal.NewFromInt(10000))
	require.NoError(t, err)

	err = s.Deposit(context.Background(), 1, decimal.NewFromInt(20000))
	require.NoError(t, err)

	err = s.Deposit(context.Background(), 1, decimal.NewFromInt(15000))
	require.NoError(t, err)

	description := "test"
	err = s.Withdrawal(context.Background(), 1, decimal.NewFromInt(10000), &description)
	require.NoError(t, err)

	err = s.Transfer(context.Background(), 1, 2, decimal.NewFromInt(10000), &description)
	require.NoError(t, err)

	user, err := s.ReadUserHistoryList(context.Background(), 1, "amount", 100, 0)
	require.NoError(t, err)

	assert.Len(t, user, len(expectedPostingTable))
	dateTime := time.Now()

	for i, expectedRecord := range expectedPostingTable {
		assert.Equal(t, expectedRecord.Posting.AccountID, user[i].AccountID)
		assert.Equal(t, expectedRecord.Posting.CashBook, user[i].CashBook)
		assert.Equal(t, expectedRecord.Posting.Amount, user[i].Amount)
		assert.LessOrEqual(t, dateTime.Format(time.RFC3339), user[i].Date.Format(time.RFC3339))
		assert.Equal(t, expectedRecord.Posting.Addressee, user[i].Addressee)
	}
}

func TestZeroSumOfAmount(t *testing.T) {
	t.Run("sum", func(t *testing.T) {
		s := bootstrap(t)
		var totalAmount decimal.Decimal

		err := s.Deposit(context.Background(), 1, decimal.NewFromInt(10000))
		require.NoError(t, err)

		err = s.Deposit(context.Background(), 1, decimal.NewFromInt(20000))
		require.NoError(t, err)

		err = s.Deposit(context.Background(), 1, decimal.NewFromInt(15000))
		require.NoError(t, err)

		description := "test"
		err = s.Withdrawal(context.Background(), 1, decimal.NewFromInt(10000), &description)
		require.NoError(t, err)

		err = s.Transfer(context.Background(), 1, 2, decimal.NewFromInt(10000), &description)
		require.NoError(t, err)

		sql := "select sum(amount) from posting;"
		err = s.DB.QueryRow(context.Background(), sql).Scan(&totalAmount)
		require.NoError(t, err)

		assert.Equal(t, decimal.NewFromInt(0), totalAmount)
	})

	t.Run("sum with limit", func(t *testing.T) {
		s := bootstrap(t)
		var totalAmount decimal.Decimal

		err := s.Deposit(context.Background(), 1, decimal.NewFromInt(10000))
		require.NoError(t, err)

		err = s.Deposit(context.Background(), 1, decimal.NewFromInt(20000))
		require.NoError(t, err)

		err = s.Deposit(context.Background(), 1, decimal.NewFromInt(15000))
		require.NoError(t, err)

		description := "test"
		err = s.Withdrawal(context.Background(), 1, decimal.NewFromInt(10000), &description)
		require.NoError(t, err)

		err = s.Transfer(context.Background(), 1, 2, decimal.NewFromInt(10000), &description)
		require.NoError(t, err)

		sql := "select sum(amount) from posting limit 100;"
		err = s.DB.QueryRow(context.Background(), sql).Scan(&totalAmount)
		require.NoError(t, err)

		assert.Equal(t, decimal.NewFromInt(0), totalAmount)
	})
}

func TestNotEnoughMoney(t *testing.T) {
	t.Run("withdrawal", func(t *testing.T) {
		s := bootstrap(t)

		err := s.Deposit(context.Background(), 1, decimal.NewFromInt(10000))
		require.NoError(t, err)

		description := "test"

		err = s.Withdrawal(context.Background(), 1, decimal.NewFromInt(20000), &description)
		assert.ErrorIs(t, ErrWithdrawal, err)
	})

	t.Run("transfer", func(t *testing.T) {
		s := bootstrap(t)

		err := s.Deposit(context.Background(), 1, decimal.NewFromInt(10000))
		require.NoError(t, err)

		description := "test"

		err = s.Transfer(context.Background(), 1, 2, decimal.NewFromInt(20000), &description)
		assert.ErrorIs(t, ErrTransfer, err)
	})
}
