package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/shopspring/decimal"
)

type jsTransferInf struct {
	ID_1   int64
	ID_2   int64
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

	err = h.Store.Transfer(hand.ID_1, hand.ID_2, newBalance)
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
