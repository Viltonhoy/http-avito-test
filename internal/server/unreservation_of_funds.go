package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"http-avito-test/internal/generated"
	"http-avito-test/internal/storage"
	"io/ioutil"
	"net/http"

	"go.uber.org/zap"
)

func (h *Handler) UnreservationOfFunds(w http.ResponseWriter, r *http.Request) {
	var hand *generated.UnreservationOfFundsRequest

	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &hand)
	if err != nil {
		http.Error(w, "malformed request body", http.StatusBadRequest)
		return
	}

	switch {
	case hand.UserId <= 1:
		http.Error(w, "wrong value of \"UserId\"", http.StatusBadRequest)
		return
	case hand.ServiceId <= 0:
		http.Error(w, "wrong value of \"ServiceId\"", http.StatusBadRequest)
		return
	case hand.OrderId <= 0:
		http.Error(w, "wrong value of \"OrderId\"", http.StatusBadRequest)
		return
	}

	var description = fmt.Sprintf(`Order number %d; Refund for the service %d by user %d`, hand.OrderId, hand.ServiceId, hand.UserId)

	err = h.Store.Unreservation(r.Context(), hand.UserId, hand.ServiceId, hand.OrderId, &description)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrSerialization):
			http.Error(w, "error updating balance", http.StatusInternalServerError)
			return
		case errors.Is(err, storage.ErrTransfer):
			http.Error(w, "not enough money in the reserve account", http.StatusInternalServerError)
			return
		case errors.Is(err, storage.ErrReserveExist):
			http.Error(w, "the reserve order does not exist", http.StatusBadRequest)
			return
		case errors.Is(err, storage.ErrRecordExist):
			http.Error(w, "unreserve or consolidated report record already exists", http.StatusBadRequest)
			return
		default:
			http.Error(w, "unreservation error", http.StatusInternalServerError)
			return
		}
	}

	result := generated.UnreservationOfFundsResponse{
		Result: struct {
			Message string "json:\"message\""
		}{
			Message: ResultMessage,
		},
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
