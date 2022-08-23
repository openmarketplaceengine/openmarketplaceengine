package queue

import (
	"context"
	"testing"
	"time"

	"github.com/openmarketplaceengine/openmarketplaceengine/dao"

	"github.com/google/uuid"
	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
	"github.com/stretchr/testify/require"
)

var jobs = []Job{
	{ID: "job0", PickupLocation: Location{Lon: -122.47304848490842, Lat: 37.65617701286946}},
	{ID: "job1", PickupLocation: Location{Lon: -122.43073395709482, Lat: 37.65046887713942}},
	{ID: "job2", PickupLocation: Location{Lon: -122.43536881409672, Lat: 37.64061451520277}},
	{ID: "job3", PickupLocation: Location{Lon: -122.48575142632102, Lat: 37.633953585386486}},
	{ID: "job4", PickupLocation: Location{Lon: -122.47708252711378, Lat: 37.63191440586108}},
	{ID: "job5", PickupLocation: Location{Lon: -122.48025826264276, Lat: 37.68681676281144}},
	{ID: "job6", PickupLocation: Location{Lon: -122.46781281326953, Lat: 37.729188812252616}},
}

var myLocation = Location{
	Lon: -122.45476654908023,
	Lat: 37.6777824094095,
}

const areaKey = "dispatch_test"

func TestQueue(t *testing.T) {
	err := cfg.Load()
	require.NoError(t, err)

	if !dao.Reds.State.Running() {
		require.NoError(t, dao.Reds.Boot())
	}

	queue := NewQueue(dao.Reds.StoreClient)

	ctx := context.Background()

	for _, job := range jobs {
		err = queue.Add(ctx, areaKey, job)
		require.NoError(t, err)
	}

	t.Run("testGetNearbyRadius5", func(t *testing.T) {
		testGetNearbyRadius5(ctx, t, queue)
	})
	t.Run("testGetNearbyRadius3", func(t *testing.T) {
		testGetNearbyRadius3(ctx, t, queue)
	})
	t.Run("testGetJob", func(t *testing.T) {
		testGetJob(ctx, t, queue)
	})

	t.Run("testGetAddedBeforeTimeJobs", func(t *testing.T) {
		testGetAddedBeforeTimeJobs(ctx, t, queue)
	})
}

func testGetNearbyRadius5(ctx context.Context, t *testing.T, queue *Queue) {
	result, err := queue.GetNearbyJobs(ctx, areaKey, myLocation, 5, "km")
	require.NoError(t, err)
	require.Len(t, result, 4)
}

func testGetNearbyRadius3(ctx context.Context, t *testing.T, queue *Queue) {
	result, err := queue.GetNearbyJobs(ctx, areaKey, myLocation, 3, "km")

	require.NoError(t, err)
	require.Len(t, result, 2)
}

func testGetJob(ctx context.Context, t *testing.T, queue *Queue) {
	start := time.Now()
	jobID := uuid.NewString()

	job := queue.GetJob(ctx, areaKey, jobID)
	require.Nil(t, job)

	job = queue.GetJob(ctx, areaKey, jobs[0].ID)
	require.NotNil(t, job)
	require.Less(t, job.EnqueuedTime.UnixNano(), start.UnixNano())
}

func testGetAddedBeforeTimeJobs(ctx context.Context, t *testing.T, queue *Queue) {
	addedBefore := time.Time{}

	for i, job := range jobs {
		time.Sleep(10 * time.Millisecond)
		err := queue.Add(ctx, areaKey, job)
		require.NoError(t, err)
		if i == 3 {
			addedBefore = time.Now()
		}
	}

	result, err := queue.GetAddedBeforeTimeJobs(ctx, areaKey, addedBefore, 5)
	require.NoError(t, err)
	require.Len(t, result, 4)

	job := result[len(result)-1]

	require.Less(t, job.EnqueuedTime.UnixNano(), time.Now().UnixNano())
}
