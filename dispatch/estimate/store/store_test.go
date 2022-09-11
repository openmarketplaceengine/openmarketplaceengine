package store

import (
	"context"
	"github.com/google/uuid"
	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
	"github.com/openmarketplaceengine/openmarketplaceengine/dao"
	"github.com/openmarketplaceengine/openmarketplaceengine/dispatch/estimate"
	"github.com/openmarketplaceengine/openmarketplaceengine/dispatch/geohash"
	"github.com/stretchr/testify/require"
	"testing"
)

var estimates = []*estimate.Estimate{
	{
		ID: uuid.NewString(),
		ToPickup: estimate.Eta{
			DistanceMeters: 1,
			Duration:       1,
		},
		PickupToDropOff: estimate.Eta{
			DistanceMeters: 2,
			Duration:       2,
		},
		From: estimate.Location{
			Address: "a",
			Lat:     0,
			Lon:     0,
		},
		PickUp: estimate.Location{
			Address: "b",
			Lat:     0,
			Lon:     0,
		},
		DropOff: estimate.Location{
			Address: "c",
			Lat:     0,
			Lon:     0,
		},
	},
	{
		ID: uuid.NewString(),
		ToPickup: estimate.Eta{
			DistanceMeters: 1,
			Duration:       1,
		},
		PickupToDropOff: estimate.Eta{
			DistanceMeters: 2,
			Duration:       2,
		},
		From: estimate.Location{
			Address: "a",
			Lat:     0,
			Lon:     0,
		},
		PickUp: estimate.Location{
			Address: "b",
			Lat:     0,
			Lon:     0,
		},
		DropOff: estimate.Location{
			Address: "c",
			Lat:     0,
			Lon:     0,
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

	store := NewEstimateStore(client)

	t.Run("testStore", func(t *testing.T) {
		testStore(t, store)
	})
}

func testStore(t *testing.T, store *EstimateStore) {
	ctx := context.Background()

	ll0 := estimate.LatLon{
		Lat: 40.636916,
		Lon: -74.195995,
	}
	_ = estimate.LatLon{
		Lat: 40.634408,
		Lon: -74.198356,
	}

	hash := geohash.ToGeoHash(ll0.Lat, ll0.Lon, geohash.Precision800)
	err := store.Store(ctx, hash, Radius2000m, estimates)
	require.NoError(t, err)

	retrieved, err := store.GetAll(ctx, hash, Radius2000m)
	require.NoError(t, err)
	require.NotNil(t, retrieved)
}
