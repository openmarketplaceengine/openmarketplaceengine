package bbox

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
	redisClient "github.com/openmarketplaceengine/openmarketplaceengine/redis/client"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	err := cfg.Load()
	require.NoError(t, err)

	client := redisClient.NewStoreClient()
	require.NotNil(t, client)

	storage := NewStorage(client)

	ctx := context.Background()

	t.Run("testVisits", func(t *testing.T) {
		testVisits(ctx, t, storage)
	})
}

func testVisits(ctx context.Context, t *testing.T, storage Storage) {
	tollgateID := "toll-123"
	subjectID := uuid.NewString()
	err := storage.Visit(ctx, tollgateID, subjectID, 0)
	require.NoError(t, err)
	err = storage.Visit(ctx, tollgateID, subjectID, 1)
	require.NoError(t, err)
	err = storage.Visit(ctx, tollgateID, subjectID, 5)
	require.NoError(t, err)

	res, err := storage.Visits(ctx, tollgateID, subjectID, 5)
	require.NoError(t, err)
	require.Equal(t, []int64{1, 1, 0, 0, 0}, res)

	res, err = storage.Visits(ctx, tollgateID, subjectID, 6)
	require.NoError(t, err)
	require.Equal(t, []int64{1, 1, 0, 0, 0, 1}, res)
}
