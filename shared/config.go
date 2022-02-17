package shared

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/nicklvsa/shorturl/shared/db"
)

type Config struct {
	DB *redis.Client
}

func InitConfig() (*Config, error) {
	// initialize our redis db
	db, err := db.InitRedis(context.Background())
	if err != nil {
		return nil, err
	}

	return &Config{
		DB: db,
	}, nil
}
