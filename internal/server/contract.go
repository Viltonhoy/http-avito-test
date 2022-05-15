package server

import (
	"context"
	"http-avito-test/internal/storage"

	"github.com/shopspring/decimal"
)

type Storager interface {
	ReadClient(int64, context.Context) (storage.User, error)
	Deposit(int64, decimal.Decimal, context.Context) error
	Withdrawal(int64, decimal.Decimal, context.Context) error
	Transfer(user_id1, user_id2 int64, amount decimal.Decimal, ctx context.Context) error
	ReadUserHistoryList(user_id int64, sort string, ctx context.Context) (l []storage.Transf, err error)
}
