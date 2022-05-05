package server

import (
	"encoding/json"
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
	//vars := mux.Vars(r)
	//key := vars["key"]

	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	err := json.Unmarshal(body, &hand)
	if err != nil {
		http.Error(w, "Empty request body", http.StatusBadRequest)
		return
	}

	user, err := h.Store.ReadClient(hand.User_id)
	if err != nil {
		log.Fatal("Error reading client", err.Error())
	}

	var newval string
	val := user.Balance.String()
	nextval := val[:len(val)-2] + "." + val[len(val)-2:]

	if hand.Currency == "RUR" || hand.Currency == "RUB" || hand.Currency == "" {
		newval = nextval
	} else {
		newBalance := exchanger.ExchangeRates(nextval, hand.Currency)
		newval = decimal.NewFromFloat32(newBalance.Result).String()
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
