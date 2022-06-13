package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

func (h *Handler) TransferCommand(w http.ResponseWriter, r *http.Request) {
	var hand *TransferCommandRequest

	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &hand)
	if err != nil {
		http.Error(w, "Empty request body", http.StatusBadRequest)
		return
	}

	if hand.UserID1 <= 0 && hand.UserID2 <= 0 {
		http.Error(w, "wrong value of \"User_id\"", http.StatusBadRequest)
		return
	}

	var newBalance = decimal.NewFromFloat32(hand.Amount).Mul(decimal.NewFromInt(100))

	switch {
	case newBalance.Exponent() < -2:
		http.Error(w, "wrong value of amount", http.StatusBadRequest)
		return
	case newBalance.LessThanOrEqual(decimal.NewFromInt(int64(0))):
		http.Error(w, "wrong value of amount", http.StatusBadRequest)
		return
	}

	if *hand.Description == "" {
		hand.Description = nil
	}

	err = h.Store.Transfer(r.Context(), hand.UserID1, hand.UserID2, newBalance, hand.Description)
	if err != nil {
		http.Error(w, "Error updating balance", http.StatusInternalServerError)
		return
	}

	result := TransferCommandResponse{
		Result: struct {
			Message string
		}{
			Message: resultMessage,
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
