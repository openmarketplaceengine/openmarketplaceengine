package geoqueue

import (
	"context"
	"github.com/go-redis/redis/v8"
)

type Member struct {
	ID     string
	PickUp LatLon
}

type LatLon struct {
	Lat float64
	Lon float64
}

type GeoQueue struct {
	client *redis.Client
}

func NewGeoQueue(client *redis.Client) *GeoQueue {
	q := &GeoQueue{
		client: client,
	}

	return q
}

func (q *GeoQueue) Enqueue(ctx context.Context, areaKey string, m Member) error {
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

	return nil
}

func (q *GeoQueue) Dequeue(ctx context.Context, areaKey string, id string) (*Member, error) {
	m, _ := q.PeekOne(ctx, areaKey, id)

	err := q.client.ZRem(ctx, areaKey, id).Err()
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (q *GeoQueue) Purge(ctx context.Context, areaKey string) error {
	err := q.client.ZRemRangeByScore(ctx, areaKey, "-inf", "+inf").Err()
	if err != nil {
		return err
	}

	return nil
}

func (q *GeoQueue) PeekOne(ctx context.Context, areaKey string, id string) (*Member, error) {
	v, err := q.client.GeoPos(ctx, areaKey, id).Result()

	if err != nil {
		return nil, err
	}

	// expect max one element
	if len(v) > 0 && v[0] != nil {
		return &Member{
			ID: id,
			PickUp: LatLon{
				Lon: v[0].Longitude,
				Lat: v[0].Latitude,
			},
		}, nil
	}
	return nil, nil

}
func (q *GeoQueue) PeekMany(ctx context.Context, areaKey string, from LatLon, radiusMeters int) ([]*Member, error) {
	locations, err := q.client.GeoSearchLocation(ctx, areaKey, &redis.GeoSearchLocationQuery{
		GeoSearchQuery: redis.GeoSearchQuery{
			Longitude:  from.Lon,
			Latitude:   from.Lat,
			Radius:     float64(radiusMeters),
			RadiusUnit: "m",
			Sort:       "ASC",
		},
		WithCoord: true,
		WithDist:  true,
		WithHash:  false,
	}).Result()

	if err != nil {
		return nil, err
	}

	length := len(locations)

	var res = make([]*Member, length)
	for i, location := range locations {
		res[i] = &Member{
			ID: location.Name,
			PickUp: LatLon{
				Lat: location.Latitude,
				Lon: location.Longitude,
			},
		}
	}

	return res, nil
}
