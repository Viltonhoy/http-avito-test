package server

import (
	"bytes"
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
	//logger, err := zap.NewDevelopment()
	//require.NoError(t, err)

	t.Run("green case", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := NewMockStorager(ctrl)

		m.EXPECT().ReadClient(int64(1)).Return(storage.User{
			ID:      1,
			Balance: decimal.NewFromInt(10000),
		}, nil)

		arg := bytes.NewBuffer([]byte(`{"User_id":1, "Currency":""}`))
		req := httptest.NewRequest(http.MethodPost, "http://loacalhost:9090/read", arg)
		w := httptest.NewRecorder()

		s := Handler{
			Store: m,
		}

		s.ReadUser(w, req)
		resptest := `{"ID":1,"Balance":"100.00"}`
		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)

		assert.Equal(t, string(resptest), string(body))
	})

	t.Run("empty request body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "http://loacalhost:9090/read", nil)
		w := httptest.NewRecorder()

		s := Handler{
			Store: nil,
		}

		s.ReadUser(w, req)
		body, _ := ioutil.ReadAll(w.Body)

		assert.Equal(t, "Empty request body\n", string(body))
	})

	t.Run("wrong User_id value", func(t *testing.T) {
		arg := bytes.NewBuffer([]byte(`{"User_id":0, "Currency":""}`))
		req := httptest.NewRequest(http.MethodPost, "http://loacalhost:9090/read", arg)
		w := httptest.NewRecorder()

		s := Handler{
			Store: nil,
		}

		s.ReadUser(w, req)
		body, _ := ioutil.ReadAll(w.Body)

		assert.Equal(t, "Missing Field \"User_id\"\n", string(body))
	})

}
