package db

import (
	"context"
	"go-admin/internal/config"

	"github.com/go-redis/redis/v8"
)

func OpenRedis(ctx context.Context, cfg config.Redis) (*redis.Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if err := redisClient.Ping(ctx).Err(); err != nil {
		_ = redisClient.Close()
		return nil, err
	}
	return redisClient, nil

}
