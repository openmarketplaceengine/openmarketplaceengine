package detector

import "context"

// Storage keeps passing through information for WorkerID
// visits is a flag bitmap array for corresponding BBox array, 1 - means visited, 0 - not visited
// Examples:
// 1 - [1,0,0,0,0] - for bbox size 5 elements array the first one was visited.
// 2 - [1,0,1,0,0]
// 3 - [1,0,1,1,0].
type Storage interface {
	Visit(ctx context.Context, key string, size int, index int) error
	Visits(ctx context.Context, key string, size int) ([]int, error)
	Del(ctx context.Context, key string) error
}

type MapStorage struct {
	data map[string][]int
}

func NewMapStorage() *MapStorage {
	s := MapStorage{
		data: make(map[string][]int),
	}
	return &s
}

func (s *MapStorage) Visit(ctx context.Context, key string, size int, index int) error {
	for len(s.data[key]) < size {
		s.data[key] = append(s.data[key], 0)
	}
	s.data[key][index] = 1
	return nil
}

func (s *MapStorage) Visits(ctx context.Context, key string, size int) ([]int, error) {
	v := s.data[key]
	return v, nil
}

func (s *MapStorage) Del(ctx context.Context, key string) error {
	delete(s.data, key)
	return nil
}
