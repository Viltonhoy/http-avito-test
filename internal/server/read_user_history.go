package server

import (
	"encoding/json"
	"errors"
	"http-avito-test/internal/generated"
	"http-avito-test/internal/storage"
	"io/ioutil"
	"net/http"

	"go.uber.org/zap"
)

func (h *Handler) ReadUserHistory(w http.ResponseWriter, r *http.Request) {
	var hand *generated.ReadUserHistoryRequest

	body, _ := ioutil.ReadAll(r.Body)

	err := json.Unmarshal(body, &hand)
	if err != nil {

		if errors.Is(err, storage.ErrBadOrderType) {
			http.Error(w, "wrong value of \"ordBy\" type", http.StatusBadRequest)
			return
		} else {
			http.Error(w, "malformed request body", http.StatusBadRequest)
			return
		}
	}

	if hand.UserId <= 0 {
		http.Error(w, "wrong value of \"User_id\"", http.StatusBadRequest)
		return
	}

	if hand.Order == "" {
		http.Error(w, "wrong value of \"ordBy\" type", http.StatusBadRequest)
		return
	}

	user, err := h.Store.ReadUserHistoryList(r.Context(), hand.UserId, hand.Order, hand.Limit, hand.Offset)
	if err != nil {
		if errors.Is(err, storage.ErrNoUser) {
			http.Error(w, "user does not exist", http.StatusBadRequest)
			return
		}
		http.Error(w, "error reading user history", http.StatusInternalServerError)
		return
	}

	if user == nil {
		http.Error(w, "wrong \"Offset\" value", http.StatusBadRequest)
		return
	}

	result := generated.ReadUserHistoryResponse{
		Result: user,
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
