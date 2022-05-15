package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type jsHistoryReader struct {
	User_id int64
	Sort    string
}

func (h *Handler) ReadUserHistory(w http.ResponseWriter, r *http.Request) {
	var hand *jsHistoryReader

	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &hand)
	if err != nil {
		http.Error(w, "Empty request body", http.StatusBadRequest)
		return
	}

	hist, err := h.Store.ReadUserHistoryList(hand.User_id, hand.Sort, r.Context())
	if err != nil {
		log.Fatal("Error reading history", err.Error())
	}

	js, err := json.Marshal(hist)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Write(js)
}
