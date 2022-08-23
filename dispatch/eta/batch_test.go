package eta

import (
	"testing"

	"github.com/openmarketplaceengine/openmarketplaceengine/dispatch/job"

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

	t.Run("testGetEstimatedJobsPickupToDropOff", func(t *testing.T) {
		testGetEstimatedJobsPickupToDropOff(t, s)
	})

	t.Run("testToChunks", func(t *testing.T) {
		testToChunks(t)
	})
}

func testGetEstimatedJobsPickupToDropOff(t *testing.T, batchService *BatchService) {
	ctx := cfg.Context()

	from := job.LatLon{
		Lat: 40.633650,
		Lon: -74.143650,
	}

	job1 := &job.Job{
		ID: uuid.NewString(),
		Pickup: job.LatLon{
			Lat: 40.636916,
			Lon: -74.195995,
		},
		DropOff: job.LatLon{
			Lat: 40.634408,
			Lon: -74.198356,
		},
	}
	job2 := &job.Job{
		ID: uuid.NewString(),
		Pickup: job.LatLon{
			Lat: 40.634408,
			Lon: -74.198356,
		},
		DropOff: job.LatLon{Lat: 40.636916,
			Lon: -74.195995,
		},
	}

	res1, err := batchService.GetEstimatedJobs(ctx, from, []*job.Job{job1, job2})
	require.NoError(t, err)
	require.Len(t, res1, 2)
}

func testToChunks(t *testing.T) {
	jobs := []*job.Job{
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

	chunks := toChunks(jobs, 5)
	assert.Len(t, chunks, 3)
	assert.Len(t, chunks[0], 5)
	assert.Len(t, chunks[1], 5)
	assert.Len(t, chunks[2], 2)
	assert.Len(t, toChunks(jobs, 12), 1)
	assert.Len(t, toChunks(jobs, 20), 1)

	assert.Len(t, toChunks([]*job.Job{}, 5), 0)
	assert.Len(t, toChunks([]*job.Job{}, 0), 0)
	assert.Len(t, toChunks([]*job.Job{}, 1), 0)
}
