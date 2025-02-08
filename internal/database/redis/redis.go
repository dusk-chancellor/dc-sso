package redis

import (
	"context"

	"github.com/dusk-chancellor/dc-sso/internal/config"
	
	"github.com/redis/go-redis/v9"
)

// redis client connection setup

func NewClient(ctx context.Context, cfg *config.Redis) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: cfg.Host + ":" + cfg.Port,
		DB: 0, // default
	})
	
	// ping pong
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return client, nil
}
