package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/shopspring/decimal"
)

type jsDeposInf struct {
	UserID int64
	Amount float32
}

func (h *Handler) AccountDeposit(w http.ResponseWriter, r *http.Request) {
	var hand *jsDeposInf

	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &hand)
	if err != nil {
		http.Error(w, "malformed request body", http.StatusBadRequest)
		return
	}

	if hand.UserID <= 0 {
		http.Error(w, "wrong value of \"User_id\"", http.StatusBadRequest)
		return
	}

	var newBalance = decimal.NewFromFloat32(hand.Amount).Mul(decimal.NewFromInt(100))

	if newBalance.Exponent() < -2 {
		http.Error(w, "wrong value of amount", http.StatusBadRequest)
		return
	}

	err = h.Store.Deposit(r.Context(), hand.UserID, newBalance)
	if err != nil {
		// log.Fatal("Error updating client", err.Error())

		return
	}

	js, err := json.Marshal("Balance updateted successfully!")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Write(js)
}
