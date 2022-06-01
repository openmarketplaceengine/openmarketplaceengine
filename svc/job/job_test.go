package job

import (
	"context"
	"github.com/google/uuid"
	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
	"github.com/openmarketplaceengine/openmarketplaceengine/dao"
	"github.com/openmarketplaceengine/openmarketplaceengine/dom"
	"github.com/openmarketplaceengine/openmarketplaceengine/dom/job"
	"github.com/openmarketplaceengine/openmarketplaceengine/pkg/detector"
	"github.com/openmarketplaceengine/openmarketplaceengine/svc/location"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestService(t *testing.T) {
	err := cfg.Load()
	require.NoError(t, err)

	dom.WillTest(t, "test", false)
	if !dao.Reds.State.Running() {
		require.NoError(t, dao.Reds.Boot())
	}

	ctx := context.Background()

	err = job.DeleteAll(ctx)
	require.NoError(t, err)

	fromLat := 40.633650
	fromLon := -74.143650

	job1 := &job.Job{
		ID:        uuid.NewString(),
		WorkerID:  "",
		State:     "AVAILABLE",
		PickupLat: 40.636916,
		PickupLon: -74.195995,
	}
	job2 := &job.Job{
		ID:        uuid.NewString(),
		WorkerID:  "",
		State:     "AVAILABLE",
		PickupLat: 40.634408,
		PickupLon: -74.198356,
	}

	for _, j := range []*job.Job{job1, job2} {
		_, _, err = j.Upsert(ctx)
		require.NoError(t, err)
	}

	storeClient := dao.Reds.StoreClient
	_d := detector.NewDetector([]*detector.Tollgate{}, nil)
	tracker := location.NewTracker(location.NewStorage(storeClient), _d)
	require.NoError(t, err)
	s := NewService(tracker)

	t.Run("testGetAvailableJobs", func(t *testing.T) {
		testGetAvailableJobs(t, tracker, s, fromLon, fromLat)
	})

	t.Run("testGetAvailableJobsNotTracked", func(t *testing.T) {
		testGetAvailableJobsNotTracked(t, s)
	})
}

func testGetAvailableJobs(t *testing.T, tracker *location.Tracker, service *Service, fromLon float64, fromLat float64) {
	ctx := context.Background()

	id := uuid.NewString()
	areaKey := "test-tracker"
	_, err := tracker.TrackLocation(ctx, areaKey, id, fromLon, fromLat)
	require.NoError(t, err)

	jobs0, err := service.GetAvailableJobs(ctx, areaKey, id, 10.0, 100)
	require.NoError(t, err)
	require.Len(t, jobs0, 2)
}

func testGetAvailableJobsNotTracked(t *testing.T, service *Service) {
	ctx := context.Background()

	id := "16eb5627-ff8e-4b35-916a-4d14191d8229"
	areaKey := "test-tracker"

	_, err := service.GetAvailableJobs(ctx, areaKey, id, 10.0, 100)
	require.EqualError(t, err, "location of worker is not known")
}
