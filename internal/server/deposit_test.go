package server

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAccountDeposit(t *testing.T) {
	ctrl := gomock.NewController(t)

	m := NewMockStorager(ctrl)
	m.EXPECT().Deposit(int64(1), float32(100.00), "Deposit").Return()

	arg := bytes.NewBuffer([]byte(`{"User_id":1, "Amount":100.00}`))
	req := httptest.NewRequest(http.MethodPost, "http://localhost:9090/deposit", arg)
	w := httptest.NewRecorder()

	s := Handler{
		Store: m,
	}

	s.AccountDeposit(w, req)
	resptest := "Balance updateted successfully!"
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, string(resptest), string(body))
}
