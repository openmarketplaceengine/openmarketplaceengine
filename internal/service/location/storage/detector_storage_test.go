package storage

import (
	"context"
	"testing"

	"github.com/openmarketplaceengine/openmarketplaceengine/dao"

	"github.com/google/uuid"
	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	err := cfg.Load()
	require.NoError(t, err)

	if !dao.Reds.State.Running() {
		require.NoError(t, dao.Reds.Boot())
	}

	storage := NewRedisStorage(dao.Reds.StoreClient)

	ctx := context.Background()

	t.Run("testVisits", func(t *testing.T) {
		testVisits(ctx, t, storage)
	})
}

func testVisits(ctx context.Context, t *testing.T, storage *RedisStorage) {
	tollgateID := "toll-123"
	subjectID := uuid.NewString()
	key := Key(tollgateID, subjectID)
	err := storage.Visit(ctx, key, 5, 0)
	require.NoError(t, err)
	err = storage.Visit(ctx, key, 5, 1)
	require.NoError(t, err)
	err = storage.Visit(ctx, key, 5, 5)
	require.NoError(t, err)

	res, err := storage.Visits(ctx, key, 5)
	require.NoError(t, err)
	require.Equal(t, []int{1, 1, 0, 0, 0}, res)

	res, err = storage.Visits(ctx, key, 6)
	require.NoError(t, err)
	require.Equal(t, []int{1, 1, 0, 0, 0, 1}, res)
}
