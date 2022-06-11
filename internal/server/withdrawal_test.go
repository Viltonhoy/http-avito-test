package server

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestAccountWithdrawal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := NewMockStorager(ctrl)
	m.EXPECT().Withdrawal(int64(1), decimal.NewFromFloat32(100).Mul(decimal.NewFromInt(100)), "", gomock.Any()).Return(nil)

	arg := bytes.NewBuffer([]byte(`{"User_id":1, "Amount":100.00}`))
	req := httptest.NewRequest(http.MethodPost, "http://localhost:9090/withdrawal", arg)
	w := httptest.NewRecorder()

	s := Handler{
		Store: m,
	}

	s.AccountWithdrawal(w, req)
	resptest := "\"Balance updateted successfully!\""
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, resptest, string(body))
}
