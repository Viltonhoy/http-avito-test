package server

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

type jsHistoryReader struct {
	UserID int64
	Order  ordBy
	Limit  int64
	Offset int64
}

type ordBy string

const (
	orderByAmount ordBy = "amount"
	orderByDate   ordBy = "date"
)

var errBadOrderType = errors.New("wrong value of ordBy type")

func (j *ordBy) UnmarshalJSON(v []byte) error {
	var s string

	if err := json.Unmarshal(v, &s); err != nil {
		return err
	}

	switch s {
	case "amount":
		*j = orderByAmount
	case "date":
		*j = orderByDate
	default:
		return errBadOrderType
	}

	return nil
}

func (h *Handler) ReadUserHistory(w http.ResponseWriter, r *http.Request) {
	var hand *jsHistoryReader

	body, _ := ioutil.ReadAll(r.Body)

	err := json.Unmarshal(body, &hand)
	if err != nil {

		if errors.Is(err, errBadOrderType) {
			http.Error(w, "wrong value of ordBy type", http.StatusBadRequest)
			return
		} else {
			http.Error(w, "malformed request body", http.StatusBadRequest)
			return
		}

	}

	hist, err := h.Store.ReadUserHistoryList(r.Context(), hand.UserID, string(hand.Order), hand.Limit, hand.Offset)
	if err != nil {
		log.Panic("Error reading history", err.Error())
		return
	}

	js, err := json.Marshal(hist)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Write(js)
}
