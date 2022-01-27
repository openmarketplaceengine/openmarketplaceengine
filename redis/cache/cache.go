package cache

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type Cache interface {
	GetName() string
	Put(ctx context.Context, key string, value string) error
	Get(ctx context.Context, key string) (value string, err error)
	Remove(ctx context.Context, key string) (err error)
	ContainsKey(ctx context.Context, key string) (containsKey bool)
}

func NewCache(c *redis.Client, name string) Cache {
	return &cache{
		c:    c,
		name: getName(name),
	}
}

func getName(name string) string {
	return fmt.Sprintf("cache:%s", name)
}

type cache struct {
	c    *redis.Client
	name string
}

func (c *cache) GetName() string {
	return c.name
}

func (c *cache) Put(ctx context.Context, key string, value string) (err error) {
	return c.c.HSet(ctx, c.name, key, value).Err()
}

func (c *cache) Get(ctx context.Context, key string) (value string, err error) {
	return c.c.HGet(ctx, c.name, key).Result()
}

func (c *cache) Remove(ctx context.Context, key string) (err error) {
	return c.c.HDel(ctx, c.name, key).Err()
}

func (c *cache) ContainsKey(ctx context.Context, key string) (containsKey bool) {
	containsKey, err := c.c.HExists(ctx, c.name, key).Result()
	if err != nil {
		panic(err)
	}
	return
}
