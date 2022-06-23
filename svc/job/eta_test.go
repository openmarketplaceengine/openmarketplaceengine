package job

import (
	"context"
	"fmt"
	"sort"
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
		ID:        uuid.NewString(),
		WorkerID:  "",
		State:     state,
		PickupLat: 40.636916,
		PickupLon: -74.195995,
	}
	job2 := &job.Job{
		ID:        uuid.NewString(),
		WorkerID:  "",
		State:     state,
		PickupLat: 40.634408,
		PickupLon: -74.198356,
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
	jobs, err := PickupDistanceEstimatedJobs(from, toJobs)

	require.NoError(t, err)
	require.NotNil(t, jobs)
	require.Equal(t, jobs.OriginAddress, "64 Innis St, Staten Island, NY 10302, USA")

	require.Equal(t, len(jobs.Jobs), len(toJobs))
	require.Len(t, jobs.Jobs, 2)
	require.Equal(t, jobs.Jobs[0].Address, "JRP3+QJ New York, NY, USA")
	require.Equal(t, jobs.Jobs[1].Address, "JRM2+QM New York, NY, USA")

	sorted0 := sort.SliceIsSorted(jobs.Jobs, func(i, j int) bool {
		return jobs.Jobs[i].Distance < jobs.Jobs[j].Distance
	})
	require.True(t, sorted0)
}
