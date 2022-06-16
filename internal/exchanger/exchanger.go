package exchanger

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/caarlos0/env/v6"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type ExchangerClient struct {
	Client *http.Client
}

func New() *ExchangerClient {
	return &ExchangerClient{Client: &http.Client{}}
}

type ExchangeResult struct {
	Result float32         `json:"result"`
	Err    *codeAndMassage `json:"error"`
}

type codeAndMassage struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ExchangerConfig struct {
	Key string `env:"API_KEY"`
}

const ErrorExchangerMessage = "You have entered an invalid \"to\" property. [Example: to=GBP]"

var ErrExchanger = errors.New(ErrorExchangerMessage)

func (e *ExchangerClient) ExchangeRates(logger *zap.Logger, value decimal.Decimal, currency string) (decimal.Decimal, error) {
	logger.Debug("starting exchanger rates")

	logger.Sync()

	var ex *ExchangeResult

	cfg := ExchangerConfig{}
	if err := env.Parse(&cfg); err != nil {
		logger.Error("parse error", zap.Error(err))
		return decimal.NewFromInt(0), err
	}

	url := fmt.Sprintf(`https://api.apilayer.com/exchangerates_data/convert?to=%s&from=RUB&amount=%s`, currency, value)

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("apikey", cfg.Key)

	if err != nil {
		logger.Error("bad request error", zap.Error(err))
		return decimal.NewFromInt(0), err
	}
	res, _ := e.Client.Do(req)
	if res.Body != nil {
		defer res.Body.Close()
	}
	body, _ := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(body, &ex)
	if err != nil {
		return decimal.NewFromInt(0), err
	}
	if ex.Err != nil {
		logger.Error(ex.Err.Code, zap.Error(ErrExchanger))
		return decimal.NewFromInt(0), ErrExchanger
	}
	return decimal.NewFromFloat32(ex.Result), nil
}
