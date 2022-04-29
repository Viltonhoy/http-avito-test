package server

import (
	"encoding/json"
	"http-avito-test/internal/storage"
	"io/ioutil"
	"log"
	"net/http"
)

type Handler struct {
	store *storage.Storage
}

type JsUserInf struct {
	ID      int64
	Balance int64
}

func (h *Handler) ReadUser(w http.ResponseWriter, r *http.Request) {
	var hand *JsUserInf

	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &hand)
	if err != nil {
		http.Error(w, "Missing Fields \"ID\", \"Balance\"", http.StatusBadRequest)
		return
	}

	er := h.store.ReadClient(hand.ID, hand.Balance)
	if er != nil {
		log.Fatal("Error reading client", err.Error())
	}

	js, err := json.Marshal(hand)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Write(js)

}

func (h *Handler) AccountFunding(w http.ResponseWriter, r *http.Request) {
	var hand *JsUserInf

	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &hand)
	if err != nil {
		http.Error(w, "Missing Fields \"ID\", \"Balance\"", http.StatusBadRequest)
		return
	}

	er := h.store.UpdateClient(hand.ID, hand.Balance)
	if er != nil {
		log.Fatal("Error updating client", err.Error())
		return
	}
}
