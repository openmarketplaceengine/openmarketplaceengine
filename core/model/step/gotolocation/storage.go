package gotolocation

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/openmarketplaceengine/openmarketplaceengine/core/model/step"
	"github.com/openmarketplaceengine/openmarketplaceengine/redis/client"
	"github.com/openmarketplaceengine/openmarketplaceengine/redis/marshalutils"
)

type Storage struct {
	client     *redis.Client
	expiration time.Duration
}

var storage = newStorage(3 * time.Hour)

func newStorage(expiration time.Duration) *Storage {
	return &Storage{
		client:     client.NewStoreClient(),
		expiration: expiration,
	}
}

func key(stepID string) string {
	return fmt.Sprintf("gotolocation-%s", stepID)
}

func (s *Storage) Store(ctx context.Context, goToLocation GoToLocation) error {
	m, err := marshalutils.StructToMap(goToLocation)

	if err != nil {
		return fmt.Errorf("store StructToMap error: %w", err)
	}

	key := key(goToLocation.StepID)

	err = s.client.HSet(ctx, key, m).Err()
	if err != nil {
		return fmt.Errorf("store HSet error: %w", err)
	}

	s.client.Expire(ctx, key, s.expiration)

	return nil
}

func (s *Storage) Retrieve(ctx context.Context, stepID step.ID) (goToLocation *GoToLocation, err error) {
	key := key(string(stepID))
	m, err := s.client.HGetAll(ctx, key).Result()
	if err != nil {
		err = fmt.Errorf("retrieve HGetAll error: %w", err)
		return
	}

	if len(m) == 0 {
		return
	}

	goToLocation = &GoToLocation{}

	err = marshalutils.MapToStruct(m, goToLocation)
	if err != nil {
		err = fmt.Errorf("retrieve MapToStruct error: %w", err)
		return
	}
	return
}
