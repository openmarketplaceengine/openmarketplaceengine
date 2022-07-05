package job

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
	"github.com/openmarketplaceengine/openmarketplaceengine/dom"
	"github.com/openmarketplaceengine/openmarketplaceengine/dom/job"
	"github.com/openmarketplaceengine/openmarketplaceengine/svc/location"
	"github.com/stretchr/testify/require"
)

func TestEstimatedJobs(t *testing.T) {
	ctx := context.Background()
	dom.WillTest(t, "test", false)
	err := cfg.Load()
	require.NoError(t, err)

	workerID := uuid.NewString()

	state := fmt.Sprintf("AVAILABLE-%s", workerID)
	job1 := &job.Job{
		ID:         uuid.NewString(),
		WorkerID:   "",
		State:      state,
		PickupLat:  40.636916,
		PickupLon:  -74.195995,
		DropoffLat: 40.634408,
		DropoffLon: -74.198356,
	}
	job2 := &job.Job{
		ID:         uuid.NewString(),
		WorkerID:   "",
		State:      state,
		PickupLat:  40.634408,
		PickupLon:  -74.198356,
		DropoffLat: 40.636916,
		DropoffLon: -74.195995,
	}

	for _, j := range []*job.Job{job1, job2} {
		_, _, err = j.Upsert(ctx)
		require.NoError(t, err)
	}

	from := &location.LastLocation{
		WorkerID:     workerID,
		Longitude:    -74.143650,
		Latitude:     40.633650,
		LastSeenTime: time.Time{},
	}
	toJobs, err := job.QueryByPickupDistance(ctx, from.Longitude, from.Latitude, state, 10000, 25)
	require.NoError(t, err)
	jobs, err := estimateJobs(ctx, from, toJobs)

	require.NoError(t, err)
	require.NotNil(t, jobs)

	require.Equal(t, len(jobs), len(toJobs))
	require.Len(t, jobs, 2)
	require.Equal(t, jobs[0].PickupLocation.Address, "JRP3+QJ New York, NY, USA")
	require.Equal(t, jobs[1].PickupLocation.Address, "JRM2+QM New York, NY, USA")
}
