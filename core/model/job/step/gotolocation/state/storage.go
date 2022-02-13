package state

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/openmarketplaceengine/openmarketplaceengine/redis/client"
	"github.com/openmarketplaceengine/openmarketplaceengine/redis/util"
	"time"
)

type Storage struct {
	client     *redis.Client
	expiration time.Duration
}

func NewStorage(expiration time.Duration) *Storage {
	return &Storage{
		client:     client.NewStoreClient(),
		expiration: expiration,
	}
}

func (s *Storage) Store(ctx context.Context, key string, state State) error {
	m, err := util.StructToMap(state)

	if err != nil {
		return fmt.Errorf("store StructToMap error: %s", err)
	}

	err = s.client.HSet(ctx, key, m).Err()
	if err != nil {
		return fmt.Errorf("store HSet error: %s", err)
	}

	s.client.Expire(ctx, key, s.expiration)

	return nil
}

func (s *Storage) Retrieve(ctx context.Context, key string) (state State, err error) {
	m, err := s.client.HGetAll(ctx, key).Result()
	if err != nil {
		err = fmt.Errorf("retrieve HGetAll error: %s", err)
		return
	}

	err = util.MapToStruct(m, &state)
	if err != nil {
		err = fmt.Errorf("retrieve MapToStruct error: %s", err)
		return
	}
	return
}
