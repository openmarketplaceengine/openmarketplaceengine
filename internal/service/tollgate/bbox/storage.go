package bbox

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

// Storage keeps passing through information for WorkerID
// visits is a flag bitmap array for corresponding BBox array, 1 - means visited, 0 - not visited
// Examples:
// 1 - [1,0,0,0,0] - for bbox 5 elements array the first one was visited.
// 2 - [1,0,1,0,0]
// 3 - [1,0,1,1,0].
type Storage interface {
	Visit(ctx context.Context, tollgateID, subjectID string, bBoxIdx int) error
	Visits(ctx context.Context, tollgateID, subjectID string, size int) ([]int64, error)
	Del(ctx context.Context, tollgateID, subjectID string) error
}

type storage struct {
	client *redis.Client
}

func NewStorage(client *redis.Client) Storage {
	s := storage{
		client: client,
	}
	return &s
}

func storageKey(tollgateID, subjectID string) string {
	return fmt.Sprintf("toll-bbox-%s-%s", tollgateID, subjectID)
}

func (s *storage) Visit(ctx context.Context, tollgateID, subjectID string, index int) error {
	key := storageKey(tollgateID, subjectID)
	return s.client.BitField(ctx, key, "SET", "i2", fmt.Sprintf("#%v", index), 1).Err()
}

func (s *storage) Visits(ctx context.Context, tollgateID, subjectID string, size int) ([]int64, error) {
	var args []interface{}
	for i := 0; i < size; i++ {
		args = append(args, "GET")
		args = append(args, "i2")
		args = append(args, fmt.Sprintf("#%v", i))
	}
	key := storageKey(tollgateID, subjectID)
	return s.client.BitField(ctx, key, args...).Result()
}

func (s *storage) Del(ctx context.Context, tollgateID, subjectID string) error {
	key := storageKey(tollgateID, subjectID)
	return s.client.Del(ctx, key).Err()
}
