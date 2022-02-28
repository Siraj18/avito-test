package rediscache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

func NewRedis(addr, pass string, db int) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {

		return nil, err
	}

	return rdb, nil
}
