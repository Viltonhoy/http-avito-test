package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/shopspring/decimal"
)

type jsDeposInf struct {
	User_id int64
	Amount  float32
}

func (h *Handler) AccountDeposit(w http.ResponseWriter, r *http.Request) {
	var hand *jsDeposInf

	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &hand)
	if err != nil {
		http.Error(w, "Empty request body", http.StatusBadRequest)
		return
	}

	var newBalance = decimal.NewFromFloat32(hand.Amount).Mul(decimal.NewFromInt(100))

	err = h.Store.Deposit(hand.User_id, newBalance, r.Context())
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
