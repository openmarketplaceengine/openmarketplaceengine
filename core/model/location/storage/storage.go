package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/openmarketplaceengine/openmarketplaceengine/core/model/location"
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

func (s *LocationStorage) updateLocation(ctx context.Context, areaKey string, location location.UpdateLocation) (err error) {
	err = s.client.GeoAdd(ctx, areaKey, &redis.GeoLocation{
		Name:      location.PrincipalID,
		Longitude: location.Longitude,
		Latitude:  location.Latitude,
		Dist:      0,
		GeoHash:   0,
	}).Err()
	if err != nil {
		return err
	}

	expireTrackingKey := expireTrackingKey(areaKey)
	s.client.ZAddArgs(ctx, expireTrackingKey, redis.ZAddArgs{
		NX: false,
		XX: false,
		LT: false,
		GT: false,
		Ch: false,
		Members: []redis.Z{{
			Score:  float64(time.Now().UnixMilli()),
			Member: location.PrincipalID,
		}},
	})

	return nil
}

func (s *LocationStorage) ForgetLocation(ctx context.Context, areaKey string, principalID string) (err error) {
	err = s.client.ZRem(ctx, areaKey, principalID).Err()
	if err != nil {
		return err
	}
	expireTrackingKey := expireTrackingKey(areaKey)
	err = s.client.ZRem(ctx, expireTrackingKey, principalID).Err()
	if err != nil {
		return err
	}

	return nil
}

func (s *LocationStorage) RemoveExpiredLocations(ctx context.Context, areaKey string, before time.Time) (err error) {
	key := expireTrackingKey(areaKey)
	min := fmt.Sprintf("%v", float64(0))
	max := fmt.Sprintf("%v", float64(before.UnixMilli()))

	keys, err := s.client.ZRangeByScore(ctx, key, &redis.ZRangeBy{
		Min: min,
		Max: max,
	}).Result()

	if err != nil {
		return err
	}

	if len(keys) > 0 {
		innerErr := s.client.ZRemRangeByScore(ctx, key, min, max).Err()
		if innerErr != nil {
			return fmt.Errorf("LocationStorage ZRemRangeByScore error: %s", innerErr)
		}
		innerErr = s.client.ZRem(ctx, areaKey, keys).Err()
		if innerErr != nil {
			return fmt.Errorf("LocationStorage HDel error: %s", innerErr)
		}
	}
	return nil
}

func (s *LocationStorage) QueryLocations(ctx context.Context, areaKey string, longitude float64, latitude float64, radius float64, radiusUnit string) (locations []*location.QueryLocation, err error) {
	geoLocations, err := s.client.GeoSearchLocation(ctx, areaKey, &redis.GeoSearchLocationQuery{
		GeoSearchQuery: redis.GeoSearchQuery{
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
		},
		WithCoord: true,
		WithDist:  true,
		WithHash:  false,
	}).Result()

	if err != nil {
		return nil, err
	}

	length := len(geoLocations)
	reasonableLimit := 30
	expireTrackingKey := expireTrackingKey(areaKey)

	for _, geoLocation := range geoLocations {
		lastSeen := time.Time{}
		if length < reasonableLimit {
			score := s.client.ZScore(ctx, expireTrackingKey, geoLocation.Name).Val()
			lastSeen = time.UnixMilli(int64(score))
		}
		locations = append(locations, &location.QueryLocation{
			PrincipalID: geoLocation.Name,
			Longitude:   geoLocation.Longitude,
			Latitude:    geoLocation.Latitude,
			Distance:    geoLocation.Dist,
			LastSeen:    lastSeen,
		})
	}

	return locations, nil
}

func expireTrackingKey(areaKey string) string {
	return fmt.Sprintf("%s-expire-tracking", areaKey)
}
