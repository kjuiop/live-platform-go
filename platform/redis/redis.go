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

func (r *Client) HSet(ctx context.Context, key string, data interface{}, expiration time.Duration) error {

	pipe := r.client.TxPipeline()
	pipe.HSet(ctx, key, data)
	if expiration > 0 {
		pipe.Expire(ctx, key, expiration)
	}

	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("redis pipeline exec error: %w", err)
	}

	return nil
}

func (r *Client) HSetField(ctx context.Context, key, field string, value interface{}) error {
	return r.client.HSet(ctx, key, field, value).Err()
}

func (r *Client) RunScript(ctx context.Context, script *redis.Script, keys []string, args ...interface{}) error {
	return script.Run(ctx, r.client, keys, args...).Err()
}

func (r *Client) RunScriptResult(ctx context.Context, script *redis.Script, keys []string, args ...interface{}) (*redis.Cmd, error) {
	cmd := script.Run(ctx, r.client, keys, args...)
	return cmd, cmd.Err()
}

func (r *Client) Close() {
	if err := r.client.Close(); err != nil {
		slog.Error("fail close redis client", "error", err)
	}
}
