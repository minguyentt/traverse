package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, payload []byte, ttl time.Duration) error
	Delete(ctx context.Context, key string)
}

type userCache struct {
	inner *redis.Client
}

func New(client *redis.Client) (Cache, error) {
	return &userCache{
		inner: client,
	}, client.Ping(context.Background()).Err()
}

func (c *userCache) Get(ctx context.Context, userKey string) ([]byte, error) {
	cmd := c.inner.Get(ctx, userKey)
	return cmd.Bytes()
}

func (c *userCache) Set(ctx context.Context, userKey string, payload []byte, ttl time.Duration) error {
	return c.inner.Set(ctx, userKey, payload, ttl).Err()
}

func (c *userCache) Delete(ctx context.Context, userKey string) {
	c.inner.Del(ctx, userKey)
}
