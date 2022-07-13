package server

import (
	"encoding/json"
	"errors"
	"http-avito-test/internal/generated"
	"http-avito-test/internal/storage"
	"io/ioutil"
	"net/http"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

func (h *Handler) TransferCommand(w http.ResponseWriter, r *http.Request) {
	var hand *generated.TransferCommandRequest

	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &hand)
	if err != nil {
		http.Error(w, "malformed request body", http.StatusBadRequest)
		return
	}

	switch {
	case hand.Sender <= 0:
		http.Error(w, "wrong value of \"Sender\"", http.StatusBadRequest)
		return
	case hand.Recipient <= 0:
		http.Error(w, "wrong value of \"Recipient\"", http.StatusBadRequest)
		return
	}

	var newBalance = decimal.NewFromFloat32(hand.Amount).Mul(decimal.NewFromInt(100))

	switch {
	case newBalance.Exponent() < -2:
		http.Error(w, "wrong value of \"Amount\"", http.StatusBadRequest)
		return
	case newBalance.LessThanOrEqual(decimal.NewFromInt(int64(0))):
		http.Error(w, "wrong value of \"Amount\"", http.StatusBadRequest)
		return
	}

	if hand.Description == nil || *hand.Description == "" {
		hand.Description = nil
	}

	err = h.Store.Transfer(r.Context(), int64(hand.Sender), int64(hand.Recipient), newBalance, hand.Description)
	if err != nil {
		if errors.Is(err, storage.ErrTransfer) {
			http.Error(w, "not enough money in the account", http.StatusBadRequest)
			return
		}
		if errors.Is(err, storage.ErrUserAvailability) {
			http.Error(w, "sender does not exist", http.StatusBadRequest)
			return
		}
		http.Error(w, "error updating balance", http.StatusInternalServerError)
		return
	}

	result := generated.TransferCommandResponse{
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
