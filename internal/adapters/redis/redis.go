package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/nlsnnn/berezhok/internal/shared/config"
)

func New(ctx context.Context, cfg config.Redis) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       0,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return client, nil
}
