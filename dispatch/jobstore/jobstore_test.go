package jobstore

import (
	"context"
	"github.com/google/uuid"
	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
	"github.com/openmarketplaceengine/openmarketplaceengine/dao"
	"github.com/stretchr/testify/require"
	"testing"
)

const areaKey = "jobstore_test"

var jobs = []*Job{
	{
		ID: uuid.NewString(),
		PickUp: LatLon{
			Lat: 0,
			Lon: 0,
		},
		DropOff: LatLon{
			Lat: 0,
			Lon: 0,
		},
	},
	{
		ID: uuid.NewString(),
		PickUp: LatLon{
			Lat: 0,
			Lon: 0,
		},
		DropOff: LatLon{
			Lat: 0,
			Lon: 0,
		},
	},
}

func TestStore(t *testing.T) {
	err := cfg.Load()
	require.NoError(t, err)

	if !dao.Reds.State.Running() {
		require.NoError(t, dao.Reds.Boot())
	}

	client := dao.Reds.StoreClient

	store := NewJobStore(client)

	t.Run("testStore", func(t *testing.T) {
		testStore(t, store)
	})
}

func testStore(t *testing.T, store *JobStore) {
	ctx := context.Background()

	err := store.StoreMany(ctx, areaKey, jobs)
	require.NoError(t, err)

	retrieved, err := store.GetAll(ctx, areaKey, jobs[0].ID, jobs[1].ID)
	require.NoError(t, err)
	require.NotNil(t, retrieved)
	require.Equal(t, jobs, retrieved)
}
