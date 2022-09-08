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

var jobs = []Request{
	{ID: "job0", PickUp: LatLon{Lon: -122.47304848490842, Lat: 37.65617701286946}},
	{ID: "job1", PickUp: LatLon{Lon: -122.43073395709482, Lat: 37.65046887713942}},
	{ID: "job2", PickUp: LatLon{Lon: -122.43536881409672, Lat: 37.64061451520277}},
	{ID: "job3", PickUp: LatLon{Lon: -122.48575142632102, Lat: 37.633953585386486}},
	{ID: "job4", PickUp: LatLon{Lon: -122.47708252711378, Lat: 37.63191440586108}},
	{ID: "job5", PickUp: LatLon{Lon: -122.48025826264276, Lat: 37.68681676281144}},
	{ID: "job6", PickUp: LatLon{Lon: -122.46781281326953, Lat: 37.729188812252616}},
}

var myLocation = LatLon{
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
		err = queue.Enqueue(ctx, areaKey, job)
		require.NoError(t, err)
	}

	t.Run("testGetNearbyRadius5000", func(t *testing.T) {
		testGetNearbyRadius5000(ctx, t, queue)
	})
	t.Run("testGetNearbyRadius3000", func(t *testing.T) {
		testGetNearbyRadius3000(ctx, t, queue)
	})
	t.Run("testPeek", func(t *testing.T) {
		testPeek(ctx, t, queue)
	})

	t.Run("testGetAddedBeforeTime", func(t *testing.T) {
		testGetAddedBeforeTime(ctx, t, queue)
	})
}

func testGetNearbyRadius5000(ctx context.Context, t *testing.T, queue *Queue) {
	result, err := queue.GetNearbyRequests(ctx, areaKey, myLocation, 5000)
	require.NoError(t, err)
	require.Len(t, result, 4)
}

func testGetNearbyRadius3000(ctx context.Context, t *testing.T, queue *Queue) {
	result, err := queue.GetNearbyRequests(ctx, areaKey, myLocation, 3000)

	require.NoError(t, err)
	require.Len(t, result, 2)
}

func testPeek(ctx context.Context, t *testing.T, queue *Queue) {
	start := time.Now()
	jobID := uuid.NewString()

	job := queue.Peek(ctx, areaKey, jobID)
	require.Nil(t, job)

	job = queue.Peek(ctx, areaKey, jobs[0].ID)
	require.NotNil(t, job)
	require.Less(t, job.EnqueuedTime.UnixNano(), start.UnixNano())
}

func testGetAddedBeforeTime(ctx context.Context, t *testing.T, queue *Queue) {
	addedBefore := time.Time{}

	for i, job := range jobs {
		time.Sleep(10 * time.Millisecond)
		err := queue.Enqueue(ctx, areaKey, job)
		require.NoError(t, err)
		if i == 3 {
			addedBefore = time.Now()
		}
	}

	result, err := queue.GetAddedBeforeTime(ctx, areaKey, addedBefore, 5)
	require.NoError(t, err)
	require.Len(t, result, 4)

	job := result[len(result)-1]

	require.Less(t, job.EnqueuedTime.UnixNano(), time.Now().UnixNano())
}
