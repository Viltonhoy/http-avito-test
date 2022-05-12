package server

import (
	"bytes"
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
	t.Run("", func(t *testing.T) {
		dt := time.Now()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := NewMockStorager(ctrl)
		m.EXPECT().ReadUserHistoryList(int64(1), "amount").Return([]storage.Transf{
			{
				ID:        1,
				Type:      "Deposit",
				Sum:       100,
				Date:      time.Date(2022, time.May, 05, 0, 0, 0, 0, time.UTC),
				Addressee: "",
			},
			{
				ID:        1,
				Type:      "Deposit",
				Sum:       120,
				Date:      time.Date(2022, time.May, 05, 0, 0, 0, 0, time.UTC),
				Addressee: "",
			},
			{
				ID:        1,
				Type:      "Deposit",
				Sum:       130,
				Date:      time.Date(2022, time.May, 05, 0, 0, 0, 0, time.UTC),
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
		resptest := `[{"ID":1,"Type":"Deposit","Sum":100,
		"Date":"2022-05-08T00:00:00Z","Addressee":""},
		{"ID":1,"Type":"Deposit","Sum":120,
		"Date":"2022-05-08T00:00:00Z","Addressee":""},
		{"ID":1,"Type":"Deposit","Sum":130,
		"Date":"2022-05-08T00:00:00Z","Addressee":""}]`

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)

		assert.Equal(t, string(resptest), string(body))

	})
}
