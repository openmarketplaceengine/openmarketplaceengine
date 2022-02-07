package client

import (
	"github.com/go-redis/redis/v8"
	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
)

func NewStoreClient() *redis.Client {
	addr := cfg.Redis.Store.Addr
	password := cfg.Redis.Store.Pass

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: string(password),
	})
	return client
}

func NewStoreClientPool() *redis.Client {
	pool := cfg.Redis.Store.Pool
	addr := cfg.Redis.Store.Addr
	password := cfg.Redis.Store.Pass

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: string(password),
		PoolSize: pool,
	})
	return client
}

func NewPubSubClient() *redis.Client {
	addr := cfg.Redis.Pubsub.Addr
	password := cfg.Redis.Pubsub.Pass

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: string(password),
	})
	return client
}

func NewPubSubClientPool() *redis.Client {
	pool := cfg.Redis.Pubsub.Pool
	addr := cfg.Redis.Pubsub.Addr
	password := cfg.Redis.Pubsub.Pass

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: string(password),
		PoolSize: pool,
	})
	return client
}
