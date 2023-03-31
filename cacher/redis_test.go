package cacher

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestCacher_Get(t *testing.T) {
	mr, err := miniredis.Run()
	assert.NoError(t, err)

	defer mr.Close()

	cacher := NewCacher(redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
		DB:   0,
	}))

	ctx := context.TODO()

	t.Run("ok", func(t *testing.T) {
		err := mr.Set("key", "value")
		assert.NoError(t, err)

		res, err := cacher.Get(ctx, "key")

		assert.NoError(t, err)
		assert.Equal(t, "value", res)
	})

	t.Run("not found", func(t *testing.T) {
		mr.SetError("redis: nil")
		res, err := cacher.Get(ctx, "not_found")

		assert.Error(t, err)
		assert.Equal(t, "", res)
		assert.Equal(t, err, redis.Nil)

		mr.SetError("")
	})

	t.Run("error", func(t *testing.T) {
		mr.SetError("err redis")

		res, err := cacher.Get(ctx, "key")

		assert.Error(t, err)
		assert.Equal(t, "", res)

		mr.SetError("")
	})
}

func TestCacher_Set(t *testing.T) {
	mr, err := miniredis.Run()
	assert.NoError(t, err)

	defer mr.Close()

	cacher := NewCacher(redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
		DB:   0,
	}))

	ctx := context.TODO()

	t.Run("ok", func(t *testing.T) {
		err := cacher.Set(ctx, "key", "value", 1000)
		assert.NoError(t, err)

		res, err := mr.Get("key")
		assert.NoError(t, err)

		assert.Equal(t, res, "value")
	})

	t.Run("error", func(t *testing.T) {
		mr.SetError("err redis")

		err := cacher.Set(ctx, "key", "value", 1000)
		assert.Error(t, err)

		mr.SetError("")
	})
}
