package queue

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/openmarketplaceengine/openmarketplaceengine/pkg/util"
	"golang.org/x/net/context"
)

type Location struct {
	Lon float64
	Lat float64
}

type Job struct {
	ID              string
	PickupLocation  Location
	CrowFlyDistance float64
	EnqueuedTime    time.Time
}

type Driver struct {
	ID        string
	Location  Location
	AddedTime time.Time
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

func (q *Queue) Add(ctx context.Context, areaKey string, job Job) error {
	err := q.client.GeoAdd(ctx, areaKey, &redis.GeoLocation{
		Name:      job.ID,
		Longitude: job.PickupLocation.Lon,
		Latitude:  job.PickupLocation.Lat,
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
			Member: job.ID,
		}},
	})

	return nil
}

func (q *Queue) GetJob(ctx context.Context, areaKey string, jobID string) *Job {
	v := q.client.GeoPos(ctx, areaKey, jobID).Val()

	enqueuedTimeKey := enqueuedTimeKey(areaKey)
	score := q.client.ZScore(ctx, enqueuedTimeKey, jobID).Val()

	enqueued := time.UnixMilli(int64(score))

	// expect max one element
	if len(v) > 0 && v[0] != nil {
		return &Job{
			ID: jobID,
			PickupLocation: Location{
				Lon: util.Round6(v[0].Longitude),
				Lat: util.Round6(v[0].Latitude),
			},
			EnqueuedTime: enqueued,
		}
	}
	return nil
}
func (q *Queue) GetNearbyJobs(ctx context.Context, areaKey string, location Location, radius float64, radiusUnit string) ([]*Job, error) {
	locations, err := q.client.GeoSearchLocation(ctx, areaKey, &redis.GeoSearchLocationQuery{
		GeoSearchQuery: redis.GeoSearchQuery{
			Member:     "",
			Longitude:  location.Lon,
			Latitude:   location.Lat,
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

	length := len(locations)
	reasonableLimit := 30
	timeKey := enqueuedTimeKey(areaKey)

	var jobs = make([]*Job, length)
	for i, location := range locations {
		t := time.Time{}
		if length < reasonableLimit {
			score := q.client.ZScore(ctx, timeKey, location.Name).Val()
			t = time.UnixMilli(int64(score))
		}
		jobs[i] = &Job{
			ID: location.Name,
			PickupLocation: Location{
				Lon: util.Round6(location.Longitude),
				Lat: util.Round6(location.Latitude),
			},
			CrowFlyDistance: location.Dist,
			EnqueuedTime:    t,
		}
	}

	return jobs, nil
}

func (q *Queue) GetAddedBeforeTimeJobs(ctx context.Context, areaKey string, addedBefore time.Time, limit int) ([]*Job, error) {
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

	var jobs = make([]*Job, length)
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
		jobs[i] = &Job{
			ID: keys[i],
			PickupLocation: Location{
				Lon: util.Round6(pos.Longitude),
				Lat: util.Round6(pos.Latitude),
			},
			CrowFlyDistance: 0,
			EnqueuedTime:    t,
		}
	}
	return jobs, nil
}

func enqueuedTimeKey(areaKey string) string {
	return fmt.Sprintf("%s_job_enqueued_time", areaKey)
}
