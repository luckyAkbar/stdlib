// Package cacher contains all cache related functionality
package cacher

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// Cacher :nodoc:
type Cacher interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, exp time.Duration) error
}

type cacher struct {
	client *redis.Client
}

// NewCacher return a new model.Chacher instance
func NewCacher(client *redis.Client) Cacher {
	return &cacher{
		client: client,
	}
}

// Get get cache value by given key. Return json string if found. Otherwise return a non nil error
func (c *cacher) Get(ctx context.Context, key string) (string, error) {
	res, err := c.client.Get(ctx, key).Result()
	switch err {
	case nil:
		return res, nil
	case redis.Nil:
		return res, err
	default:
		return res, err
	}
}

// Set set a cache value by key with the given expiry time. The val should be a json string
func (c *cacher) Set(ctx context.Context, key string, val string, exp time.Duration) error {
	err := c.client.Set(ctx, key, val, exp).Err()
	if err != nil {
		return err
	}

	return nil
}
