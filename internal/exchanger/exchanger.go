package exchanger

import (
	"encoding/json"
	"fmt"
	"http-avito-test/internal/storage"
	"io/ioutil"
	"net/http"
)

type ExchangeResult struct {
	Result float32
}

func ExchangeRates(value string, currency string) *ExchangeResult {
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
		return nil
	}
	return exch
}
