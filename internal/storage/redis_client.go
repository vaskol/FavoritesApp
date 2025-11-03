package storage

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *redis.Client
	TTL    time.Duration
}

// NewRedisClient creates a redis client and sets a default TTL for cached keys.
func NewRedisClient(addr string) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
		// Password: "", // add if needed
		// DB: 0,
	})

	// Optionally: ping to fail early (omitted error handling here to keep minimal)
	// _ = rdb.Ping(context.Background()).Err()

	return &RedisClient{
		Client: rdb,
		TTL:    5 * time.Minute,
	}
}

func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return r.Client.Get(ctx, key).Result()
}

func (r *RedisClient) Set(ctx context.Context, key string, value string) error {
	return r.Client.Set(ctx, key, value, r.TTL).Err()
}

func (r *RedisClient) Del(ctx context.Context, key string) error {
	return r.Client.Del(ctx, key).Err()
}
