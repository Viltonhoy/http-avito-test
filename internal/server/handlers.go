package server

import (
	"encoding/json"
	"http-avito-test/internal/exchanger"
	"http-avito-test/internal/storage"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/shopspring/decimal"
)

type Handler struct {
	Store *storage.Storage
}

type jsReaderInf struct {
	User_id  int64
	Currency string
}

type returnStruct struct {
	ID      int64
	Balance string
}

type jsFundingInf struct {
	User_id int64
	Balance float32
}

type jsTransferInf struct {
	ID_1    int64
	ID_2    int64
	Balance float32
}

func (h *Handler) ReadUser(w http.ResponseWriter, r *http.Request) {
	var hand *jsReaderInf

	body, _ := ioutil.ReadAll(r.Body)
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

	readUser := returnStruct{
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

func (h *Handler) TransferCommand(w http.ResponseWriter, r *http.Request) {
	var hand *jsTransferInf

	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &hand)
	if err != nil {
		http.Error(w, "Empty request body", http.StatusBadRequest)
		return
	}

	var newBalance = decimal.NewFromFloat32(hand.Balance).Mul(decimal.NewFromInt(100))

	err = h.Store.MoneyTransfer(hand.ID_1, hand.ID_2, newBalance)
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
