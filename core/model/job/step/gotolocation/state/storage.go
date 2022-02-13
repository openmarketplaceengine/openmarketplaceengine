package state

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/openmarketplaceengine/openmarketplaceengine/redis/util"
)

type Storage struct {
	client redis.Client
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

	return nil
}

func (s *Storage) Retrieve(ctx context.Context, key string) (*State, error) {
	m, err := s.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("retrieve HGetAll error: %s", err)
	}

	var state State
	err = util.MapToStruct(m, &state)
	if err != nil {
		return nil, fmt.Errorf("retrieve MapToStruct error: %s", err)
	}
	return &state, nil
}
