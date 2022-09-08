package timed

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

type Member struct {
	ID            string
	PickUp        LatLon
	DropOff       LatLon
	EnqueuedTime  time.Time
	CrowFlyMeters int
}

type LatLon struct {
	Lat float64
	Lon float64
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

func (q *Queue) Enqueue(ctx context.Context, areaKey string, m Member) error {
	err := q.client.GeoAdd(ctx, areaKey, &redis.GeoLocation{
		Name:      m.ID,
		Longitude: m.PickUp.Lon,
		Latitude:  m.PickUp.Lat,
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
			Member: m.ID,
		}},
	})

	return nil
}

func (q *Queue) Dequeue(ctx context.Context, areaKey string, id string) (*Member, error) {
	m, err := q.PeekOne(ctx, areaKey, id)

	if err != nil {
		return nil, err
	}

	err = q.client.ZRem(ctx, areaKey, id).Err()
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (q *Queue) Purge(ctx context.Context, areaKey string) error {
	err := q.client.ZRemRangeByScore(ctx, areaKey, "-inf", "+inf").Err()
	if err != nil {
		return err
	}

	timeKey := enqueuedTimeKey(areaKey)
	err = q.client.ZRemRangeByScore(ctx, timeKey, "-inf", "+inf").Err()
	if err != nil {
		return err
	}
	return nil
}

func (q *Queue) PeekOne(ctx context.Context, areaKey string, id string) (*Member, error) {
	v, err := q.client.GeoPos(ctx, areaKey, id).Result()

	if err != nil {
		return nil, err
	}

	enqueuedTimeKey := enqueuedTimeKey(areaKey)
	score := q.client.ZScore(ctx, enqueuedTimeKey, id).Val()

	enqueued := time.UnixMilli(int64(score))

	// expect max one element
	if len(v) > 0 && v[0] != nil {
		return &Member{
			ID: id,
			PickUp: LatLon{
				Lon: v[0].Longitude,
				Lat: v[0].Latitude,
			},
			EnqueuedTime: enqueued,
		}, nil
	}
	return nil, nil
}
func (q *Queue) PeekMany(ctx context.Context, areaKey string, from LatLon, radiusMeters int) ([]*Member, error) {
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

	var res = make([]*Member, length)
	for i, location := range locations {
		t := time.Time{}
		if length < reasonableLimit {
			score := q.client.ZScore(ctx, timeKey, location.Name).Val()
			t = time.UnixMilli(int64(score))
		}
		res[i] = &Member{
			ID: location.Name,
			PickUp: LatLon{
				Lat: location.Latitude,
				Lon: location.Longitude,
			},
			CrowFlyMeters: int(location.Dist),
			EnqueuedTime:  t,
		}
	}

	return res, nil
}

func (q *Queue) PeekManyAddedBefore(ctx context.Context, areaKey string, addedBefore time.Time, limit int) ([]*Member, error) {
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

	var res = make([]*Member, 0)
	result, err := q.client.GeoPos(ctx, areaKey, keys...).Result()
	if err != nil {
		return nil, err
	}
	for i, pos := range result {
		if pos == nil {
			continue
		}
		t := time.Time{}
		if length < reasonableLimit {
			score := q.client.ZScore(ctx, timeKey, keys[i]).Val()
			t = time.UnixMilli(int64(score))
		}
		res = append(res, &Member{
			ID: keys[i],
			PickUp: LatLon{
				Lon: pos.Longitude,
				Lat: pos.Latitude,
			},
			CrowFlyMeters: 0,
			EnqueuedTime:  t,
		})
	}
	return res, nil
}

func enqueuedTimeKey(areaKey string) string {
	return fmt.Sprintf("%s_request_enqueued_time", areaKey)
}
