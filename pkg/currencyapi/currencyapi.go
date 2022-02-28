package currencyapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var ErrorWrongCurrency = fmt.Errorf("wrong currency")

func GetCurrencyRate(basicCurrency, currency, apiToken string) (float64, error) {
	apiUrl := fmt.Sprintf("http://api.exchangeratesapi.io/v1/latest?access_key=%s", apiToken)

	resp, err := http.Get(apiUrl)
	if err != nil {
		return 0, err
	}

	var body struct {
		Rates map[string]float64 `json:"rates"`
	}

	if err = json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	rate, ok := body.Rates[currency]
	if !ok {
		return 0, ErrorWrongCurrency
	}

	basicRate, ok := body.Rates[basicCurrency]
	if !ok {
		return 0, ErrorWrongCurrency
	}

	rate = basicRate / rate

	return rate, nil
}
