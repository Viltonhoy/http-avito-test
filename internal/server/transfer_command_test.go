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

func TestTransferCommand(t *testing.T) {
	t.Run("green case", func(t *testing.T) {
		var testTransfer = generated.TransferCommandResponse{
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
		m.EXPECT().Transfer(gomock.Any(), int64(1), int64(2), decimal.NewFromFloat32(100).Mul(decimal.NewFromInt(100)), &description).Return(nil)

		arg := bytes.NewBuffer([]byte(`{"Sender":1, "Recipient":2, "Amount":100.00, "Description":"test"}`))
		req := httptest.NewRequest(http.MethodPost, "http://localhost:9090/transf", arg)
		w := httptest.NewRecorder()

		s := Handler{
			Store: m,
		}

		s.TransferCommand(w, req)

		resp := w.Result()
		body, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err)

		js, err := json.Marshal(testTransfer)
		assert.NoError(t, err)

		assert.Equal(t, string(js), string(body))
	})

	t.Run("malformed request body", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := NewMockStorager(ctrl)

		req := httptest.NewRequest(http.MethodPost, "http://localhost:9090/transf", nil)
		w := httptest.NewRecorder()

		s := Handler{
			Store: m,
		}

		s.TransferCommand(w, req)

		body, err := ioutil.ReadAll(w.Body)
		assert.NoError(t, err)

		assert.Equal(t, "malformed request body\n", string(body))
	})

	t.Run("wrong incoming values", func(t *testing.T) {
		t.Run("wrong sender value", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := NewMockStorager(ctrl)

			arg := bytes.NewBuffer([]byte(`{"Sender":0, "Recipient":2, "Amount":100.00, "Description":"test"}`))

			req := httptest.NewRequest(http.MethodPost, "http://localhost:9090/transf", arg)
			w := httptest.NewRecorder()

			s := Handler{
				Store: m,
			}

			s.TransferCommand(w, req)

			body, err := ioutil.ReadAll(w.Body)
			assert.NoError(t, err)

			assert.Equal(t, "wrong value of \"Sender\"\n", string(body))
		})

		t.Run("wrong recipient value", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := NewMockStorager(ctrl)

			arg := bytes.NewBuffer([]byte(`{"Sender":1, "Recipient":0, "Amount":100.00, "Description":"test"}`))

			req := httptest.NewRequest(http.MethodPost, "http://localhost:9090/transf", arg)
			w := httptest.NewRecorder()

			s := Handler{
				Store: m,
			}

			s.TransferCommand(w, req)

			body, err := ioutil.ReadAll(w.Body)
			assert.NoError(t, err)

			assert.Equal(t, "wrong value of \"Recipient\"\n", string(body))
		})
	})

	t.Run("wrong value of amount", func(t *testing.T) {
		t.Run("amount exponent greater than 2", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := NewMockStorager(ctrl)

			arg := bytes.NewBuffer([]byte(`{"Sender":1, "Recipient":2, "Amount":100.111, "Description":"test"}`))

			req := httptest.NewRequest(http.MethodPost, "http://localhost:9090/transf", arg)
			w := httptest.NewRecorder()

			s := Handler{
				Store: m,
			}

			s.TransferCommand(w, req)

			body, err := ioutil.ReadAll(w.Body)
			assert.NoError(t, err)

			assert.Equal(t, "wrong value of \"Amount\"\n", string(body))
		})

		t.Run("amount less than or equal to zero", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := NewMockStorager(ctrl)

			arg := bytes.NewBuffer([]byte(`{"Sender":1, "Recipient":2, "Amount":-100.00, "Description":"test"}`))

			req := httptest.NewRequest(http.MethodPost, "http://localhost:9090/transf", arg)
			w := httptest.NewRecorder()

			s := Handler{
				Store: m,
			}

			s.TransferCommand(w, req)

			body, err := ioutil.ReadAll(w.Body)
			assert.NoError(t, err)

			assert.Equal(t, "wrong value of \"Amount\"\n", string(body))
		})
	})

	t.Run("transfer errors", func(t *testing.T) {
		t.Run("not enough money in the account", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			description := "test"

			m := NewMockStorager(ctrl)
			m.EXPECT().Transfer(gomock.Any(), int64(1), int64(2), decimal.NewFromFloat32(100).Mul(decimal.NewFromInt(100)), &description).Return(storage.ErrTransfer)

			arg := bytes.NewBuffer([]byte(`{"Sender":1, "Recipient":2, "Amount":100.00, "Description":"test"}`))
			req := httptest.NewRequest(http.MethodPost, "http://localhost:9090/transf", arg)
			w := httptest.NewRecorder()

			s := Handler{
				Store: m,
			}

			s.TransferCommand(w, req)

			resp := w.Result()
			body, err := ioutil.ReadAll(resp.Body)
			assert.NoError(t, err)

			assert.Equal(t, "not enough money in the account\n", string(body))
		})

		t.Run("sender does not exist", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			description := "test"

			m := NewMockStorager(ctrl)
			m.EXPECT().Transfer(gomock.Any(), int64(1000000), int64(2), decimal.NewFromFloat32(100).Mul(decimal.NewFromInt(100)), &description).Return(storage.ErrUserAvailability)

			arg := bytes.NewBuffer([]byte(`{"Sender":1000000, "Recipient":2, "Amount":100.00, "Description":"test"}`))
			req := httptest.NewRequest(http.MethodPost, "http://localhost:9090/transf", arg)
			w := httptest.NewRecorder()

			s := Handler{
				Store: m,
			}

			s.TransferCommand(w, req)

			resp := w.Result()
			body, err := ioutil.ReadAll(resp.Body)
			assert.NoError(t, err)

			assert.Equal(t, "sender does not exist\n", string(body))
		})

		t.Run("error updating balance", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			description := "test"

			m := NewMockStorager(ctrl)
			m.EXPECT().Transfer(gomock.Any(), int64(1000000), int64(2), decimal.NewFromFloat32(100).Mul(decimal.NewFromInt(100)), &description).Return(errors.New(""))

			arg := bytes.NewBuffer([]byte(`{"Sender":1000000, "Recipient":2, "Amount":100.00, "Description":"test"}`))
			req := httptest.NewRequest(http.MethodPost, "http://localhost:9090/transf", arg)
			w := httptest.NewRecorder()

			s := Handler{
				Store: m,
			}

			s.TransferCommand(w, req)

			resp := w.Result()
			body, err := ioutil.ReadAll(resp.Body)
			assert.NoError(t, err)

			assert.Equal(t, "error updating balance\n", string(body))
		})
		t.Run("isolation level error", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			description := "test"

			err := storage.ErrSerialization

			m := NewMockStorager(ctrl)
			m.EXPECT().Transfer(gomock.Any(), int64(1000000), int64(2), decimal.NewFromFloat32(100).Mul(decimal.NewFromInt(100)), &description).Return(err)

			arg := bytes.NewBuffer([]byte(`{"Sender":1000000, "Recipient":2, "Amount":100.00, "Description":"test"}`))
			req := httptest.NewRequest(http.MethodPost, "http://localhost:9090/transf", arg)
			w := httptest.NewRecorder()

			s := Handler{
				Store: m,
			}

			s.TransferCommand(w, req)

			resp := w.Result()
			body, err := ioutil.ReadAll(resp.Body)
			assert.NoError(t, err)

			assert.Equal(t, "error updating balance\n", string(body))
		})
	})
}
