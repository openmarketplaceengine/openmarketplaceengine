package queue

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/openmarketplaceengine/openmarketplaceengine/pkg/util"
	"golang.org/x/net/context"
)

type EnqueuedRequest struct {
	ID            string
	PickUp        LatLon
	DropOff       LatLon
	EnqueuedTime  time.Time
	CrowFlyMeters int
}

type Queue struct {
	client *redis.Client
}

func NewQueue(client *redis.Client) *Queue {
	q := &Queue{
		client: client,
	}

	return q
}

func (q *Queue) Enqueue(ctx context.Context, areaKey string, request Request) error {
	err := q.client.GeoAdd(ctx, areaKey, &redis.GeoLocation{
		Name:      request.ID,
		Longitude: request.PickUp.Lon,
		Latitude:  request.PickUp.Lat,
		Dist:      0,
		GeoHash:   0,
	}).Err()
	if err != nil {
		return err
	}

	t := time.Now()
	timeKey := enqueuedTimeKey(areaKey)
	q.client.ZAddArgs(ctx, timeKey, redis.ZAddArgs{
		NX: false,
		XX: false,
		LT: false,
		GT: false,
		Ch: false,
		Members: []redis.Z{{
			Score:  float64(t.UnixMilli()),
			Member: request.ID,
		}},
	})

	return nil
}

func (q *Queue) Peek(ctx context.Context, areaKey string, id string) *EnqueuedRequest {
	v := q.client.GeoPos(ctx, areaKey, id).Val()

	enqueuedTimeKey := enqueuedTimeKey(areaKey)
	score := q.client.ZScore(ctx, enqueuedTimeKey, id).Val()

	enqueued := time.UnixMilli(int64(score))

	// expect max one element
	if len(v) > 0 && v[0] != nil {
		return &EnqueuedRequest{
			ID: id,
			PickUp: LatLon{
				Lon: util.Round6(v[0].Longitude),
				Lat: util.Round6(v[0].Latitude),
			},
			EnqueuedTime: enqueued,
		}
	}
	return nil
}
func (q *Queue) GetNearbyRequests(ctx context.Context, areaKey string, from LatLon, radiusMeters int) ([]*EnqueuedRequest, error) {
	locations, err := q.client.GeoSearchLocation(ctx, areaKey, &redis.GeoSearchLocationQuery{
		GeoSearchQuery: redis.GeoSearchQuery{
			Member:     "",
			Longitude:  from.Lon,
			Latitude:   from.Lat,
			Radius:     float64(radiusMeters),
			RadiusUnit: "m",
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

	length := len(locations)
	reasonableLimit := 30
	timeKey := enqueuedTimeKey(areaKey)

	var res = make([]*EnqueuedRequest, length)
	for i, location := range locations {
		t := time.Time{}
		if length < reasonableLimit {
			score := q.client.ZScore(ctx, timeKey, location.Name).Val()
			t = time.UnixMilli(int64(score))
		}
		res[i] = &EnqueuedRequest{
			ID: location.Name,
			PickUp: LatLon{
				Lat: util.Round6(location.Latitude),
				Lon: util.Round6(location.Longitude),
			},
			CrowFlyMeters: int(location.Dist),
			EnqueuedTime:  t,
		}
	}

	return res, nil
}

func (q *Queue) GetAddedBeforeTime(ctx context.Context, areaKey string, addedBefore time.Time, limit int) ([]*EnqueuedRequest, error) {
	key := enqueuedTimeKey(areaKey)

	min := fmt.Sprintf("%v", float64(0))
	max := fmt.Sprintf("%v", float64(addedBefore.UnixMilli()))

	keys, err := q.client.ZRangeByScore(ctx, key, &redis.ZRangeBy{
		Min:   min,
		Max:   max,
		Count: int64(limit),
	}).Result()

	if err != nil {
		return nil, err
	}

	length := len(keys)
	reasonableLimit := 30
	timeKey := enqueuedTimeKey(areaKey)

	var res = make([]*EnqueuedRequest, length)
	result, err := q.client.GeoPos(ctx, areaKey, keys...).Result()
	if err != nil {
		return nil, err
	}
	for i, pos := range result {
		t := time.Time{}
		if length < reasonableLimit {
			score := q.client.ZScore(ctx, timeKey, keys[i]).Val()
			t = time.UnixMilli(int64(score))
		}
		res[i] = &EnqueuedRequest{
			ID: keys[i],
			PickUp: LatLon{
				Lon: util.Round6(pos.Longitude),
				Lat: util.Round6(pos.Latitude),
			},
			CrowFlyMeters: 0,
			EnqueuedTime:  t,
		}
	}
	return res, nil
}

func enqueuedTimeKey(areaKey string) string {
	return fmt.Sprintf("%s_request_enqueued_time", areaKey)
}
