package server

import (
	"bytes"
	"context"
	"errors"
	"http-avito-test/internal/exchanger"
	"http-avito-test/internal/storage"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestReadUser(t *testing.T) {
	t.Run("green case", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := NewMockStorager(ctrl)

		m.EXPECT().ReadUserByID(context.Background(), int64(2)).Return(storage.User{
			AccountID: 2,
			Balance:   decimal.NewFromInt(10000),
		}, nil)

		arg := bytes.NewBuffer([]byte(`{"User_id":2, "Currency":"RUB"}`))
		req := httptest.NewRequest(http.MethodPost, "http://loacalhost:9090/read", arg)
		w := httptest.NewRecorder()

		s := Handler{
			Store: m,
		}

		s.ReadUser(w, req)
		resptest := "{\"result\":{\"balance\":\"100\",\"user_id\":2},\"status\":\"ok\"}"
		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)

		assert.Equal(t, resptest, string(body))
	})

	t.Run("empty request body", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := NewMockStorager(ctrl)

		req := httptest.NewRequest(http.MethodPost, "http://loacalhost:9090/read", nil)
		w := httptest.NewRecorder()

		s := Handler{
			Store: m,
		}

		s.ReadUser(w, req)
		body, err := ioutil.ReadAll(w.Body)
		assert.NoError(t, err)

		assert.Equal(t, "malformed request body\n", string(body))
	})

	t.Run("wrong User_id value", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := NewMockStorager(ctrl)

		arg := bytes.NewBuffer([]byte(`{"UserID":0, "Currency":"RUB"}`))
		req := httptest.NewRequest(http.MethodPost, "http://loacalhost:9090/read", arg)
		w := httptest.NewRecorder()

		s := Handler{
			Store: m,
		}

		s.ReadUser(w, req)
		body, err := ioutil.ReadAll(w.Body)
		assert.NoError(t, err)

		assert.Equal(t, "wrong value of \"User_id\"\n", string(body))
	})

	t.Run("empty Currency value", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := NewMockStorager(ctrl)

		arg := bytes.NewBuffer([]byte(`{"User_id":2, "Currency":""}`))
		req := httptest.NewRequest(http.MethodPost, "http://loacalhost:9090/read", arg)
		w := httptest.NewRecorder()

		s := Handler{
			Store: m,
		}

		s.ReadUser(w, req)
		resptest := "incorrect currency code value\n"
		body, err := ioutil.ReadAll(w.Body)
		assert.NoError(t, err)

		assert.Equal(t, resptest, string(body))
	})

	t.Run("user does not exist", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := NewMockStorager(ctrl)
		m.EXPECT().ReadUserByID(context.Background(), int64(1000000)).Return(
			storage.User{},
			storage.ErrUserAvailability)

		arg := bytes.NewBuffer([]byte(`{"User_id":1000000, "Currency":"RUB"}`))
		req := httptest.NewRequest(http.MethodPost, "http://loacalhost:9090/read", arg)
		w := httptest.NewRecorder()

		s := Handler{
			Store: m,
		}

		s.ReadUser(w, req)

		resp := w.Result()
		body, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err)

		resptest := "user does not exist\n"
		assert.Equal(t, resptest, string(body))
	})

	t.Run("cannot read user with specified id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := NewMockStorager(ctrl)
		m.EXPECT().ReadUserByID(context.Background(), int64(2)).Return(
			storage.User{},
			errors.New("cannot read user with specified id"))

		arg := bytes.NewBuffer([]byte(`{"User_id":2, "Currency":"RUB"}`))
		req := httptest.NewRequest(http.MethodPost, "http://loacalhost:9090/read", arg)
		w := httptest.NewRecorder()

		s := Handler{
			Store: m,
		}

		s.ReadUser(w, req)

		resp := w.Result()
		body, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err)

		resptest := "cannot read user with specified id\n"
		assert.Equal(t, resptest, string(body))
	})

	t.Run("exchanger errors", func(t *testing.T) {
		t.Run("incorrect currency code value", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			var logger *zap.Logger

			newStorage := storage.User{
				AccountID: 2,
				Balance:   decimal.NewFromInt(10000),
			}

			m := NewMockStorager(ctrl)
			m.EXPECT().ReadUserByID(context.Background(), int64(2)).Return(
				newStorage,
				nil)

			e := NewMockExchanger(ctrl)
			e.EXPECT().ExchangeRates(logger, decimal.New(newStorage.Balance.IntPart(), int32(-2)), "RUBBB").Return(decimal.NewFromInt(0),
				exchanger.ErrExchanger)

			arg := bytes.NewBuffer([]byte(`{"User_id":2, "Currency":"RUBBB"}`))
			req := httptest.NewRequest(http.MethodPost, "http://loacalhost:9090/read", arg)
			w := httptest.NewRecorder()

			s := Handler{
				Store:     m,
				Exchanger: e,
			}

			s.ReadUser(w, req)

			resp := w.Result()
			body, err := ioutil.ReadAll(resp.Body)
			assert.NoError(t, err)

			resptest := "incorrect currency code value\n"
			assert.Equal(t, resptest, string(body))
		})

		t.Run("cannot convert the value to the specified currency", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			var logger *zap.Logger

			newStorage := storage.User{
				AccountID: 2,
				Balance:   decimal.NewFromInt(10000),
			}

			m := NewMockStorager(ctrl)
			m.EXPECT().ReadUserByID(context.Background(), int64(2)).Return(
				newStorage,
				nil)

			e := NewMockExchanger(ctrl)
			e.EXPECT().ExchangeRates(logger, decimal.New(newStorage.Balance.IntPart(), int32(-2)), "EUR").Return(decimal.NewFromInt(0),
				errors.New(""))

			arg := bytes.NewBuffer([]byte(`{"User_id":2, "Currency":"EUR"}`))
			req := httptest.NewRequest(http.MethodPost, "http://loacalhost:9090/read", arg)
			w := httptest.NewRecorder()

			s := Handler{
				Store:     m,
				Exchanger: e,
			}

			s.ReadUser(w, req)

			resp := w.Result()
			body, err := ioutil.ReadAll(resp.Body)
			assert.NoError(t, err)

			resptest := "cannot convert the value to the specified currency\n"
			assert.Equal(t, resptest, string(body))
		})
	})
}
