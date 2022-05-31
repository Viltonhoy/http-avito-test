package exchanger

import (
	"encoding/json"
	"fmt"
	"http-avito-test/internal/storage"
	"io/ioutil"
	"net/http"

	"github.com/shopspring/decimal"
)

type ExchangeResult struct {
	Result float32
}

func (e *ExchangeResult) ExchangeRates(value decimal.Decimal, currency string) (decimal.Decimal, error) {
	var exch *ExchangeResult
	var conf = storage.NewExch()

	url := fmt.Sprintf(`https://api.apilayer.com/exchangerates_data/convert?to=%s&from=RUB&amount=%s`, currency, value)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("apikey", conf.Key)

	if err != nil {
		fmt.Println(err)
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
