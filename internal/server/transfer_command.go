package server

import (
	"encoding/json"
	"http-avito-test/internal/generated"
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
		http.Error(w, "Empty request body", http.StatusBadRequest)
		return
	}

	if hand.Userid1 <= 0 && hand.Userid2 <= 0 {
		http.Error(w, "wrong value of \"Userid\"", http.StatusBadRequest)
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

	err = h.Store.Transfer(r.Context(), int64(hand.Userid1), int64(hand.Userid2), newBalance, hand.Description)
	if err != nil {
		http.Error(w, "Error updating balance", http.StatusInternalServerError)
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
