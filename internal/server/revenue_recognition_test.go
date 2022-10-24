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

func TestRevenueRecognition(t *testing.T) {
	t.Run("green case", func(t *testing.T) {
		var testReservation = generated.ReservationOfFundsResponse{
			Result: struct {
				Message string "json:\"message\""
			}{
				Message: "balance updated successfully",
			},
			Status: "ok",
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		description := "Order number 1; Transferring money for the service 1 from a reserve account to a company account in the sum 100.000000"

		m := NewMockStorager(ctrl)
		m.EXPECT().Revenue(gomock.Any(), int64(2), int64(1), int64(1), decimal.NewFromFloat32(100).Mul(decimal.NewFromInt(100)), &description).Return(nil)

		arg := bytes.NewBuffer([]byte(`{"user_id":2, "service_id":1, "order_id":1, "sum":100.00}`))
		req := httptest.NewRequest(http.MethodPost, "http://localhost:9090/revenue", arg)
		w := httptest.NewRecorder()

		h := Handler{
			Store: m,
		}

		h.RevenueRecognition(w, req)

		resp := w.Result()
		body, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err)

		js, err := json.Marshal(testReservation)
		assert.NoError(t, err)

		assert.Equal(t, string(js), string(body))
	})

	t.Run("wrong incoming values", func(t *testing.T) {
		t.Run("wrong user_id", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := NewMockStorager(ctrl)

			arg := bytes.NewBuffer([]byte(`{"user_id":0, "service_id":1, "order_id":1, "sum":100.00}`))

			req := httptest.NewRequest(http.MethodPost, "http://localhost:9090/revenue", arg)
			w := httptest.NewRecorder()

			s := Handler{
				Store: m,
			}

			s.RevenueRecognition(w, req)

			body, err := ioutil.ReadAll(w.Body)
			assert.NoError(t, err)

			assert.Equal(t, "wrong value of \"UserId\"\n", string(body))
		})

		t.Run("wrong service_id", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := NewMockStorager(ctrl)

			arg := bytes.NewBuffer([]byte(`{"user_id":2, "service_id":0, "order_id":1, "sum":100.00}`))

			req := httptest.NewRequest(http.MethodPost, "http://localhost:9090/revenue", arg)
			w := httptest.NewRecorder()

			s := Handler{
				Store: m,
			}

			s.RevenueRecognition(w, req)

			body, err := ioutil.ReadAll(w.Body)
			assert.NoError(t, err)

			assert.Equal(t, "wrong value of \"ServiceId\"\n", string(body))
		})

		t.Run("wrong order_id", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := NewMockStorager(ctrl)

			arg := bytes.NewBuffer([]byte(`{"user_id":2, "service_id":1, "order_id":0, "sum":100.00}`))

			req := httptest.NewRequest(http.MethodPost, "http://localhost:9090/revenue", arg)
			w := httptest.NewRecorder()

			s := Handler{
				Store: m,
			}

			s.RevenueRecognition(w, req)

			body, err := ioutil.ReadAll(w.Body)
			assert.NoError(t, err)

			assert.Equal(t, "wrong value of \"OrderId\"\n", string(body))
		})
	})

	t.Run("wrong price value", func(t *testing.T) {
		t.Run("price exponent greater than 2", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := NewMockStorager(ctrl)

			arg := bytes.NewBuffer([]byte(`{"user_id":2, "service_id":1, "order_id":1, "sum":100.111}`))

			req := httptest.NewRequest(http.MethodPost, "http://localhost:9090/revenue", arg)
			w := httptest.NewRecorder()

			s := Handler{
				Store: m,
			}

			s.RevenueRecognition(w, req)

			body, err := ioutil.ReadAll(w.Body)
			assert.NoError(t, err)

			assert.Equal(t, "wrong value of \"Price\"\n", string(body))
		})

		t.Run("price less than or equal to zero", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := NewMockStorager(ctrl)

			arg := bytes.NewBuffer([]byte(`{"user_id":2, "service_id":1, "order_id":1, "sum":0}`))

			req := httptest.NewRequest(http.MethodPost, "http://localhost:9090/revenue", arg)
			w := httptest.NewRecorder()

			s := Handler{
				Store: m,
			}

			s.RevenueRecognition(w, req)

			body, err := ioutil.ReadAll(w.Body)
			assert.NoError(t, err)

			assert.Equal(t, "wrong value of \"Price\"\n", string(body))
		})

		t.Run("revenue recognition errors", func(t *testing.T) {
			t.Run("isolation level error", func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				description := "Order number 1; Transferring money for the service 1 from a reserve account to a company account in the sum 100.000000"

				err := storage.ErrSerialization

				m := NewMockStorager(ctrl)
				m.EXPECT().Revenue(gomock.Any(), int64(2), int64(1), int64(1), decimal.NewFromFloat32(100).Mul(decimal.NewFromInt(100)), &description).Return(err)

				arg := bytes.NewBuffer([]byte(`{"user_id":2, "service_id":1, "order_id":1, "sum":100.00}`))
				req := httptest.NewRequest(http.MethodPost, "http://localhost:9090/revenue", arg)
				w := httptest.NewRecorder()

				s := Handler{
					Store: m,
				}

				s.RevenueRecognition(w, req)

				resp := w.Result()
				body, err := ioutil.ReadAll(resp.Body)
				assert.NoError(t, err)

				assert.Equal(t, "error updating balance\n", string(body))
			})

			t.Run("not enough money in the account", func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				description := "Order number 1; Transferring money for the service 1 from a reserve account to a company account in the sum 100.000000"

				m := NewMockStorager(ctrl)
				m.EXPECT().Revenue(gomock.Any(), int64(2), int64(1), int64(1), decimal.NewFromFloat32(100).Mul(decimal.NewFromInt(100)), &description).Return(storage.ErrTransfer)
				arg := bytes.NewBuffer([]byte(`{"user_id":2, "service_id":1, "order_id":1, "sum":100.00}`))
				req := httptest.NewRequest(http.MethodPost, "http://localhost:9090/revenue", arg)
				w := httptest.NewRecorder()

				s := Handler{
					Store: m,
				}

				s.RevenueRecognition(w, req)

				resp := w.Result()
				body, err := ioutil.ReadAll(resp.Body)
				assert.NoError(t, err)

				assert.Equal(t, "not enough money in the account\n", string(body))
			})

			t.Run("sender does not exist", func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				description := "Order number 1; Transferring money for the service 1 from a reserve account to a company account in the sum 100.000000"

				m := NewMockStorager(ctrl)
				m.EXPECT().Revenue(gomock.Any(), int64(2), int64(1), int64(1), decimal.NewFromFloat32(100).Mul(decimal.NewFromInt(100)), &description).Return(storage.ErrUserAvailability)
				arg := bytes.NewBuffer([]byte(`{"user_id":2, "service_id":1, "order_id":1, "sum":100.00}`))
				req := httptest.NewRequest(http.MethodPost, "http://localhost:9090/revenue", arg)
				w := httptest.NewRecorder()

				s := Handler{
					Store: m,
				}

				s.RevenueRecognition(w, req)

				resp := w.Result()
				body, err := ioutil.ReadAll(resp.Body)
				assert.NoError(t, err)

				assert.Equal(t, "sender does not exist\n", string(body))
			})

			t.Run("order does not exist", func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				description := "Order number 1; Transferring money for the service 1 from a reserve account to a company account in the sum 100.000000"

				m := NewMockStorager(ctrl)
				m.EXPECT().Revenue(gomock.Any(), int64(2), int64(1), int64(1), decimal.NewFromFloat32(100).Mul(decimal.NewFromInt(100)), &description).Return(storage.ErrReserveExist)
				arg := bytes.NewBuffer([]byte(`{"user_id":2, "service_id":1, "order_id":1, "sum":100.00}`))
				req := httptest.NewRequest(http.MethodPost, "http://localhost:9090/revenue", arg)
				w := httptest.NewRecorder()

				s := Handler{
					Store: m,
				}

				s.RevenueRecognition(w, req)

				resp := w.Result()
				body, err := ioutil.ReadAll(resp.Body)
				assert.NoError(t, err)

				assert.Equal(t, "thе order does not exist\n", string(body))
			})

			t.Run("amount on the deferred expenses error", func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				description := "Order number 1; Transferring money for the service 1 from a reserve account to a company account in the sum 100.000000"

				m := NewMockStorager(ctrl)
				m.EXPECT().Revenue(gomock.Any(), int64(2), int64(1), int64(1), decimal.NewFromFloat32(100).Mul(decimal.NewFromInt(100)), &description).Return(storage.ErrRevenue)
				arg := bytes.NewBuffer([]byte(`{"user_id":2, "service_id":1, "order_id":1, "sum":100.00}`))
				req := httptest.NewRequest(http.MethodPost, "http://localhost:9090/revenue", arg)
				w := httptest.NewRecorder()

				s := Handler{
					Store: m,
				}

				s.RevenueRecognition(w, req)

				resp := w.Result()
				body, err := ioutil.ReadAll(resp.Body)
				assert.NoError(t, err)

				assert.Equal(t, "the sum greater then amount on the deferred expenses\n", string(body))
			})

			t.Run("unreserve or consolidated report record error", func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				description := "Order number 1; Transferring money for the service 1 from a reserve account to a company account in the sum 100.000000"

				m := NewMockStorager(ctrl)
				m.EXPECT().Revenue(gomock.Any(), int64(2), int64(1), int64(1), decimal.NewFromFloat32(100).Mul(decimal.NewFromInt(100)), &description).Return(storage.ErrRecordExist)
				arg := bytes.NewBuffer([]byte(`{"user_id":2, "service_id":1, "order_id":1, "sum":100.00}`))
				req := httptest.NewRequest(http.MethodPost, "http://localhost:9090/revenue", arg)
				w := httptest.NewRecorder()

				s := Handler{
					Store: m,
				}

				s.RevenueRecognition(w, req)

				resp := w.Result()
				body, err := ioutil.ReadAll(resp.Body)
				assert.NoError(t, err)

				assert.Equal(t, "unreserve or consolidated report record already exists\n", string(body))
			})

			t.Run("revenue recognition error", func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				description := "Order number 1; Transferring money for the service 1 from a reserve account to a company account in the sum 100.000000"

				m := NewMockStorager(ctrl)
				m.EXPECT().Revenue(gomock.Any(), int64(2), int64(1), int64(1), decimal.NewFromFloat32(100).Mul(decimal.NewFromInt(100)), &description).Return(errors.New(""))
				arg := bytes.NewBuffer([]byte(`{"user_id":2, "service_id":1, "order_id":1, "sum":100.00}`))
				req := httptest.NewRequest(http.MethodPost, "http://localhost:9090/revenue", arg)
				w := httptest.NewRecorder()

				s := Handler{
					Store: m,
				}

				s.RevenueRecognition(w, req)

				resp := w.Result()
				body, err := ioutil.ReadAll(resp.Body)
				assert.NoError(t, err)

				assert.Equal(t, "recognition error\n", string(body))
			})
		})
	})
}
