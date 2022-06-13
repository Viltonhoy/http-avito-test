package storage

import (
	"time"

	"github.com/shopspring/decimal"
)

type UserBalance struct {
	AccountID int64           `json:"UserID"`
	Balance   decimal.Decimal `json:"Balance"`
}

type ReadUserHistoryResult struct {
	AcountID    int64           `json:"userID"`
	CBjournal   string          `json:"cbJournal"`
	Amount      decimal.Decimal `json:"amount"`
	Date        time.Time       `json:"date"`
	Addressee   *int64          `json:"addressee"`
	Description *string         `json:"description"`
}

type operationType string

const (
	operationTypeDeposit    operationType = "deposit"
	operationTypeWithdrawal operationType = "withdrawal"
	operationTypeTransfer   operationType = "transfer"
)
