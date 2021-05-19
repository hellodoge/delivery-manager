package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"strconv"
)

type RedisConnectionOptions struct {
	Addr string
	Port uint16
	Password string
}

func RedisConnect(options RedisConnectionOptions, db int) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: options.Addr + ":" + strconv.Itoa(int(options.Port)),
		Password: options.Password,
		DB: db,
	})
	err := client.Ping(context.Background()).Err()
	return client, err
}