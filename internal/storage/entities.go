package storage

import (
	"time"

	"github.com/shopspring/decimal"
)

type User struct {
	AccountID int64           `json:"ID"`
	Balance   decimal.Decimal `json:"Balance"`
}

type Transf struct {
	AcountID  int64     `json:"ID"`
	CBjournal string    `json:"Type"`
	Amount    int64     `json:"Sum"`
	Date      time.Time `json:"Date"`
	Addressee string    `json:"Addressee"`
}

type OperationType string

const (
	operationTypeDeposit    OperationType = "deposit"
	operationTypeWithdrawal OperationType = "withdrawal"
	operationTypeTransfer   OperationType = "transfer"
)

type TransfList []Transf

type UserBalance struct {
	AccountID int64
	Balance   decimal.Decimal
}
