package estimate

import (
	"testing"

	"github.com/google/uuid"
	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService(t *testing.T) {
	if apiKey == "" {
		t.Skip("OME_GOOGLE_API_KEY env var is not set, skipping.")
	}

	s := NewBatchService(apiKey, 3, 10)

	t.Run("testGetEstimatePickupToDropOff", func(t *testing.T) {
		testGetEstimatePickupToDropOff(t, s)
	})

	t.Run("testToChunks", func(t *testing.T) {
		testToChunks(t)
	})
}

func testGetEstimatePickupToDropOff(t *testing.T, batchService *BatchService) {
	ctx := cfg.Context()

	from := LatLon{
		Lat: 40.633650,
		Lon: -74.143650,
	}

	d1 := &Request{
		ID: uuid.NewString(),
		PickUp: LatLon{
			Lat: 40.636916,
			Lon: -74.195995,
		},
		DropOff: LatLon{
			Lat: 40.634408,
			Lon: -74.198356,
		},
	}
	d2 := &Request{
		ID: uuid.NewString(),
		PickUp: LatLon{
			Lat: 40.634408,
			Lon: -74.198356,
		},
		DropOff: LatLon{Lat: 40.636916,
			Lon: -74.195995,
		},
	}

	res1, err := batchService.GetEstimates(ctx, from, []*Request{d1, d2})
	require.NoError(t, err)
	require.Len(t, res1, 2)
}

func testToChunks(t *testing.T) {
	requests := []*Request{
		{ID: uuid.NewString()},
		{ID: uuid.NewString()},
		{ID: uuid.NewString()},
		{ID: uuid.NewString()},
		{ID: uuid.NewString()},
		{ID: uuid.NewString()},
		{ID: uuid.NewString()},
		{ID: uuid.NewString()},
		{ID: uuid.NewString()},
		{ID: uuid.NewString()},
		{ID: uuid.NewString()},
		{ID: uuid.NewString()},
	}

	chunks := toChunks(requests, 5)
	assert.Len(t, chunks, 3)
	assert.Len(t, chunks[0], 5)
	assert.Len(t, chunks[1], 5)
	assert.Len(t, chunks[2], 2)
	assert.Len(t, toChunks(requests, 12), 1)
	assert.Len(t, toChunks(requests, 20), 1)

	assert.Len(t, toChunks([]*Request{}, 5), 0)
	assert.Len(t, toChunks([]*Request{}, 0), 0)
	assert.Len(t, toChunks([]*Request{}, 1), 0)
}
