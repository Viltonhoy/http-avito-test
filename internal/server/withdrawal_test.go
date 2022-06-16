package server

// import (
// 	"bytes"
// 	"encoding/json"
// 	"io/ioutil"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/golang/mock/gomock"
// 	"github.com/shopspring/decimal"
// 	"github.com/stretchr/testify/assert"
// )

// func TestAccountWithdrawal(t *testing.T) {
// 	t.Run("green case", func(t *testing.T) {
// 		var testWithdrawal = AccountWithdrawalResponse{
// 			Result: struct {
// 				Message string
// 			}{
// 				Message: "balance updated successfully",
// 			},
// 			Status: "ok",
// 		}

// 		ctrl := gomock.NewController(t)
// 		defer ctrl.Finish()

// 		description := ""

// 		m := NewMockStorager(ctrl)
// 		m.EXPECT().Withdrawal(gomock.Any(), int64(1), decimal.NewFromFloat32(100).Mul(decimal.NewFromInt(100)), description).Return(nil)

// 		arg := bytes.NewBuffer([]byte(`{"UserID":1, "Amount":100.00}`))
// 		req := httptest.NewRequest(http.MethodPost, "http://localhost:9090/withdrawal", arg)
// 		w := httptest.NewRecorder()

// 		s := Handler{
// 			Store: m,
// 		}

// 		s.AccountDeposit(w, req)

// 		resp := w.Result()
// 		body, _ := ioutil.ReadAll(resp.Body)

// 		js, err := json.Marshal(testWithdrawal)
// 		assert.NoError(t, err)

// 		assert.Equal(t, string(js), string(body))
// 	})
// }
