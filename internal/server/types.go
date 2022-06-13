package server

import (
	"http-avito-test/internal/storage"

	"github.com/shopspring/decimal"
)

// ReadUserRequest defines model for ReadUser
type ReadUserRequest struct {
	UserID   int64  `json:"userID"`
	Currency string `json:"currency"`
}

// ReadUserResponse defines model for ReadUser
type ReadUserResponse struct {
	Result struct {
		UserID  int64           `json:"userID"`
		Balance decimal.Decimal `json:"balance"`
	} `json:"result"`
	Status string `json:"status"`
}

// ReadUserHistoryRequest defines model for ReadUserHistory
type ReadUserHistoryRequest struct {
	UserID int64 `json:"userID"`
	Order  ordBy `json:"order"`
	Limit  int64 `json:"limit"`
	Offset int64 `json:"offset"`
}

// ReadUserHistoryResponse defines model for ReadUserHistory
type ReadUserHistoryResponse struct {
	Result []storage.ReadUserHistoryResult `json:"result"`
	Status string                          `json:"status"`
}

// AccountDepositRequest defines model for AccountDeposit
type AccountDepositRequest struct {
	UserID int64   `json:"userID"`
	Amount float32 `json:"amount"`
}

const resultMessage = "balance updated successfully"

// AccountDepositResponse defines model for AccountDeposit
type AccountDepositResponse struct {
	Result struct {
		Message string
	} `json:"result"`
	Status string `json:"status"`
}

// AccountWithdrawalRequest defines model for AccountWithdrawal
type AccountWithdrawalRequest struct {
	UserID      int64   `json:"userID"`
	Amount      float32 `json:"amount"`
	Description *string `json:"description"`
}

// AccountWithdrawalResponse defines model for AccountWithdrawal
type AccountWithdrawalResponse struct {
	Result struct {
		Message string
	} `json:"result"`
	Status string `json:"status"`
}

// TransferCommandRequest defines model for TransferCommand
type TransferCommandRequest struct {
	UserID1     int64   `json:"userID1"`
	UserID2     int64   `json:"userID2"`
	Amount      float32 `json:"amount"`
	Description *string `json:"description"`
}

// TransferCommandResponse defines model for TransferCommand
type TransferCommandResponse struct {
	Result struct {
		Message string
	} `json:"result"`
	Status string `json:"status"`
}
