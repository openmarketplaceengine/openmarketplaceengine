package client

import (
	"github.com/go-redis/redis/v8"
	"github.com/openmarketplaceengine/openmarketplaceengine/config"
)

func NewStoreClient() *redis.Client {
	addr := config.GetString(config.RedisStoreAddr)
	password := config.GetString(config.RedisStorePassword)

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
	})
	return client
}

func NewStoreClientPool() *redis.Client {
	pool := config.GetInt(config.RedisStorePool)
	addr := config.GetString(config.RedisStoreAddr)
	password := config.GetString(config.RedisStorePassword)

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		PoolSize: pool,
	})
	return client
}

func NewPubSubClient() *redis.Client {
	addr := config.GetString(config.RedisPubSubAddr)
	password := config.GetString(config.RedisPubSubPassword)

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
	})
	return client
}

func NewPubSubClientPool() *redis.Client {
	pool := config.GetInt(config.RedisPubSubPool)
	addr := config.GetString(config.RedisPubSubAddr)
	password := config.GetString(config.RedisPubSubPassword)

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		PoolSize: pool,
	})
	return client
}
