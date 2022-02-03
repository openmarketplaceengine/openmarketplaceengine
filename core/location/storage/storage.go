package storage

import (
	"context"
	"github.com/go-redis/redis/v8"
)

type LocationStorage struct {
	client *redis.Client
}

func NewLocationStorage(client *redis.Client) *LocationStorage {
	s := LocationStorage{
		client: client,
	}
	return &s
}

func (s *LocationStorage) GeoAdd(ctx context.Context, key, name string, longitude, latitude float64) (err error) {
	err = s.client.GeoAdd(ctx, key, &redis.GeoLocation{
		Name:      name,
		Longitude: longitude,
		Latitude:  latitude,
		Dist:      0,
		GeoHash:   0,
	}).Err()
	if err != nil {
		return err
	}

	return nil
}

func (s *LocationStorage) GeoSearch(ctx context.Context, key string, longitude, latitude, radius float64, radiusUnit string) (result []string, err error) {
	result, err = s.client.GeoSearch(ctx, key, &redis.GeoSearchQuery{
		Member:     "",
		Longitude:  longitude,
		Latitude:   latitude,
		Radius:     radius,
		RadiusUnit: radiusUnit,
		BoxWidth:   0,
		BoxHeight:  0,
		BoxUnit:    "",
		Sort:       "",
		Count:      0,
		CountAny:   false,
	}).Result()

	return
}
