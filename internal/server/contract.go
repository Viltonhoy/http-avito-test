package server

import (
	"context"
	"http-avito-test/internal/storage"

	"github.com/shopspring/decimal"
)

type Storager interface {
	ReadUser(context.Context, int64) (storage.UserBalance, error)
	Deposit(context.Context, int64, decimal.Decimal) error
	Withdrawal(context.Context, int64, decimal.Decimal) error
	Transfer(ctx context.Context, user_id1, user_id2 int64, amount decimal.Decimal) error
	ReadUserHistoryList(ctx context.Context, user_id int64, order string, limit, offset int64) (l []storage.Transfer, err error)
}
