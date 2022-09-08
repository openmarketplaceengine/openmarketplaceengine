package geoqueue

import (
	"context"
	"github.com/google/uuid"
	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
	"github.com/openmarketplaceengine/openmarketplaceengine/dao"
	"github.com/stretchr/testify/require"
	"testing"
)

var members = []Member{
	{ID: "id0", PickUp: LatLon{Lon: -122.47304848490842, Lat: 37.65617701286946}},
	{ID: "id1", PickUp: LatLon{Lon: -122.43073395709482, Lat: 37.65046887713942}},
	{ID: "id2", PickUp: LatLon{Lon: -122.43536881409672, Lat: 37.64061451520277}},
	{ID: "id3", PickUp: LatLon{Lon: -122.48575142632102, Lat: 37.633953585386486}},
	{ID: "id4", PickUp: LatLon{Lon: -122.47708252711378, Lat: 37.63191440586108}},
	{ID: "id5", PickUp: LatLon{Lon: -122.48025826264276, Lat: 37.68681676281144}},
	{ID: "id6", PickUp: LatLon{Lon: -122.46781281326953, Lat: 37.729188812252616}},
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

	queue := NewGeoQueue(dao.Reds.StoreClient)

	ctx := context.Background()

	err = queue.Purge(ctx, areaKey)
	require.NoError(t, err)
	for _, member := range members {
		err = queue.Enqueue(ctx, areaKey, member)
		require.NoError(t, err)
	}

	t.Run("testPeekManyRadius5000", func(t *testing.T) {
		testPeekManyRadius5000(ctx, t, queue)
	})
	t.Run("testPeekManyRadius3000", func(t *testing.T) {
		testPeekManyRadius3000(ctx, t, queue)
	})
	t.Run("testPeek", func(t *testing.T) {
		testPeek(ctx, t, queue)
	})
	t.Run("testDequeue", func(t *testing.T) {
		testDequeue(ctx, t, queue)
	})
}

func testPeekManyRadius5000(ctx context.Context, t *testing.T, queue *GeoQueue) {
	result, err := queue.PeekMany(ctx, areaKey, myLocation, 5000)
	require.NoError(t, err)
	require.Len(t, result, 4)
}

func testPeekManyRadius3000(ctx context.Context, t *testing.T, queue *GeoQueue) {
	result, err := queue.PeekMany(ctx, areaKey, myLocation, 3000)

	require.NoError(t, err)
	require.Len(t, result, 2)
}

func testPeek(ctx context.Context, t *testing.T, queue *GeoQueue) {
	id := uuid.NewString()

	one, err := queue.PeekOne(ctx, areaKey, id)
	require.NoError(t, err)
	require.Nil(t, one)

	one, err = queue.PeekOne(ctx, areaKey, members[0].ID)
	require.NoError(t, err)
	require.NotNil(t, one)
}

func testDequeue(ctx context.Context, t *testing.T, queue *GeoQueue) {
	id := uuid.NewString()
	err := queue.Enqueue(ctx, areaKey, Member{
		ID: id,
		PickUp: LatLon{
			Lon: -122.47304848490842,
			Lat: 37.65617701286946,
		},
	})
	require.NoError(t, err)

	m0, err := queue.PeekOne(ctx, areaKey, id)
	require.NoError(t, err)
	require.NotNil(t, m0)

	m1, err := queue.Dequeue(ctx, areaKey, id)
	require.NoError(t, err)
	require.Equal(t, m0, m1)

	m2, _ := queue.PeekOne(ctx, areaKey, id)
	require.Nil(t, m2)
}
