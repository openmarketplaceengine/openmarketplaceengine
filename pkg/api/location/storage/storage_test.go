package storage

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
	redisClient "github.com/openmarketplaceengine/openmarketplaceengine/redis/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var locations = []*Location{
	{WorkerID: "user0", Longitude: -122.47304848490842, Latitude: 37.65617701286946},
	{WorkerID: "user1", Longitude: -122.43073395709482, Latitude: 37.65046887713942},
	{WorkerID: "user2", Longitude: -122.43536881409672, Latitude: 37.64061451520277},
	{WorkerID: "user3", Longitude: -122.48575142632102, Latitude: 37.633953585386486},
	{WorkerID: "user4", Longitude: -122.47708252711378, Latitude: 37.63191440586108},
	{WorkerID: "user5", Longitude: -122.48025826264276, Latitude: 37.68681676281144},
	{WorkerID: "user6", Longitude: -122.46781281326953, Latitude: 37.729188812252616},
}

var myLocation = Location{
	WorkerID:  "userMe",
	Longitude: -122.45476654908023,
	Latitude:  37.6777824094095,
}

const areaKey = "san_fran"

func TestLocationStorage(t *testing.T) {
	err := cfg.Load()
	require.NoError(t, err)

	client := redisClient.NewStoreClient()
	require.NotNil(t, client)

	storage := New(client)

	ctx := context.Background()

	for _, loc := range locations {
		err = storage.Update(ctx, areaKey, loc)
		require.NoError(t, err)
	}

	t.Run("testGeoSearchRadius5", func(t *testing.T) {
		testGeoSearchRadius5(ctx, t, storage)
	})
	t.Run("testGeoSearchRadius3", func(t *testing.T) {
		testGeoSearchRadius3(ctx, t, storage)
	})
	t.Run("testForgetLocation", func(t *testing.T) {
		testForgetLocation(ctx, t, storage)
	})
	t.Run("testRemoveExpiredLocations", func(t *testing.T) {
		testRemoveExpiredLocations(ctx, t, storage)
	})
	t.Run("testRangeLocationsLastSeen", func(t *testing.T) {
		testRangeLocationsLastSeen(ctx, t, storage)
	})
	t.Run("testLastLocation", func(t *testing.T) {
		testLastLocation(ctx, t, storage)
	})
}

func testGeoSearchRadius5(ctx context.Context, t *testing.T, storage *Storage) {
	result, err := storage.RangeLocations(ctx, areaKey, myLocation.Longitude, myLocation.Latitude, 5, "km")
	require.NoError(t, err)
	require.Len(t, result, 4)
}

func testGeoSearchRadius3(ctx context.Context, t *testing.T, storage *Storage) {
	result, err := storage.RangeLocations(ctx, areaKey, myLocation.Longitude, myLocation.Latitude, 3, "km")

	require.NoError(t, err)
	require.Len(t, result, 2)
}

func testForgetLocation(ctx context.Context, t *testing.T, storage *Storage) {
	result, err := storage.RangeLocations(ctx, areaKey, myLocation.Longitude, myLocation.Latitude, 3, "km")
	require.NoError(t, err)
	require.Len(t, result, 2)
	assert.Equal(t, "user5", result[0].WorkerID)
	assert.Equal(t, "user0", result[1].WorkerID)

	err = storage.ForgetLocation(ctx, areaKey, "user0")
	require.NoError(t, err)

	result, err = storage.RangeLocations(ctx, areaKey, myLocation.Longitude, myLocation.Latitude, 3, "km")
	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, "user5", result[0].WorkerID)
}

func testRemoveExpiredLocations(ctx context.Context, t *testing.T, storage *Storage) {
	var start time.Time
	for i, loc := range locations {
		err := storage.Update(ctx, areaKey, loc)
		require.NoError(t, err)
		time.Sleep(10 * time.Millisecond)
		if i == len(locations)/2 {
			start = time.Now()
		}
	}

	result1, err := storage.RangeLocations(ctx, areaKey, myLocation.Longitude, myLocation.Latitude, 8, "km")
	require.NoError(t, err)
	require.Len(t, result1, 7)

	err = storage.RemoveExpiredLocations(ctx, areaKey, start)
	require.NoError(t, err)

	result2, err := storage.RangeLocations(ctx, areaKey, myLocation.Longitude, myLocation.Latitude, 8, "km")
	require.NoError(t, err)
	require.Less(t, len(result2), len(result1))
}

func testRangeLocationsLastSeen(ctx context.Context, t *testing.T, storage *Storage) {
	start := time.Now()
	time.Sleep(10 * time.Millisecond)

	for _, loc := range locations {
		err := storage.Update(ctx, areaKey, loc)
		require.NoError(t, err)
	}

	result, err := storage.RangeLocations(ctx, areaKey, myLocation.Longitude, myLocation.Latitude, 5, "km")
	require.NoError(t, err)
	require.Len(t, result, 4)

	loc := result[0]

	require.Greater(t, loc.LastSeenTime.UnixMilli(), start.UnixMilli())
	require.Less(t, loc.LastSeenTime.UnixMilli(), time.Now().UnixMilli())
}

func testLastLocation(ctx context.Context, t *testing.T, storage *Storage) {
	workerID := uuid.NewString()

	start := time.Now()
	time.Sleep(5 * time.Millisecond)

	err := storage.Update(ctx, areaKey, &Location{
		WorkerID:  workerID,
		Longitude: 13,
		Latitude:  14,
	})
	require.NoError(t, err)

	time.Sleep(5 * time.Millisecond)
	loc := storage.LastLocation(ctx, areaKey, workerID)

	require.Less(t, loc.LastSeenTime.UnixMilli(), time.Now().UnixMilli())
	require.Greater(t, loc.LastSeenTime.UnixMilli(), start.UnixMilli())
}
