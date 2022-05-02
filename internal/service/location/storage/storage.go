package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/openmarketplaceengine/openmarketplaceengine/pkg/util"

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

func (s *Storage) Update(ctx context.Context, areaKey string, l *Location, t time.Time) (err error) {
	err = s.client.GeoAdd(ctx, areaKey, &redis.GeoLocation{
		Name:      l.WorkerID,
		Longitude: l.Longitude,
		Latitude:  l.Latitude,
		Dist:      0,
		GeoHash:   0,
	}).Err()
	if err != nil {
		return err
	}

	updateTimeKey := updateTimeKey(areaKey)
	s.client.ZAddArgs(ctx, updateTimeKey, redis.ZAddArgs{
		NX: false,
		XX: false,
		LT: false,
		GT: false,
		Ch: false,
		Members: []redis.Z{{
			Score:  float64(t.UnixMilli()),
			Member: l.WorkerID,
		}},
	})

	return nil
}

func (s *Storage) LastLocation(ctx context.Context, areaKey string, workerID string) *LastLocation {
	v := s.client.GeoPos(ctx, areaKey, workerID).Val()

	updateTimeKey := updateTimeKey(areaKey)
	score := s.client.ZScore(ctx, updateTimeKey, workerID).Val()

	lastSeen := time.UnixMilli(int64(score))
	// At the moment we expect max one element
	if len(v) > 0 && v[0] != nil {
		return &LastLocation{
			WorkerID:     workerID,
			Longitude:    util.Round6(v[0].Longitude),
			Latitude:     util.Round6(v[0].Latitude),
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
	updateTimeKey := updateTimeKey(areaKey)
	err = s.client.ZRem(ctx, updateTimeKey, workerID).Err()
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) RemoveExpiredLocations(ctx context.Context, areaKey string, before time.Time) (err error) {
	key := updateTimeKey(areaKey)
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
	updateTimeKey := updateTimeKey(areaKey)

	for _, geoLocation := range geoLocations {
		lastSeen := time.Time{}
		if length < reasonableLimit {
			score := s.client.ZScore(ctx, updateTimeKey, geoLocation.Name).Val()
			lastSeen = time.UnixMilli(int64(score))
		}
		locations = append(locations, &RangeLocation{
			WorkerID:      geoLocation.Name,
			Longitude:     util.Round6(geoLocation.Longitude),
			Latitude:      util.Round6(geoLocation.Latitude),
			Distance:      geoLocation.Dist,
			FromLatitude:  fromLatitude,
			FromLongitude: fromLongitude,
			LastSeenTime:  lastSeen,
		})
	}

	return locations, nil
}

func updateTimeKey(areaKey string) string {
	return fmt.Sprintf("%s_update_time", areaKey)
}
