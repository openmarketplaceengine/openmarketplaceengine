package storage

import (
	"testing"

	"context"

	"github.com/openmarketplaceengine/openmarketplaceengine/config"
	redisClient "github.com/openmarketplaceengine/openmarketplaceengine/redis/client"
	"github.com/stretchr/testify/require"
)

func TestLocationStorage(t *testing.T) {
	err := config.Read()
	require.NoError(t, err)

	client := redisClient.NewStoreClient()
	require.NotNil(t, client)

	storage := NewLocationStorage(client)

	ctx := context.Background()
	err = storage.GeoAdd(ctx, "italy", "Palermo", 13.361389, 38.115556)
	require.NoError(t, err)
	err = storage.GeoAdd(ctx, "italy", "Catania", 15.087269, 37.502669)
	require.NoError(t, err)

	t.Run("testGeoSearchRadius1", func(t *testing.T) {
		testGeoSearchRadius1(t, storage)
	})
	t.Run("testGeoSearchRadius2", func(t *testing.T) {
		testGeoSearchRadius2(t, storage)
	})
}

func testGeoSearchRadius1(t *testing.T, storage *LocationStorage) {
	ctx := context.Background()
	result, err := storage.GeoSearch(ctx, "italy", 15, 37, 200, "km")

	require.NoError(t, err)
	require.Len(t, result, 2)
	require.Contains(t, result, "Palermo", "Catania")
}

func testGeoSearchRadius2(t *testing.T, storage *LocationStorage) {
	ctx := context.Background()
	result, err := storage.GeoSearch(ctx, "italy", 15, 37, 100, "km")

	require.NoError(t, err)
	require.Len(t, result, 1)
	require.Contains(t, result, "Catania")
}
