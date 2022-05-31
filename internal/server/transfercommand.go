package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/shopspring/decimal"
)

type jsTransferInf struct {
	ID1    int64
	ID2    int64
	Amount float32
}

func (h *Handler) TransferCommand(w http.ResponseWriter, r *http.Request) {
	var hand *jsTransferInf

	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &hand)
	if err != nil {
		http.Error(w, "Empty request body", http.StatusBadRequest)
		return
	}

	var newBalance = decimal.NewFromFloat32(hand.Amount).Mul(decimal.NewFromInt(100))

	if newBalance.Exponent() < -2 {
		http.Error(w, "wrong value of amount", http.StatusBadRequest)
		return
	}

	err = h.Store.Transfer(r.Context(), hand.ID1, hand.ID2, newBalance)
	if err != nil {
		log.Fatal("Error transfer client", err.Error())
		return
	}

	js, err := json.Marshal("Transfer was successfull!")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Write(js)
}
