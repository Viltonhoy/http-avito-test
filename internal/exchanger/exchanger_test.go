package exchanger

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExchangeRates(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		expected := "100.00"
		assert.Equal(t, req.URL.String(), "api.apilayer.com/exchangerates_data/convert?to=RUB&from=RUB&amount=100.00")
		fmt.Fprintf(rw, expected)
	}))

	defer server.Close()

	e := ExchangeResult{}
	body, _ := e.ExchangeRates("100.00", "RUB")

	body = strings.TrimSpace(body)
}
