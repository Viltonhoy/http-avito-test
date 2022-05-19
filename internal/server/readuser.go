package server

import (
	"encoding/json"
	"fmt"
	"http-avito-test/internal/exchanger"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/shopspring/decimal"
)

type jsReaderInf struct {
	User_id  int64
	Currency string
}

type returnReader struct {
	ID      int64
	Balance string
}

func (h *Handler) ReadUser(w http.ResponseWriter, r *http.Request) {
	var hand *jsReaderInf
	var exch *exchanger.ExchangeResult

	r.Context()

	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	err := json.Unmarshal(body, &hand)
	if err != nil {
		fmt.Println("Error", err)
		http.Error(w, "Empty request body", http.StatusBadRequest)
		return
	}

	if hand.User_id <= 0 {
		http.Error(w, "Missing Field \"User_id\"", http.StatusBadRequest)
		return
	}

	user, err := h.Store.ReadClient(hand.User_id, r.Context())
	if err != nil {
		log.Fatal("Error reading client", err.Error())
	}

	var newval string
	val := user.Balance.String()
	nextval := val[:len(val)-2] + "." + val[len(val)-2:]

	if hand.Currency == "RUR" || hand.Currency == "RUB" || hand.Currency == "" {
		newval = nextval
	} else {
		exchval, _ := exch.ExchangeRates(nextval, hand.Currency)
		newval = decimal.NewFromFloat32(exchval).String()
	}

	readUser := returnReader{
		ID:      user.ID,
		Balance: newval,
	}

	js, err := json.Marshal(readUser)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Write(js)
}
