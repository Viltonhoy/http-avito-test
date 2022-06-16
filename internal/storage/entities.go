package storage

import (
	"time"

	"github.com/shopspring/decimal"
)

type UserBalance struct {
	AccountID int64           `json:"userID"`
	Balance   decimal.Decimal `json:"balance"`
}

type ReadUserHistoryResult struct {
	AccountID   int64           `json:"userID"`
	CBjournal   OperationType   `json:"cbJournal"`
	Amount      decimal.Decimal `json:"amount"`
	Date        time.Time       `json:"date"`
	Addressee   *int64          `json:"addressee"`
	Description *string         `json:"description"`
}

type OperationType string

const (
	OperationTypeDeposit    OperationType = "deposit"
	OperationTypeWithdrawal OperationType = "withdrawal"
	OperationTypeTransfer   OperationType = "transfer"
)
