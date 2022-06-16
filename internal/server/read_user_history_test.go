package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"http-avito-test/internal/storage"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestReadUserHostory(t *testing.T) {
	t.Run("green case", func(t *testing.T) {
		var testHistoryList = ReadUserHistoryResponse{
			Result: []storage.ReadUserHistoryResult{
				{
					AccountID:   1,
					CBjournal:   "deposit",
					Amount:      decimal.NewFromInt(100),
					Date:        time.Date(2022, time.May, 05, 1, 0, 0, 0, time.UTC),
					Addressee:   nil,
					Description: nil,
				},
				{
					AccountID:   1,
					CBjournal:   "deposit",
					Amount:      decimal.NewFromInt(120),
					Date:        time.Date(2022, time.May, 05, 2, 0, 0, 0, time.UTC),
					Addressee:   nil,
					Description: nil,
				},
				{
					AccountID:   1,
					CBjournal:   "deposit",
					Amount:      decimal.NewFromInt(130),
					Date:        time.Date(2022, time.May, 05, 3, 0, 0, 0, time.UTC),
					Addressee:   nil,
					Description: nil,
				},
			},
			Status: "ok",
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := NewMockStorager(ctrl)
		m.EXPECT().ReadUserHistoryList(context.Background(), int64(1), "amount", int64(100), int64(0)).Return([]storage.ReadUserHistoryResult{
			{
				AccountID:   1,
				CBjournal:   "deposit",
				Amount:      decimal.NewFromInt(100),
				Date:        time.Date(2022, time.May, 05, 1, 0, 0, 0, time.UTC),
				Addressee:   nil,
				Description: nil,
			},
			{
				AccountID:   1,
				CBjournal:   "deposit",
				Amount:      decimal.NewFromInt(120),
				Date:        time.Date(2022, time.May, 05, 2, 0, 0, 0, time.UTC),
				Addressee:   nil,
				Description: nil,
			},
			{
				AccountID:   1,
				CBjournal:   "deposit",
				Amount:      decimal.NewFromInt(130),
				Date:        time.Date(2022, time.May, 05, 3, 0, 0, 0, time.UTC),
				Addressee:   nil,
				Description: nil,
			},
		},
			nil,
		)

		arg := bytes.NewBuffer([]byte(`{"UserID":1, "Order": "amount", "Limit":100, "Offset":0}`))

		req := httptest.NewRequest(http.MethodPost, "http://localhost:9090/history", arg)
		w := httptest.NewRecorder()

		s := Handler{
			Store: m,
		}

		s.ReadUserHistory(w, req)

		resp := w.Result()
		body, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err)

		js, err := json.Marshal(testHistoryList)
		assert.NoError(t, err)

		assert.Equal(t, string(js), string(body))

	})

	t.Run("empty request body", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		ctrl.Finish()

		m := NewMockStorager(ctrl)

		req := httptest.NewRequest(http.MethodPost, "http://localhost:9090/history", nil)
		w := httptest.NewRecorder()

		s := Handler{
			Store: m,
		}

		s.ReadUserHistory(w, req)

		body, err := ioutil.ReadAll(w.Body)
		assert.NoError(t, err)

		assert.Equal(t, "malformed request body\n", string(body))
	})

	t.Run("wrong UserID value", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		ctrl.Finish()

		m := NewMockStorager(ctrl)

		arg := bytes.NewBuffer([]byte(`{"UserID":0, "Order": "amount", "Limit":100, "Offset":0}`))

		req := httptest.NewRequest(http.MethodPost, "http://localhost:9090/history", arg)
		w := httptest.NewRecorder()

		s := Handler{
			Store: m,
		}

		s.ReadUserHistory(w, req)

		body, err := ioutil.ReadAll(w.Body)
		assert.NoError(t, err)

		assert.Equal(t, "wrong value of \"User_id\"\n", string(body))
	})

	t.Run("green case", func(t *testing.T) {

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		err := errors.New("can not read user history")

		m := NewMockStorager(ctrl)
		m.EXPECT().ReadUserHistoryList(context.Background(), int64(1), "amount", int64(100), int64(0)).Return(nil, err)

		arg := bytes.NewBuffer([]byte(`{"UserID":1, "Order": "amount", "Limit":100, "Offset":0}`))

		req := httptest.NewRequest(http.MethodPost, "http://localhost:9090/history", arg)
		w := httptest.NewRecorder()

		s := Handler{
			Store: m,
		}

		s.ReadUserHistory(w, req)

		resp := w.Result()
		body, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err)

		result := "can not read user history\n"

		assert.Equal(t, result, string(body))

	})
}
