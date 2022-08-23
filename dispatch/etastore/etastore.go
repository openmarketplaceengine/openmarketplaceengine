package etastore

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/openmarketplaceengine/openmarketplaceengine/dispatch/job"
)

type Handler func(ctx context.Context, geoHash string, eJobs []*job.EstimatedJob) error

type EtaStore struct {
	client *redis.Client
}

func NewEtaStore(client *redis.Client) *EtaStore {
	return &EtaStore{
		client: client,
	}
}

func (s *EtaStore) Get(ctx context.Context, geoHash string) ([]*job.EstimatedJob, error) {
	jobs, err := s.client.HGetAll(ctx, geoHash).Result()

	if err != nil {
		return nil, fmt.Errorf("hgetall error: %w", err)
	}

	l := len(jobs)
	if l == 0 {
		return nil, nil
	}

	result := make([]*job.EstimatedJob, l)

	for k, v := range jobs {
		var ej job.EstimatedJob
		err := ej.UnmarshalBinary([]byte(v))
		if err != nil {
			return nil, fmt.Errorf("unmarshal geohash %s job %s error: %w", geoHash, k, err)
		}
		result = append(result, &ej)
	}

	return result, nil
}

func (s *EtaStore) Store(ctx context.Context, geoHash string, jobs []*job.EstimatedJob) error {
	values := toValues(jobs)
	err := s.client.HSet(ctx, geoHash, values...).Err()

	if err != nil {
		return fmt.Errorf("hset error: %w", err)
	}

	return nil
}

func toValues(jobs []*job.EstimatedJob) []interface{} {
	res := make([]interface{}, 0)
	for _, j := range jobs {
		res = append(res, j.ID)
		res = append(res, j)
	}
	return res
}
