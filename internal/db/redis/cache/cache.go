package cache

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

var ErrCacheMiss = errors.New("cache miss")

type Redis interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, payload []byte, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}

type userCache struct {
	inner *redis.Client
}

func New(client *redis.Client) (Redis, error) {
	return &userCache{
		inner: client,
	}, client.Ping(context.Background()).Err()
}

func (c *userCache) Get(ctx context.Context, userKey string) ([]byte, error) {
	val, err := c.inner.Get(ctx, userKey).Bytes()
	if err == redis.Nil {
		return nil, ErrCacheMiss // your own sentinel error
	}
	return val, err
}

func (c *userCache) Set(ctx context.Context, userKey string, payload []byte, ttl time.Duration) error {
	return c.inner.Set(ctx, userKey, payload, ttl).Err()
}

func (c *userCache) Delete(ctx context.Context, userKey string) error {
	return c.inner.Del(ctx, userKey).Err()
}

