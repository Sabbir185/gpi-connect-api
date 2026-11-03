package redis

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

const pingTimeout = 5 * time.Second

func Connect(redisUrl string) (*redis.Client, error) {
	opts, err := parseOptions(redisUrl)
	if err != nil {
		return nil, fmt.Errorf("invalid redis url %q: %w", redisUrl, err)
	}

	opts.DialTimeout = pingTimeout
	opts.ReadTimeout = 3 * time.Second
	opts.WriteTimeout = 3 * time.Second

	client := redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), pingTimeout)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		client.Close()
		return nil, fmt.Errorf("redis.Ping at %s: %w", opts.Addr, err)
	}

	return client, nil
}

func parseOptions(raw string) (*redis.Options, error) {
	raw = strings.TrimSpace(raw)
	if strings.HasPrefix(raw, "redis://") || strings.HasPrefix(raw, "rediss://") || strings.HasPrefix(raw, "unix://") {
		return redis.ParseURL(raw)
	}
	if strings.Contains(raw, "@") {
		return redis.ParseURL("redis://" + raw)
	}
	return &redis.Options{
		Addr: raw,
	}, nil
}
