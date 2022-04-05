package detector

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

	storage := newStorage(dao.Reds.StoreClient)

	ctx := context.Background()

	t.Run("testVisits", func(t *testing.T) {
		testVisits(ctx, t, storage)
	})
}

func testVisits(ctx context.Context, t *testing.T, storage *storage) {
	tollgateID := "toll-123"
	subjectID := uuid.NewString()
	err := storage.visit(ctx, tollgateID, subjectID, 0)
	require.NoError(t, err)
	err = storage.visit(ctx, tollgateID, subjectID, 1)
	require.NoError(t, err)
	err = storage.visit(ctx, tollgateID, subjectID, 5)
	require.NoError(t, err)

	res, err := storage.visits(ctx, tollgateID, subjectID, 5)
	require.NoError(t, err)
	require.Equal(t, []int64{1, 1, 0, 0, 0}, res)

	res, err = storage.visits(ctx, tollgateID, subjectID, 6)
	require.NoError(t, err)
	require.Equal(t, []int64{1, 1, 0, 0, 0, 1}, res)
}
