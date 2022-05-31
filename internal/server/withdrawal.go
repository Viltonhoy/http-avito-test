package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/shopspring/decimal"
)

type jsWithdrInf struct {
	UserID int64
	Amount float32
}

func (h *Handler) AccountWithdrawal(w http.ResponseWriter, r *http.Request) {
	var hand *jsWithdrInf

	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &hand)
	if err != nil {
		http.Error(w, "Empty request body", http.StatusBadRequest)
		//http.Error(w, "Missing Fields \"ID\", \"Balance\"", http.StatusBadRequest)
		return
	}

	var newBalance = decimal.NewFromFloat32(hand.Amount).Mul(decimal.NewFromInt(100))

	if newBalance.Exponent() < -2 {
		http.Error(w, "wrong value of amount", http.StatusBadRequest)
		return
	}

	err = h.Store.Withdrawal(r.Context(), hand.UserID, newBalance)
	if err != nil {
		log.Fatal("Error updating client", err.Error())
		return
	}

	js, err := json.Marshal("Balance updateted successfully!")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Write(js)
}
