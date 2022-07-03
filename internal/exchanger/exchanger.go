package exchanger

import (
	"encoding/json"
	"errors"
	"fmt"
	"http-avito-test/internal/generated"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

const envApiKey = "API_KEY"

type ExchangerClient struct {
	apiKey string
	Client *http.Client
}

const ErrorExchangerMessage = "You have entered an invalid \"to\" property. [Example: to=GBP]"

var ErrExchanger = errors.New(ErrorExchangerMessage)

func New() *ExchangerClient {
	return &ExchangerClient{
		apiKey: os.Getenv(envApiKey),
		Client: &http.Client{
			Timeout: 5 * time.Second},
	}
}

func (e *ExchangerClient) ExchangeRates(logger *zap.Logger, value decimal.Decimal, currency string) (decimal.Decimal, error) {
	logger.Debug("starting exchanger rates")

	var ex *generated.ExchangerResult

	url := fmt.Sprintf(`https://api.apilayer.com/exchangerates_data/convert?to=%s&from=RUB&amount=%s`, currency, value)

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("apikey", e.apiKey)

	if err != nil {
		logger.Error("bad request error", zap.Error(err))
		return decimal.NewFromInt(0), err
	}
	res, err := e.Client.Do(req)
	if err != nil {
		return decimal.NewFromInt(0), err
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return decimal.NewFromInt(0), err
	}
	err = json.Unmarshal(body, &ex)
	if err != nil {
		return decimal.NewFromInt(0), err
	}
	if ex.Error != nil {
		logger.Error(ex.Error.Code, zap.Error(ErrExchanger))
		return decimal.NewFromInt(0), ErrExchanger
	}
	return decimal.NewFromFloat32(ex.Result), nil
}
