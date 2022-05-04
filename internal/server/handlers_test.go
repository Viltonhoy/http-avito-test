package server

// import (
// 	"bytes"
// 	"http-avito-test/internal/storage"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"
// )

// func TestHandlerReadUser(t *testing.T) {
// 	t.Run("green case", func(t *testing.T){
// 		arg:=bytes.NewBuffer(`{"ID":1}`)
// 		req:=httptest.NewRequest(http.MethodGet, "http://loacalhost:9090/ReadUser", arg)
// 		w:=httptest.NewRecorder()

// 		var s = storage.NewStore()
// 		h:=Handler{
// 			Store: s,
// 		}

// 		h.ReadUser(w, req)
// 		resptest:=`{"ID":1, "Balance":10}`
// 	}
// }
