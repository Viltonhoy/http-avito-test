package server

import (
	"bytes"
	"encoding/json"
	"http-avito-test/internal/generated"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestMonthlyReport(t *testing.T) {
	t.Run("green case", func(t *testing.T) {
		var testReservation = generated.MonthlyReportResponse{
			Result: struct {
				Link string "json:\"link\""
			}{
				Link: reportLink,
			},
			Status: "ok",
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := NewMockStorager(ctrl)
		m.EXPECT().MonthlyReport(gomock.Any(), int64(2022), int64(10)).Return([][]string{}, nil)

		arg := bytes.NewBuffer([]byte(`{"year":2022, "month":10}`))
		req := httptest.NewRequest(http.MethodPost, "http://localhost:9090/report", arg)
		w := httptest.NewRecorder()

		h := Handler{
			Store: m,
		}

		h.MonthlyReport(w, req)

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

		req := httptest.NewRequest(http.MethodPost, "http://localhost:9090/report", nil)
		w := httptest.NewRecorder()

		s := Handler{
			Store: m,
		}

		s.MonthlyReport(w, req)

		body, err := ioutil.ReadAll(w.Body)
		assert.NoError(t, err)

		assert.Equal(t, "malformed request body\n", string(body))
	})

}
