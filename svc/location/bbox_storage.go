package location

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type BBoxStorage struct {
	client *redis.Client
}

func NewBBoxStorage(client *redis.Client) *BBoxStorage {
	s := BBoxStorage{
		client: client,
	}
	return &s
}

func Key(tollgateID, subjectID string) string {
	return fmt.Sprintf("toll-bbox-%s-%s", tollgateID, subjectID)
}

func (s *BBoxStorage) Visit(ctx context.Context, key string, size int, index int) error {
	return s.client.BitField(ctx, key, "SET", "i2", fmt.Sprintf("#%v", index), 1).Err()
}

func (s *BBoxStorage) Visits(ctx context.Context, key string, size int) ([]int, error) {
	var args []interface{}
	for i := 0; i < size; i++ {
		args = append(args, "GET")
		args = append(args, "i2")
		args = append(args, fmt.Sprintf("#%v", i))
	}
	int64s, err := s.client.BitField(ctx, key, args...).Result()
	if err != nil {
		return nil, err
	}
	res := make([]int, len(int64s))
	for i, v := range int64s {
		res[i] = int(v)
	}
	return res, err
}

func (s *BBoxStorage) Del(ctx context.Context, key string) error {
	return s.client.Del(ctx, key).Err()
}
