package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/shopspring/decimal"
)

type jsWithdrInf struct {
	User_id int64
	Amount  float32
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

	var newBalance = decimal.NewFromFloat32(hand.Amount).Mul(decimal.NewFromInt(-100))
	var tp = "Withdrawal"

	err = h.Store.DepositOrWithdrawal(hand.User_id, newBalance, tp)
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
