package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"http-avito-test/internal/generated"
	"http-avito-test/internal/storage"
	"io/ioutil"
	"net/http"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

func (h *Handler) RevenueRecognition(w http.ResponseWriter, r *http.Request) {
	var hand *generated.RevenueRecognitionRequest

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

	var newSum = decimal.NewFromFloat32(hand.Sum).Mul(decimal.NewFromInt(100))

	switch {
	case newSum.Exponent() < -2:
		http.Error(w, "wrong value of \"Price\"", http.StatusBadRequest)
		return
	case newSum.LessThanOrEqual(decimal.NewFromInt(int64(0))):
		http.Error(w, "wrong value of \"Price\"", http.StatusBadRequest)
		return
	}

	var description = fmt.Sprintf(`Order number %d; Transferring money for the service %d from a reserve account to a company account in the sum %f`, hand.OrderId, hand.ServiceId, hand.Sum)

	err = h.Store.Revenue(r.Context(), hand.UserId, hand.ServiceId, hand.OrderId, newSum, &description)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrSerialization):
			http.Error(w, "error updating balance", http.StatusInternalServerError)
			return
		case errors.Is(err, storage.ErrTransfer):
			http.Error(w, "not enough money in the account", http.StatusBadRequest)
			return
		case errors.Is(err, storage.ErrUserAvailability):
			http.Error(w, "sender does not exist", http.StatusBadRequest)
			return
		case errors.Is(err, storage.ErrReserveExist):
			http.Error(w, "thе order does not exist", http.StatusBadRequest)
			return
		case errors.Is(err, storage.ErrRevenue):
			http.Error(w, "the sum greater then amount on the deferred expenses", http.StatusBadRequest)
			return
		case errors.Is(err, storage.ErrRecordExist):
			http.Error(w, "unreserve or consolidated report record already exists", http.StatusBadRequest)
			return
		default:
			http.Error(w, "recognition error", http.StatusInternalServerError)
			return
		}
	}

	result := generated.RevenueRecognitionResponse{
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
