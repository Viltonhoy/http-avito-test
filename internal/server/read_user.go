package server

import (
	"encoding/json"
	"errors"
	"http-avito-test/internal/exchanger"
	"http-avito-test/internal/generated"
	"http-avito-test/internal/storage"
	"io/ioutil"
	"net/http"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

const (
	rubleCurrencyCode    = "RUB"
	oldRubleCurrensyCode = "RUR"
)

func (h *Handler) ReadUser(w http.ResponseWriter, r *http.Request) {
	var hand *generated.ReadUserRequest

	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	err := json.Unmarshal(body, &hand)
	if err != nil {
		http.Error(w, "malformed request body", http.StatusBadRequest)
		return
	}

	if hand.UserId <= 0 {
		http.Error(w, "wrong value of \"User_id\"", http.StatusBadRequest)
		return
	}

	if *hand.Currency == "" {
		http.Error(w, "incorrect currency code value", http.StatusBadRequest)
		return
	}

	if hand.Currency == nil {
		*hand.Currency = "RUB"
	}

	user, err := h.Store.ReadUserByID(r.Context(), int64(hand.UserId))
	if err != nil {
		if errors.Is(err, storage.ErrUserAvailability) {
			http.Error(w, "user does not exist", http.StatusBadRequest)
			return
		}
		http.Error(w, "cannot read user with specified id", http.StatusInternalServerError)
		return
	}

	var newBalance decimal.Decimal

	expBalance := decimal.New(user.Balance.IntPart(), int32(-2))

	if *hand.Currency == oldRubleCurrensyCode || *hand.Currency == rubleCurrencyCode {
		newBalance = expBalance
	} else {
		var newCurrency = *hand.Currency
		exchval, err := h.Exchanger.ExchangeRates(h.Logger, expBalance, newCurrency)
		if err != nil {
			if errors.Is(err, exchanger.ErrExchanger) {
				http.Error(w, "incorrect currency code value", http.StatusBadRequest)
				return
			}

			http.Error(w, "cannot convert the value to the specified currency", http.StatusInternalServerError)
			return
		}
		newBalance = exchval
	}

	result := generated.ReadUserResponse{
		Result: struct {
			Balance decimal.Decimal "json:\"balance\""
			UserId  int             "json:\"user_id\""
		}{
			Balance: newBalance,
			UserId:  int(user.AccountID),
		},
		Status: "ok",
	}

	marshalledRequest, err := json.Marshal(result)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, writeErr := w.Write(marshalledRequest)
	if err != nil {
		h.Logger.Error("failed to write connection", zap.Error(writeErr))
		return
	}
}
