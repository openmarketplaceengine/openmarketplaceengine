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

func (s *JobStore) GetByIds(ctx context.Context, areaKey string, ids ...string) ([]*Job, error) {
	k := key(areaKey)
	v, err := s.client.HMGet(ctx, k, ids...).Result()

	if err != nil {
		return nil, fmt.Errorf("get all error: %w", err)
	}

	l := len(v)
	if l == 0 {
		return nil, nil
	}

	result := make([]*Job, 0)

	for _, v := range v {
		if v == nil {
			continue
		}
		var m Job
		err := json.Unmarshal([]byte(v.(string)), &m)
		if err != nil {
			return nil, fmt.Errorf("unmarshal error: %w", err)
		}
		result = append(result, &m)
	}

	return result, nil
}

func (s *JobStore) GetAll(ctx context.Context, areaKey string) ([]*Job, error) {
	k := key(areaKey)
	v, err := s.client.HGetAll(ctx, k).Result()

	if err != nil {
		return nil, fmt.Errorf("get all error: %w", err)
	}

	l := len(v)
	if l == 0 {
		return nil, nil
	}

	result := make([]*Job, l)
	i := 0
	for _, v := range v {
		var m Job
		err := json.Unmarshal([]byte(v), &m)
		if err != nil {
			return nil, fmt.Errorf("unmarshal error: %w", err)
		}
		result[i] = &m
		i++
	}

	return result, nil
}

func (s *JobStore) DeleteAll(ctx context.Context, areaKey string) (int64, error) {
	k := key(areaKey)
	v, err := s.client.Del(ctx, k).Result()

	if err != nil {
		return 0, fmt.Errorf("delete error: %w", err)
	}

	return v, nil
}

func (s *JobStore) StoreMany(ctx context.Context, areaKey string, jobs []*Job) error {
	values := toValues(jobs...)
	k := key(areaKey)
	err := s.client.HSet(ctx, k, values...).Err()

	if err != nil {
		return fmt.Errorf("store many error: %w", err)
	}

	return nil
}

func (s *JobStore) StoreOne(ctx context.Context, areaKey string, job *Job) error {
	values := toValues(job)
	k := key(areaKey)
	err := s.client.HSet(ctx, k, values...).Err()

	if err != nil {
		return fmt.Errorf("store one error: %w", err)
	}

	return nil
}

func (s *JobStore) Remove(ctx context.Context, areaKey string, ids ...string) error {
	k := key(areaKey)
	err := s.client.HDel(ctx, k, ids...).Err()

	if err != nil {
		return fmt.Errorf("remove one error: %w", err)
	}

	return nil
}

func key(areaKey string) string {
	return fmt.Sprintf("%s-%s", areaKey, "jobs")
}

func toValues(jobs ...*Job) []interface{} {
	res := make([]interface{}, 0)
	for _, e := range jobs {
		res = append(res, e.ID)
		res = append(res, e)
	}
	return res
}
