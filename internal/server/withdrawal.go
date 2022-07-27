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

func (h *Handler) AccountWithdrawal(w http.ResponseWriter, r *http.Request) {
	var hand *generated.AccountWithdrawalRequest

	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &hand)
	if err != nil {
		http.Error(w, "malformed request body", http.StatusBadRequest)
		return
	}

	if hand.UserId <= 0 {
		http.Error(w, "wrong value of \"User_id\"", http.StatusBadRequest)
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

	newErr := h.Store.Withdrawal(r.Context(), int64(hand.UserId), newBalance, hand.Description)
	// Обработка 40001 ошибки в хендлере

	// При обновлении данных при конкурентном нагрузочном тестировании происходит ограничение доступа
	// к базе, из-за одновременного обновления. База данных должна обработать первый запрос и записать новый результат в таблицу balances, чтобы
	// следующий запрос выполнялся с новым актуальным значением баланса пользователя.

	// Ошибка: при выполнении ретрая остальные запросы выполняются некорректно, увеличивая баланс пользователя, при этом действует ограничение
	// в виде сериализуемого уровня изоляции транзакции. То есть несколько транзакций выполняются одновременно.

	// count := 0
	// for newErr != nil {
	// 	switch {
	// 	case errors.Is(newErr, storage.ErrSerialization) && count < 5:
	// 		newErr = h.Store.Withdrawal(r.Context(), int64(hand.UserId), newBalance, hand.Description)
	// 		if newErr == storage.ErrSerialization {
	// 			count++
	// 		}
	// 	case count >= 5:
	// 		http.Error(w, "error updating balance", http.StatusTeapot)
	// 		return
	// 	}
	// }

	if newErr != nil {
		if errors.Is(newErr, storage.ErrSerialization) {
			http.Error(w, "error updating balance", http.StatusTeapot)
			return
		}
		if errors.Is(newErr, storage.ErrWithdrawal) {
			http.Error(w, "not enough money in the account", http.StatusBadRequest)
			return
		}
		if errors.Is(newErr, storage.ErrUserAvailability) {
			http.Error(w, "user does not exist", http.StatusBadRequest)
			return
		}
		http.Error(w, "error updating balance", http.StatusInternalServerError)
		return
	}

	result := generated.AccountWithdrawalResponse{
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
