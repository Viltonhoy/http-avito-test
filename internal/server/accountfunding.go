package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/shopspring/decimal"
)

type jsFundingInf struct {
	User_id int64
	Balance float32
}

func (h *Handler) AccountFunding(w http.ResponseWriter, r *http.Request) {
	var hand *jsFundingInf

	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &hand)
	if err != nil {
		http.Error(w, "Empty request body", http.StatusBadRequest)
		//http.Error(w, "Missing Fields \"ID\", \"Balance\"", http.StatusBadRequest)
		return
	}

	var newBalance = decimal.NewFromFloat32(hand.Balance).Mul(decimal.NewFromInt(100))

	err = h.Store.UpdateClient(hand.User_id, newBalance)
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
