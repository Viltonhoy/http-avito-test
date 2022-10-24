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
	"github.com/stretchr/testify/assert"
)

func TestUnreservationOfFunds(t *testing.T) {
	t.Run("green case", func(t *testing.T) {
		var testReservation = generated.UnreservationOfFundsResponse{
			Result: struct {
				Message string "json:\"message\""
			}{
				Message: "balance updated successfully",
			},
			Status: "ok",
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		description := "Order number 1; Refund for the service 1 by user 2"

		m := NewMockStorager(ctrl)
		m.EXPECT().Unreservation(gomock.Any(), int64(2), int64(1), int64(1), &description).Return(nil)

		arg := bytes.NewBuffer([]byte(`{"user_id":2, "service_id":1, "order_id":1}`))
		req := httptest.NewRequest(http.MethodPost, "http://localhost:9090/unreserv", arg)
		w := httptest.NewRecorder()

		h := Handler{
			Store: m,
		}

		h.UnreservationOfFunds(w, req)

		resp := w.Result()
		body, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err)

		js, err := json.Marshal(testReservation)
		assert.NoError(t, err)

		assert.Equal(t, string(js), string(body))
	})

	t.Run("malformed request body", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := NewMockStorager(ctrl)

		req := httptest.NewRequest(http.MethodPost, "http://localhost:9090/unreserv", nil)
		w := httptest.NewRecorder()

		s := Handler{
			Store: m,
		}

		s.UnreservationOfFunds(w, req)

		body, err := ioutil.ReadAll(w.Body)
		assert.NoError(t, err)

		assert.Equal(t, "malformed request body\n", string(body))
	})

	t.Run("wrong incoming values", func(t *testing.T) {
		t.Run("wrong user_id", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := NewMockStorager(ctrl)

			arg := bytes.NewBuffer([]byte(`{"user_id":0, "service_id":1, "order_id":1}`))

			req := httptest.NewRequest(http.MethodPost, "http://localhost:9090/unreserv", arg)
			w := httptest.NewRecorder()

			s := Handler{
				Store: m,
			}

			s.UnreservationOfFunds(w, req)

			body, err := ioutil.ReadAll(w.Body)
			assert.NoError(t, err)

			assert.Equal(t, "wrong value of \"UserId\"\n", string(body))
		})

		t.Run("wrong service_id", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := NewMockStorager(ctrl)

			arg := bytes.NewBuffer([]byte(`{"user_id":2, "service_id":0, "order_id":1}`))

			req := httptest.NewRequest(http.MethodPost, "http://localhost:9090/unreserv", arg)
			w := httptest.NewRecorder()

			s := Handler{
				Store: m,
			}

			s.UnreservationOfFunds(w, req)

			body, err := ioutil.ReadAll(w.Body)
			assert.NoError(t, err)

			assert.Equal(t, "wrong value of \"ServiceId\"\n", string(body))
		})

		t.Run("wrong order_id", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := NewMockStorager(ctrl)

			arg := bytes.NewBuffer([]byte(`{"user_id":2, "service_id":1, "order_id":0}`))

			req := httptest.NewRequest(http.MethodPost, "http://localhost:9090/unreserv", arg)
			w := httptest.NewRecorder()

			s := Handler{
				Store: m,
			}

			s.UnreservationOfFunds(w, req)

			body, err := ioutil.ReadAll(w.Body)
			assert.NoError(t, err)

			assert.Equal(t, "wrong value of \"OrderId\"\n", string(body))
		})
	})

	t.Run("unreservation errors", func(t *testing.T) {
		t.Run("isolation level error", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			description := "Order number 1; Refund for the service 1 by user 2"

			err := storage.ErrSerialization

			m := NewMockStorager(ctrl)
			m.EXPECT().Unreservation(gomock.Any(), int64(2), int64(1), int64(1), &description).Return(err)

			arg := bytes.NewBuffer([]byte(`{"user_id":2, "service_id":1, "order_id":1}`))
			req := httptest.NewRequest(http.MethodPost, "http://localhost:9090/unreserv", arg)
			w := httptest.NewRecorder()

			s := Handler{
				Store: m,
			}

			s.UnreservationOfFunds(w, req)

			resp := w.Result()
			body, err := ioutil.ReadAll(resp.Body)
			assert.NoError(t, err)

			assert.Equal(t, "error updating balance\n", string(body))
		})

		t.Run("not enough money in the account", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			description := "Order number 1; Refund for the service 1 by user 2"

			m := NewMockStorager(ctrl)
			m.EXPECT().Unreservation(gomock.Any(), int64(2), int64(1), int64(1), &description).Return(storage.ErrTransfer)
			arg := bytes.NewBuffer([]byte(`{"user_id":2, "service_id":1, "order_id":1}`))
			req := httptest.NewRequest(http.MethodPost, "http://localhost:9090/unreserv", arg)
			w := httptest.NewRecorder()

			s := Handler{
				Store: m,
			}

			s.UnreservationOfFunds(w, req)

			resp := w.Result()
			body, err := ioutil.ReadAll(resp.Body)
			assert.NoError(t, err)

			assert.Equal(t, "not enough money in the reserve account\n", string(body))
		})

		t.Run("the reserve order error", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			description := "Order number 1; Refund for the service 1 by user 2"

			m := NewMockStorager(ctrl)
			m.EXPECT().Unreservation(gomock.Any(), int64(2), int64(1), int64(1), &description).Return(storage.ErrReserveExist)
			arg := bytes.NewBuffer([]byte(`{"user_id":2, "service_id":1, "order_id":1}`))
			req := httptest.NewRequest(http.MethodPost, "http://localhost:9090/unreserv", arg)
			w := httptest.NewRecorder()

			s := Handler{
				Store: m,
			}

			s.UnreservationOfFunds(w, req)

			resp := w.Result()
			body, err := ioutil.ReadAll(resp.Body)
			assert.NoError(t, err)

			assert.Equal(t, "the reserve order does not exist\n", string(body))
		})

		t.Run("unreservation or consolidated report error", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			description := "Order number 1; Refund for the service 1 by user 2"

			m := NewMockStorager(ctrl)
			m.EXPECT().Unreservation(gomock.Any(), int64(2), int64(1), int64(1), &description).Return(storage.ErrRecordExist)
			arg := bytes.NewBuffer([]byte(`{"user_id":2, "service_id":1, "order_id":1}`))
			req := httptest.NewRequest(http.MethodPost, "http://localhost:9090/unreserv", arg)
			w := httptest.NewRecorder()

			s := Handler{
				Store: m,
			}

			s.UnreservationOfFunds(w, req)

			resp := w.Result()
			body, err := ioutil.ReadAll(resp.Body)
			assert.NoError(t, err)

			assert.Equal(t, "unreserve or consolidated report record already exists\n", string(body))
		})

		t.Run("unreservation error", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			description := "Order number 1; Refund for the service 1 by user 2"

			m := NewMockStorager(ctrl)
			m.EXPECT().Unreservation(gomock.Any(), int64(2), int64(1), int64(1), &description).Return(errors.New(""))
			arg := bytes.NewBuffer([]byte(`{"user_id":2, "service_id":1, "order_id":1}`))
			req := httptest.NewRequest(http.MethodPost, "http://localhost:9090/unreserv", arg)
			w := httptest.NewRecorder()

			s := Handler{
				Store: m,
			}

			s.UnreservationOfFunds(w, req)

			resp := w.Result()
			body, err := ioutil.ReadAll(resp.Body)
			assert.NoError(t, err)

			assert.Equal(t, "unreservation error\n", string(body))
		})
	})
}
