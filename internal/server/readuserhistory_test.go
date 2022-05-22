package server

import (
	"bytes"
	"context"
	"encoding/json"
	"http-avito-test/internal/storage"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestReadUserHostory(t *testing.T) {
	t.Run("green case", func(t *testing.T) {
		var testHistoryList = []storage.Transf{
			{
				AcountID:  1,
				CBjournal: "Deposit",
				Amount:    100,
				Date:      time.Date(2022, time.May, 05, 1, 0, 0, 0, time.UTC),
				Addressee: "",
			},
			{
				AcountID:  1,
				CBjournal: "Deposit",
				Amount:    120,
				Date:      time.Date(2022, time.May, 05, 2, 0, 0, 0, time.UTC),
				Addressee: "",
			},
			{
				AcountID:  1,
				CBjournal: "Deposit",
				Amount:    130,
				Date:      time.Date(2022, time.May, 05, 3, 0, 0, 0, time.UTC),
				Addressee: "",
			},
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := NewMockStorager(ctrl)
		m.EXPECT().ReadUserHistoryList(int64(1), "amount", context.Background()).Return([]storage.Transf{
			{
				AcountID:  1,
				CBjournal: "Deposit",
				Amount:    100,
				Date:      time.Date(2022, time.May, 05, 1, 0, 0, 0, time.UTC),
				Addressee: "",
			},
			{
				AcountID:  1,
				CBjournal: "Deposit",
				Amount:    120,
				Date:      time.Date(2022, time.May, 05, 2, 0, 0, 0, time.UTC),
				Addressee: "",
			},
			{
				AcountID:  1,
				CBjournal: "Deposit",
				Amount:    130,
				Date:      time.Date(2022, time.May, 05, 3, 0, 0, 0, time.UTC),
				Addressee: "",
			},
		}, nil)

		arg := bytes.NewBuffer([]byte(`{"User_id":1, "Sort": "amount"}`))

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

		assert.Equal(t, js, body)

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

		assert.Equal(t, "Empty request body\n", string(body))
	})
}
