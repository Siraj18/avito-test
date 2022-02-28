package currency

import (
	"context"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/siraj18/avito-test/pkg/currencyapi"
)

type Currency struct {
	rdb      *redis.Client
	apiToken string
	lifeTime time.Duration
}

func NewCurrency(rdb *redis.Client, apiToken string, lifeTime time.Duration) *Currency {
	return &Currency{
		rdb:      rdb,
		apiToken: apiToken,
		lifeTime: lifeTime,
	}
}

func (cur *Currency) GetCurrency(currency string) (float64, error) {

	value, err := cur.rdb.Get(context.Background(), currency).Result()
	if err != nil {
		if err == redis.Nil {
			rate, err := currencyapi.GetCurrencyRate("RUB", currency, cur.apiToken)
			if err != nil {
				return 0, err
			}

			err = cur.rdb.Set(context.Background(), currency, rate, cur.lifeTime).Err()
			if err != nil {
				return 0, err
			}

			return rate, nil
		}
		return 0, err
	}

	rate, _ := strconv.ParseFloat(value, 64)

	return rate, nil
}
