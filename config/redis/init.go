package redis

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
)

func (r Redis) Init() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     r.Addr,
		Password: r.Password, // no password set
		DB:       r.DB,       // use default DB
	})
	result, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	if result != "PONG" {
		return nil, errors.New("redis ping error")
	}
	return client, nil
}
