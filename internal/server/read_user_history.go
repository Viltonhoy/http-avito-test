package server

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"go.uber.org/zap"
)

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
	case string(orderByAmount):
		*j = orderByAmount
	case string(orderByDate):
		*j = orderByDate
	default:
		return errBadOrderType
	}

	return nil
}

func (h *Handler) ReadUserHistory(w http.ResponseWriter, r *http.Request) {
	var hand *ReadUserHistoryRequest

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

	if hand.UserID <= 0 {
		http.Error(w, "wrong value of \"User_id\"", http.StatusBadRequest)
		return
	}

	hist, err := h.Store.ReadUserHistoryList(r.Context(), hand.UserID, string(hand.Order), hand.Limit, hand.Offset)
	if err != nil {
		http.Error(w, "can not read user history", http.StatusInternalServerError)
		return
	}

	result := ReadUserHistoryResponse{
		Result: hist,
		Status: "ok",
	}

	marshalledRequest, err := json.Marshal(result)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, writeErr := w.Write(marshalledRequest)
	if err != nil {
		h.Logger.Error("failed to write connection", zap.Error(writeErr))
		return
	}
}