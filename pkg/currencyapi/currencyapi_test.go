package currencyapi_test

import (
	"testing"

	"github.com/siraj18/avito-test/pkg/currencyapi"
	"github.com/stretchr/testify/assert"
)

var correctApiToken = "05802b752dcc15626f922580104bad3a"
var wrongApiToken = "05802b743dcc15626f922580104bad3a"

func TestGetCurrencyRate(t *testing.T) {
	baseCurrency := "RUB"
	finalCurrency := "USD"

	rate, err := currencyapi.GetCurrencyRate(baseCurrency, finalCurrency, correctApiToken)
	assert.Nil(t, err)

	assert.NotEqual(t, rate, 0)
}

func TestGetCurrencyRateWrongBasicCurrency(t *testing.T) {
	baseCurrency := "HER"
	finalCurrency := "USD"

	rate, err := currencyapi.GetCurrencyRate(baseCurrency, finalCurrency, correctApiToken)

	assert.EqualError(t, err, currencyapi.ErrorWrongCurrency.Error())
	assert.Equal(t, rate, 0.0)
}
func TestGetCurrencyRateWrongFinalCurrency(t *testing.T) {
	baseCurrency := "RUB"
	finalCurrency := "HER"

	rate, err := currencyapi.GetCurrencyRate(baseCurrency, finalCurrency, correctApiToken)

	assert.EqualError(t, err, currencyapi.ErrorWrongCurrency.Error())
	assert.Equal(t, rate, 0.0)
}

func TestGetCurrencyRateWrongApiToken(t *testing.T) {
	baseCurrency := "RUB"
	finalCurrency := "USD"

	rate, err := currencyapi.GetCurrencyRate(baseCurrency, finalCurrency, wrongApiToken)

	assert.NotNil(t, err)
	assert.Equal(t, rate, 0.0)
}
