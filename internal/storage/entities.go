package storage

import (
	"database/sql"
	"time"

	"github.com/shopspring/decimal"
)

type User struct {
	AccountID int64           `json:"userID"`
	Balance   decimal.Decimal `json:"balance"`
}

type ReadUserHistoryResult struct {
	AccountID   int64           `json:"userID"`
	CashBook    OperationType   `json:"cashebook"`
	Amount      decimal.Decimal `json:"amount"`
	Date        time.Time       `json:"date"`
	Addressee   sql.NullInt64   `json:"addressee"`
	Description sql.NullString  `json:"description"`
}

type OperationType string

const (
	OperationTypeDeposit    OperationType = "deposit"
	OperationTypeWithdrawal OperationType = "withdrawal"
	OperationTypeTransfer   OperationType = "transfer"
)

type OrdBy string

const (
	OrderByAmount OrdBy = "amount"
	OrderByDate   OrdBy = "date"
)
