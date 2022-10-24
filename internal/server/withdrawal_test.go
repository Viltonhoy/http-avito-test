package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"http-avito-test/internal/generated"
	"http-avito-test/internal/storage"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestAccountWithdrawal(t *testing.T) {
	t.Run("green case", func(t *testing.T) {
		var testWithdrawal = generated.AccountWithdrawalResponse{
			Result: struct {
				Message string "json:\"message\""
			}{
				Message: "balance updated successfully",
			},
			Status: "ok",
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		description := "test"

		m := NewMockStorager(ctrl)
		m.EXPECT().Withdrawal(gomock.Any(), int64(2), decimal.NewFromFloat32(100).Mul(decimal.NewFromInt(100)), &description).Return(nil)

		arg := bytes.NewBuffer([]byte(`{"User_id":2, "Amount":100.00, "Description":"test"}`))
		req := httptest.NewRequest(http.MethodPost, "http://localhost:9090/withdrawal", arg)
		w := httptest.NewRecorder()

		s := Handler{
			Store: m,
		}

		s.AccountWithdrawal(w, req)

		resp := w.Result()
		body, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err)

		js, err := json.Marshal(testWithdrawal)
		assert.NoError(t, err)

		assert.Equal(t, string(js), string(body))
	})

	t.Run("empty request body", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := NewMockStorager(ctrl)

		req := httptest.NewRequest(http.MethodPost, "http://localhost:9090/withdrawal", nil)
		w := httptest.NewRecorder()

		s := Handler{
			Store: m,
		}

		s.AccountWithdrawal(w, req)

		body, err := ioutil.ReadAll(w.Body)
		assert.NoError(t, err)

		assert.Equal(t, "malformed request body\n", string(body))
	})

	t.Run("wrong UserID value", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := NewMockStorager(ctrl)

		arg := bytes.NewBuffer([]byte(`{"User_id":0, "Amount":100.00}`))

		req := httptest.NewRequest(http.MethodPost, "http://localhost:9090/withdrawal", arg)
		w := httptest.NewRecorder()

		s := Handler{
			Store: m,
		}

		s.AccountWithdrawal(w, req)

		body, err := ioutil.ReadAll(w.Body)
		assert.NoError(t, err)

		assert.Equal(t, "wrong value of \"User_id\"\n", string(body))
	})

	t.Run("wrong amount value", func(t *testing.T) {
		t.Run("amount exponent greater than 2", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := NewMockStorager(ctrl)

			arg := bytes.NewBuffer([]byte(`{"User_id":2, "Amount":100.345}`))
			req := httptest.NewRequest(http.MethodPost, "http://localhost:9090/withdrawal", arg)
			w := httptest.NewRecorder()

			s := Handler{
				Store: m,
			}

			s.AccountWithdrawal(w, req)

			resp := w.Result()
			body, err := ioutil.ReadAll(resp.Body)
			assert.NoError(t, err)

			result := "wrong value of \"Amount\"\n"

			assert.Equal(t, result, string(body))
		})

		t.Run("amount less than or equal to zero", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := NewMockStorager(ctrl)

			arg := bytes.NewBuffer([]byte(`{"User_id":2, "Amount":0}`))
			req := httptest.NewRequest(http.MethodPost, "http://localhost:9090/withdrawwal", arg)
			w := httptest.NewRecorder()

			s := Handler{
				Store: m,
			}

			s.AccountWithdrawal(w, req)

			resp := w.Result()
			body, err := ioutil.ReadAll(resp.Body)
			assert.NoError(t, err)

			result := "wrong value of \"Amount\"\n"

			assert.Equal(t, result, string(body))
		})
	})

	t.Run("withdrawal errors", func(t *testing.T) {
		t.Run("not enough money in the account", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			description := "test"

			m := NewMockStorager(ctrl)
			m.EXPECT().Withdrawal(gomock.Any(), int64(2), decimal.NewFromFloat32(100).Mul(decimal.NewFromInt(100)), &description).Return(storage.ErrWithdrawal)

			arg := bytes.NewBuffer([]byte(`{"User_id":2, "Amount":100.00, "Description":"test"}`))
			req := httptest.NewRequest(http.MethodPost, "http://localhost:9090/withdrawal", arg)
			w := httptest.NewRecorder()

			s := Handler{
				Store: m,
			}

			s.AccountWithdrawal(w, req)

			resp := w.Result()
			body, err := ioutil.ReadAll(resp.Body)
			assert.NoError(t, err)

			assert.Equal(t, "not enough money in the account\n", string(body))
		})

		t.Run("user does not exist", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			description := "test"

			m := NewMockStorager(ctrl)
			m.EXPECT().Withdrawal(gomock.Any(), int64(2), decimal.NewFromFloat32(100).Mul(decimal.NewFromInt(100)), &description).Return(storage.ErrUserAvailability)

			arg := bytes.NewBuffer([]byte(`{"User_id":2, "Amount":100.00, "Description":"test"}`))
			req := httptest.NewRequest(http.MethodPost, "http://localhost:9090/withdrawal", arg)
			w := httptest.NewRecorder()

			s := Handler{
				Store: m,
			}

			s.AccountWithdrawal(w, req)

			resp := w.Result()
			body, err := ioutil.ReadAll(resp.Body)
			assert.NoError(t, err)

			assert.Equal(t, "user does not exist\n", string(body))
		})

		t.Run("error updating balance", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			err := errors.New("error updating balance")

			description := "test"

			m := NewMockStorager(ctrl)
			m.EXPECT().Withdrawal(gomock.Any(), int64(2), decimal.NewFromFloat32(100).Mul(decimal.NewFromInt(100)), &description).Return(err)

			arg := bytes.NewBuffer([]byte(`{"User_id":2, "Amount":100.00, "Description":"test"}`))
			req := httptest.NewRequest(http.MethodPost, "http://localhost:9090/withdrawal", arg)
			w := httptest.NewRecorder()

			s := Handler{
				Store: m,
			}

			s.AccountWithdrawal(w, req)

			resp := w.Result()
			body, err := ioutil.ReadAll(resp.Body)
			assert.NoError(t, err)

			assert.Equal(t, "error updating balance\n", string(body))
		})
		t.Run("level isolation error", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			err := storage.ErrSerialization

			description := "test"

			m := NewMockStorager(ctrl)
			m.EXPECT().Withdrawal(gomock.Any(), int64(2), decimal.NewFromFloat32(100).Mul(decimal.NewFromInt(100)), &description).Return(err)

			arg := bytes.NewBuffer([]byte(`{"User_id":2, "Amount":100.00, "Description":"test"}`))
			req := httptest.NewRequest(http.MethodPost, "http://localhost:9090/withdrawal", arg)
			w := httptest.NewRecorder()

			s := Handler{
				Store: m,
			}

			s.AccountWithdrawal(w, req)

			resp := w.Result()
			body, err := ioutil.ReadAll(resp.Body)
			assert.NoError(t, err)

			assert.Equal(t, "error updating balance\n", string(body))
		})
	})

}
