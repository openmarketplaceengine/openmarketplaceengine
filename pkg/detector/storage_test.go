package detector

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	storage := NewMapStorage()

	t.Run("testVisits", func(t *testing.T) {
		testVisits(t, storage)
	})
}

func testVisits(t *testing.T, storage Storage) {
	ctx := context.Background()
	key := uuid.NewString()
	_ = storage.Visit(ctx, key, 5, 0)
	_ = storage.Visit(ctx, key, 5, 1)
	_ = storage.Visit(ctx, key, 5, 4)

	visits, _ := storage.Visits(ctx, key, 5)
	require.Equal(t, []int{1, 1, 0, 0, 1}, visits)
	_ = storage.Visit(ctx, key, 6, 5)
	visits, _ = storage.Visits(ctx, key, 5)
	require.Equal(t, []int{1, 1, 0, 0, 1, 1}, visits)
}
