package server

import (
	"bytes"
	"context"
	"http-avito-test/internal/storage"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestReadUser(t *testing.T) {
	t.Run("green case", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := NewMockStorager(ctrl)

		m.EXPECT().ReadUser(context.Background(), int64(1)).Return(storage.UserBalance{
			AccountID: 1,
			Balance:   decimal.NewFromInt(10000),
		}, nil)

		arg := bytes.NewBuffer([]byte(`{"UserID":1, "Currency":"RUB"}`))
		req := httptest.NewRequest(http.MethodPost, "http://loacalhost:9090/read", arg)
		w := httptest.NewRecorder()

		s := Handler{
			Store: m,
		}

		s.ReadUser(w, req)
		resptest := `{"result":{"userID":1,"balance":"100"},"status":"ok"}`
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

		arg := bytes.NewBuffer([]byte(`{"UserID":1, "Currency":""}`))
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

	// t.Run("user does not exist", func(t *testing.T) {
	// 	ctrl := gomock.NewController(t)
	// 	defer ctrl.Finish()

	// 	err := errors.New("user does not exist\n")

	// 	m := NewMockStorager(ctrl)
	// 	m.EXPECT().ReadUser(context.Background(), int64(1000000)).Return(storage.UserBalance{}, err)

	// 	arg := bytes.NewBuffer([]byte(`{"UserID":1000000, "Currency":"RUB"}`))
	// 	req := httptest.NewRequest(http.MethodPost, "http://loacalhost:9090/read", arg)
	// 	w := httptest.NewRecorder()

	// 	s := Handler{
	// 		Store: m,
	// 	}

	// 	s.ReadUser(w, req)

	// 	resp := w.Result()
	// 	body, err := ioutil.ReadAll(resp.Body)
	// 	assert.NoError(t, err)

	// 	resptest := `user does not exist`
	// 	assert.Equal(t, resptest, string(body))
	// })

}
