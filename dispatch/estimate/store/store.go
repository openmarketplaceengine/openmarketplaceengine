package store

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/openmarketplaceengine/openmarketplaceengine/dispatch/estimate"
)

const (
	Radius2000m = 2000
	Radius4000m = 4000
	Radius9000m = 9000
)

type Handler func(ctx context.Context, geoHash string, estimates []*estimate.Estimate) error

type EstimateStore struct {
	client *redis.Client
}

func NewEstimateStore(client *redis.Client) *EstimateStore {
	return &EstimateStore{
		client: client,
	}
}

func (s *EstimateStore) GetAll(ctx context.Context, geoHash string, radiusMeters int) ([]*estimate.Estimate, error) {
	k := key(geoHash, radiusMeters)
	estimates, err := s.client.HGetAll(ctx, k).Result()

	if err != nil {
		return nil, fmt.Errorf("hgetall error: %w", err)
	}

	l := len(estimates)
	if l == 0 {
		return nil, nil
	}

	result := make([]*estimate.Estimate, l)

	for k, v := range estimates {
		var e estimate.Estimate
		err := e.UnmarshalBinary([]byte(v))
		if err != nil {
			return nil, fmt.Errorf("unmarshal geohash %s of %s error: %w", geoHash, k, err)
		}
		result = append(result, &e)
	}

	return result, nil
}

func (s *EstimateStore) Store(ctx context.Context, geoHash string, radiusMeters int, estimates []*estimate.Estimate) error {
	values := toValues(estimates)
	k := key(geoHash, radiusMeters)
	err := s.client.HSet(ctx, k, values...).Err()

	if err != nil {
		return fmt.Errorf("hset error: %w", err)
	}

	return nil
}

func key(geoHash string, radiusMeters int) string {
	return fmt.Sprintf("%s-%v", geoHash, radiusMeters)
}

func toValues(estimates []*estimate.Estimate) []interface{} {
	res := make([]interface{}, 0)
	for _, e := range estimates {
		res = append(res, e.ID)
		res = append(res, e)
	}
	return res
}
