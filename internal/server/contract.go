package server

import (
	"http-avito-test/internal/storage"

	"github.com/shopspring/decimal"
)

type Storager interface {
	ReadClient(int64) (storage.User, error)
	DepositOrWithdrawal(int64, decimal.Decimal, string) error
	Transfer(user_id1, user_id2 int64, amount decimal.Decimal) error
	ReadUserHistoryList(user_id int64, sort string) (l []storage.Transf, err error)
}
