package db

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/nicklvsa/shorturl/shared/logger"
)

func InitRedis(ctx context.Context) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	ping, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	logger.Infof("Successfully initialized redis! Received ping message: %s", ping)
	return client, nil
}
