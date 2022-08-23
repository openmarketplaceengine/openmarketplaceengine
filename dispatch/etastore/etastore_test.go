package etastore

import (
	"context"
	"testing"

	"github.com/openmarketplaceengine/openmarketplaceengine/dao"
	"github.com/openmarketplaceengine/openmarketplaceengine/dispatch/job"

	"github.com/google/uuid"
	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
	"github.com/stretchr/testify/require"
)

var jobs = []*job.EstimatedJob{
	{
		ID: uuid.NewString(),
		ToPickup: job.Estimate{
			DistanceMeters: 1,
			Duration:       1,
		},
		PickupToDropOff: job.Estimate{
			DistanceMeters: 2,
			Duration:       2,
		},
		WorkerLocation: job.Location{
			Address: "a",
			Lat:     0,
			Lon:     0,
		},
		Pickup: job.Location{
			Address: "b",
			Lat:     0,
			Lon:     0,
		},
		DropOff: job.Location{
			Address: "c",
			Lat:     0,
			Lon:     0,
		},
	},
	{
		ID: uuid.NewString(),
		ToPickup: job.Estimate{
			DistanceMeters: 1,
			Duration:       1,
		},
		PickupToDropOff: job.Estimate{
			DistanceMeters: 2,
			Duration:       2,
		},
		WorkerLocation: job.Location{
			Address: "a",
			Lat:     0,
			Lon:     0,
		},
		Pickup: job.Location{
			Address: "b",
			Lat:     0,
			Lon:     0,
		},
		DropOff: job.Location{
			Address: "c",
			Lat:     0,
			Lon:     0,
		},
	},
}

func TestEtaStore(t *testing.T) {
	err := cfg.Load()
	require.NoError(t, err)

	if !dao.Reds.State.Running() {
		require.NoError(t, dao.Reds.Boot())
	}

	client := dao.Reds.StoreClient

	store := NewEtaStore(client)

	t.Run("testStore", func(t *testing.T) {
		testStore(t, store)
	})
}

func testStore(t *testing.T, store *EtaStore) {
	ctx := context.Background()

	ll0 := job.LatLon{
		Lat: 40.636916,
		Lon: -74.195995,
	}
	_ = job.LatLon{
		Lat: 40.634408,
		Lon: -74.198356,
	}

	hash := job.ToGeoHash(ll0)
	err := store.Store(ctx, hash, jobs)
	require.NoError(t, err)

	estimatedJobs, err := store.Get(ctx, hash)
	require.NoError(t, err)
	require.NotNil(t, estimatedJobs)
}
