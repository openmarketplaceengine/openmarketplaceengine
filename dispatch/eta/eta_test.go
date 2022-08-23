package eta

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/openmarketplaceengine/openmarketplaceengine/dispatch/job"
	"github.com/stretchr/testify/require"
)

var apiKey = os.Getenv("OME_GOOGLE_API_KEY")

func TestEstimateJobs(t *testing.T) {
	if apiKey == "" {
		t.Skip("OME_GOOGLE_API_KEY env var is not set, skipping.")
	}
	ctx := context.Background()

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

	from := job.LatLon{
		Lon: -74.143650,
		Lat: 40.633650,
	}
	jobs, err := EstimateJobs(ctx, apiKey, from, []*job.Job{job1, job2})
	require.NoError(t, err)

	require.NoError(t, err)
	require.NotNil(t, jobs)

	require.Equal(t, len(jobs), 2)
	require.Len(t, jobs, 2)
	require.Equal(t, jobs[0].Pickup.Address, "JRP3+QJ New York, NY, USA")
	require.Equal(t, jobs[1].Pickup.Address, "JRM2+QM New York, NY, USA")
}
