package gotolocation

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
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

func key(driverID string) string {
	return fmt.Sprintf("gotolocation-%s", driverID)
}

func (s *Storage) Store(ctx context.Context, goToLocation GoToLocation) error {
	m, err := marshalutils.StructToMap(goToLocation)

	if err != nil {
		return fmt.Errorf("store StructToMap error: %s", err)
	}

	key := key(goToLocation.DriverID)

	err = s.client.HSet(ctx, key, m).Err()
	if err != nil {
		return fmt.Errorf("store HSet error: %s", err)
	}

	s.client.Expire(ctx, key, s.expiration)

	return nil
}

func (s *Storage) Retrieve(ctx context.Context, driverID string) (goToLocation GoToLocation, err error) {
	key := key(driverID)
	m, err := s.client.HGetAll(ctx, key).Result()
	if err != nil {
		err = fmt.Errorf("retrieve HGetAll error: %s", err)
		return
	}

	if len(m) == 0 {
		err = fmt.Errorf("goToLocation not found by driverID %q", driverID)
		return
	}

	err = marshalutils.MapToStruct(m, &goToLocation)
	if err != nil {
		err = fmt.Errorf("retrieve MapToStruct error: %s", err)
		return
	}
	return
}
