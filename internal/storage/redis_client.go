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

func NewRedisClient(addr string) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
		// Password: "", // TODO: add password
		// DB: 0,
	})

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
