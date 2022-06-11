package server

import (
	"encoding/json"
	"http-avito-test/internal/exchanger"
	"io/ioutil"
	"net/http"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type jsReaderInf struct {
	UserID   int64
	Currency string
}

type returnReader struct {
	UserID  int64
	Balance decimal.Decimal
}

func (h *Handler) ReadUser(w http.ResponseWriter, r *http.Request) {
	var hand *jsReaderInf
	var exch *exchanger.ExchangeResult

	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	err := json.Unmarshal(body, &hand)
	if err != nil {
		http.Error(w, "malformed request body", http.StatusBadRequest)
		return
	}

	if hand.UserID <= 0 {
		http.Error(w, "wrong value of \"User_id\"", http.StatusBadRequest)
		return
	}

	user, err := h.Store.ReadUser(r.Context(), hand.UserID)
	if err != nil {
		h.Logger.Error("error reading client", zap.Error(err))
		http.Error(w, "error reading client", http.StatusBadRequest)
		return
	}

	var newval decimal.Decimal

	nextval := decimal.New(user.Balance.IntPart(), int32(-2))

	if hand.Currency == "RUR" || hand.Currency == "RUB" || hand.Currency == "" {
		newval = nextval
	} else {
		exchval, _ := exch.ExchangeRates(nextval, hand.Currency)
		newval = exchval
	}

	readUser := returnReader{
		UserID:  user.AccountID,
		Balance: newval,
	}

	js, err := json.Marshal(readUser)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Write(js)
}
