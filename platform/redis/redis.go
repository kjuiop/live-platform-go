package redis

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/kjuiop/live-platform-go/config"
)

type Client struct {
	cfg    config.Redis
	client *redis.Client
}

func NewRedisSingleClient(ctx context.Context, cfg config.Redis) (*Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		PoolSize:     cfg.PoolSize,
		DialTimeout:  time.Second * 3,
		ReadTimeout:  time.Second * 3,
		WriteTimeout: time.Second * 3,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		if cerr := client.Close(); cerr != nil {
			slog.Error("fail close redis client after ping error", "error", cerr)
		}
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	return &Client{
		cfg:    cfg,
		client: client,
	}, nil
}

func (r *Client) Close() {
	if err := r.client.Close(); err != nil {
		slog.Error("fail close redis client", "error", err)
	}
}
