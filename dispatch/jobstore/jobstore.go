package jobstore

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
)

type Job struct {
	ID      string
	PickUp  LatLon
	DropOff LatLon
}

func (e Job) MarshalBinary() ([]byte, error) {
	return json.Marshal(e)
}

func (e *Job) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &e)
}

type LatLon struct {
	Lat float64
	Lon float64
}

type JobStore struct {
	client *redis.Client
}

func NewJobStore(client *redis.Client) *JobStore {
	return &JobStore{
		client: client,
	}
}

func (s *JobStore) GetAll(ctx context.Context, areaKey string, members ...string) ([]*Job, error) {
	k := key(areaKey)
	v, err := s.client.HMGet(ctx, k, members...).Result()

	if err != nil {
		return nil, fmt.Errorf("hgetall error: %w", err)
	}

	l := len(v)
	if l == 0 {
		return nil, nil
	}

	result := make([]*Job, l)

	for i, v := range v {
		var m Job
		err := json.Unmarshal([]byte(v.(string)), &m)
		if err != nil {
			return nil, fmt.Errorf("unmarshal error: %w", err)
		}
		result[i] = &m
	}

	return result, nil
}

func (s *JobStore) Store(ctx context.Context, areaKey string, members []*Job) error {
	values := toValues(members)
	k := key(areaKey)
	err := s.client.HSet(ctx, k, values...).Err()

	if err != nil {
		return fmt.Errorf("hset error: %w", err)
	}

	return nil
}

func key(areaKey string) string {
	return fmt.Sprintf("%s-%s", areaKey, "jobs")
}

func toValues(estimates []*Job) []interface{} {
	res := make([]interface{}, 0)
	for _, e := range estimates {
		res = append(res, e.ID)
		res = append(res, e)
	}
	return res
}
