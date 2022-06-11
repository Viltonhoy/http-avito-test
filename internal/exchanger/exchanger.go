package exchanger

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/caarlos0/env/v6"
	"github.com/shopspring/decimal"
)

type ExchangeResult struct {
	Result float32
}

type ExchangerConfig struct {
	Key string `env:"API_KEY"`
}

func (e *ExchangeResult) ExchangeRates(value decimal.Decimal, currency string) (decimal.Decimal, error) {
	var exch *ExchangeResult

	cfg := ExchangerConfig{}
	if err := env.Parse(&cfg); err != nil {
		return decimal.NewFromInt(0), err
	}

	url := fmt.Sprintf(`https://api.apilayer.com/exchangerates_data/convert?to=%s&from=RUB&amount=%s`, currency, value)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("apikey", cfg.Key)

	if err != nil {
		return decimal.NewFromInt(0), err
	}
	res, _ := client.Do(req)
	if res.Body != nil {
		defer res.Body.Close()
	}
	body, _ := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(body, &exch)
	if err != nil {
		return decimal.NewFromInt(0), err
	}
	return decimal.NewFromFloat32(exch.Result), nil
}
