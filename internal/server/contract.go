package server

import (
	"context"
	"http-avito-test/internal/storage"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type Storager interface {
	ReadUserByID(context.Context, int64) (storage.User, error)
	Deposit(context.Context, int64, decimal.Decimal) error
	Withdrawal(context.Context, int64, decimal.Decimal, *string) error
	Transfer(ctx context.Context, user_id1, user_id2 int64, amount decimal.Decimal, description *string) error
	ReadUserHistoryList(ctx context.Context, user_id int64, order storage.OrdBy, limit, offset int64) ([]storage.ReadUserHistoryResult, error)
}

type Exchanger interface {
	ExchangeRates(logger *zap.Logger, value decimal.Decimal, currency string) (decimal.Decimal, error)
}
