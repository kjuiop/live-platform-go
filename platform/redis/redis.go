package redis

import (
	"context"
	"errors"
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

func (r *Client) HVals(ctx context.Context, key string) ([]string, error) {
	vals, err := r.client.HVals(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("HVals %s: %w", key, err)
	}
	return vals, nil
}

func (r *Client) HGetAllPipeline(ctx context.Context, keys []string) ([]map[string]string, error) {
	if len(keys) == 0 {
		return nil, nil
	}

	pipe := r.client.TxPipeline()
	// command 를 실행할 목록 배열을 저장
	cmds := make([]*redis.MapStringStringCmd, len(keys))
	for i, key := range keys {
		cmds[i] = pipe.HGetAll(ctx, key)
	}

	// command 실행
	if _, err := pipe.Exec(ctx); err != nil && !errors.Is(err, redis.Nil) {
		return nil, fmt.Errorf("HGetAllPipeline exec: %w", err)
	}

	// 결과를 map[string]string 형태로 변환하여 반환
	results := make([]map[string]string, len(keys))
	for i, cmd := range cmds {
		m, err := cmd.Result()
		if err != nil && !errors.Is(err, redis.Nil) {
			return nil, fmt.Errorf("HGetAllPipeline cmd[%d]: %w", i, err)
		}
		results[i] = m
	}
	return results, nil
}

func (r *Client) Close() {
	if err := r.client.Close(); err != nil {
		slog.Error("fail close redis client", "error", err)
	}
}
