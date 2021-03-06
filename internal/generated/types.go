// Package generated provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version (devel) DO NOT EDIT.
package generated

import (
	"http-avito-test/internal/storage"

	"github.com/shopspring/decimal"
)

// AccountDepositRequest defines model for AccountDepositRequest.
type AccountDepositRequest struct {
	Amount float32 `json:"amount"`
	UserId int     `json:"user_id"`
}

// AccountDepositResponse defines model for AccountDepositResponse.
type AccountDepositResponse struct {
	Result struct {
		Message string `json:"message"`
	} `json:"result"`
	Status string `json:"status"`
}

// AccountWithdrawalRequest defines model for AccountWithdrawalRequest.
type AccountWithdrawalRequest struct {
	Amount      float32 `json:"amount"`
	Description *string `json:"description"`
	UserId      int     `json:"user_id"`
}

// AccountWithdrawalResponse defines model for AccountWithdrawalResponse.
type AccountWithdrawalResponse = AccountDepositResponse

// ReadUserHistoryRequest defines model for ReadUserHistoryRequest.
type ReadUserHistoryRequest struct {
	Limit  int           `json:"limit"`
	Offset int           `json:"offset"`
	Order  storage.OrdBy `json:"order"`
	UserId int           `json:"user_id"`
}

// ReadUserHistoryResponse defines model for ReadUserHistoryResponse.
type ReadUserHistoryResponse struct {
	Result []storage.ReadUserHistoryResult `json:"result"`
	Status string                          `json:"status"`
}

// ReadUserRequest defines model for ReadUserRequest.
type ReadUserRequest struct {
	Currency *string `json:"currency"`
	UserId   int    `json:"user_id"`
}

// ReadUserResponse defines model for ReadUserResponse.
type ReadUserResponse struct {
	Result struct {
		Balance decimal.Decimal `json:"balance"`
		UserId  int             `json:"user_id"`
	} `json:"result"`
	Status string `json:"status"`
}

// TransferCommandRequest defines model for TransferCommandRequest.
type TransferCommandRequest struct {
	Amount      float32 `json:"amount"`
	Description *string `json:"description"`
	Recipient   int     `json:"recipient"`
	Sender      int     `json:"sender"`
}

// TransferCommandResponse defines model for TransferCommandResponse.
type TransferCommandResponse = AccountDepositResponse

// Version defines model for Version.
type Version = int

// AccountDepositJSONBody defines parameters for AccountDeposit.
type AccountDepositJSONBody = AccountDepositRequest

// AccountWithdrawalJSONBody defines parameters for AccountWithdrawal.
type AccountWithdrawalJSONBody = AccountWithdrawalRequest

// ReadUserJSONBody defines parameters for ReadUser.
type ReadUserJSONBody = ReadUserRequest

// ReadUserHistoryJSONBody defines parameters for ReadUserHistory.
type ReadUserHistoryJSONBody = ReadUserHistoryRequest

// TransferCommandJSONBody defines parameters for TransferCommand.
type TransferCommandJSONBody = TransferCommandRequest

// AccountDepositJSONRequestBody defines body for AccountDeposit for application/json ContentType.
type AccountDepositJSONRequestBody = AccountDepositJSONBody

// AccountWithdrawalJSONRequestBody defines body for AccountWithdrawal for application/json ContentType.
type AccountWithdrawalJSONRequestBody = AccountWithdrawalJSONBody

// ReadUserJSONRequestBody defines body for ReadUser for application/json ContentType.
type ReadUserJSONRequestBody = ReadUserJSONBody

// ReadUserHistoryJSONRequestBody defines body for ReadUserHistory for application/json ContentType.
type ReadUserHistoryJSONRequestBody = ReadUserHistoryJSONBody

// TransferCommandJSONRequestBody defines body for TransferCommand for application/json ContentType.
type TransferCommandJSONRequestBody = TransferCommandJSONBody
