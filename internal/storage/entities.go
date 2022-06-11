package storage

import (
	"time"

	"github.com/shopspring/decimal"
)

type UserBalance struct {
	AccountID int64           `json:"ID"`
	Balance   decimal.Decimal `json:"Balance"`
}

type Transfer struct {
	AcountID    int64           `json:"ID"`
	CBjournal   string          `json:"Type"`
	Amount      decimal.Decimal `json:"Sum"`
	Date        time.Time       `json:"Date"`
	Addressee   *int64          `json:"Addressee"`
	Description *string         `json:"Description"`
}

type operationType string

const (
	operationTypeDeposit    operationType = "deposit"
	operationTypeWithdrawal operationType = "withdrawal"
	operationTypeTransfer   operationType = "transfer"
)
