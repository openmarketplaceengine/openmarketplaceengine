package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type Storage struct {
	client *redis.Client
}

func New(client *redis.Client) *Storage {
	s := Storage{
		client: client,
	}
	return &s
}

func (s *Storage) Update(ctx context.Context, areaKey string, location *Location) (err error) {
	err = s.client.GeoAdd(ctx, areaKey, &redis.GeoLocation{
		Name:      location.WorkerID,
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
			Member: location.WorkerID,
		}},
	})

	return nil
}

func (s *Storage) LastLocation(ctx context.Context, areaKey string, workerID string) *LastLocation {
	v := s.client.GeoPos(ctx, areaKey, workerID).Val()

	expireTrackingKey := expireTrackingKey(areaKey)
	score := s.client.ZScore(ctx, expireTrackingKey, workerID).Val()

	lastSeen := time.UnixMilli(int64(score))
	// At the moment we expect max one element
	if len(v) > 0 && v[0] != nil {
		return &LastLocation{
			WorkerID:     workerID,
			Longitude:    v[0].Longitude,
			Latitude:     v[0].Latitude,
			LastSeenTime: lastSeen,
		}
	}
	return nil
}

func (s *Storage) ForgetLocation(ctx context.Context, areaKey string, workerID string) (err error) {
	err = s.client.ZRem(ctx, areaKey, workerID).Err()
	if err != nil {
		return err
	}
	expireTrackingKey := expireTrackingKey(areaKey)
	err = s.client.ZRem(ctx, expireTrackingKey, workerID).Err()
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) RemoveExpiredLocations(ctx context.Context, areaKey string, before time.Time) (err error) {
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
			return fmt.Errorf("Storage ZRemRangeByScore error: %s", innerErr)
		}
		innerErr = s.client.ZRem(ctx, areaKey, keys).Err()
		if innerErr != nil {
			return fmt.Errorf("Storage HDel error: %s", innerErr)
		}
	}
	return nil
}

func (s *Storage) RangeLocations(ctx context.Context, areaKey string, fromLongitude float64, fromLatitude float64, radius float64, radiusUnit string) (locations []*RangeLocation, err error) {
	geoLocations, err := s.client.GeoSearchLocation(ctx, areaKey, &redis.GeoSearchLocationQuery{
		GeoSearchQuery: redis.GeoSearchQuery{
			Member:     "",
			Longitude:  fromLongitude,
			Latitude:   fromLatitude,
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
		locations = append(locations, &RangeLocation{
			WorkerID:      geoLocation.Name,
			Longitude:     geoLocation.Longitude,
			Latitude:      geoLocation.Latitude,
			Distance:      geoLocation.Dist,
			FromLatitude:  fromLatitude,
			FromLongitude: fromLongitude,
			LastSeenTime:  lastSeen,
		})
	}

	return locations, nil
}

func expireTrackingKey(areaKey string) string {
	return fmt.Sprintf("%s-expire-tracking", areaKey)
}
