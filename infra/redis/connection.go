package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const pingTimeout = 5 * time.Second

func Connect(redisUrl string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         redisUrl,
		DialTimeout:  pingTimeout,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), pingTimeout)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		client.Close()
		return nil, fmt.Errorf("redis.Ping at %s: %w", redisUrl, err)
	}

	return client, nil
}
