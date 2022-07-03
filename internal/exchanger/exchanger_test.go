package exchanger

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

const response = `{
	"date": "2022-06-25",
	"info": {
	  "rate": 0.017477,
	  "timestamp": 1656145743
	},
	"query": {
	  "amount": 100,
	  "from": "RUB",
	  "to": "EUR"
	},
	"result": 1.7477,
	"success": true
  }`

const responseError = `
{
    "error": {
        "code": "invalid_to_currency",
        "message": "You have entered an invalid \"to\" property. [Example: to=GBP]"
    }
}`

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func NewTestClient(statusCode int, response string) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(func(req *http.Request) *http.Response {
			return &http.Response{
				StatusCode: statusCode,
				Body:       ioutil.NopCloser(bytes.NewBufferString(response)),
			}
		}),
	}
}

var newClient = ExchangerClient{
	Client: NewTestClient(http.StatusOK, response),
}

var newClientErr = ExchangerClient{
	Client: NewTestClient(http.StatusOK, responseError),
}

var badRequest = ExchangerClient{
	Client: NewTestClient(http.StatusBadRequest, ""),
}

func TestExchangeRates(t *testing.T) {
	t.Run("green case", func(t *testing.T) {
		var logger, err = zap.NewDevelopment()
		assert.NoError(t, err)
		value, err := newClient.ExchangeRates(logger, decimal.NewFromInt(100), "EUR")
		assert.NoError(t, err)
		result := decimal.NewFromFloat32(1.7477)
		assert.Equal(t, result, value)
	})

	t.Run("wrong currency code", func(t *testing.T) {
		var logger, err = zap.NewDevelopment()
		assert.NoError(t, err)
		_, err = newClientErr.ExchangeRates(logger, decimal.NewFromInt(100), "test")
		result := errors.New("You have entered an invalid \"to\" property. [Example: to=GBP]")
		assert.Equal(t, result, err)
	})
}
