package google

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/openmarketplaceengine/openmarketplaceengine/dispatch/job"
	"github.com/stretchr/testify/require"
	"googlemaps.github.io/maps"
)

var apiKey = os.Getenv("OME_GOOGLE_API_KEY")

func TestGetMatrix(t *testing.T) {
	if apiKey == "" {
		t.Skip("OME_GOOGLE_API_KEY env var is not set, skipping.")
	}

	t.Run("testGetMatrix", func(t *testing.T) {
		testGetMatrix(t)
	})

	t.Run("testGetMatrixFromPlaces", func(t *testing.T) {
		testGetMatrixFromPlaces(t)
	})
}

func testGetMatrix(t *testing.T) {
	a := job.LatLon{
		Lat: 40.791680675548136,
		Lon: -73.9650115649754,
	}
	b := job.LatLon{
		Lat: 40.76866089218841,
		Lon: -73.98145413365043,
	}

	ctx := context.Background()
	client, err := maps.NewClient(maps.WithAPIKey(apiKey))
	require.NoError(t, err)
	matrix, err := Matrix(ctx, client, MatrixPointsInput{
		Origins:      []job.LatLon{a},
		Destinations: []job.LatLon{b},
	})
	require.NoError(t, err)
	require.Len(t, matrix.Rows, 1)
	require.Len(t, matrix.Rows[0].Elements, 1)
	require.Greater(t, matrix.Rows[0].Elements[0].Duration, time.Duration(0))
	require.Greater(t, matrix.Rows[0].Elements[0].Distance, 0)
	require.Equal(t, matrix.OriginAddresses[0], "96 St, New York, NY 10025, USA")
	require.Equal(t, matrix.DestinationAddresses[0], "59 St-Columbus Circle, Columbus Cir, New York, NY 10023, USA")
}

func testGetMatrixFromPlaces(t *testing.T) {
	ctx := context.Background()
	client, err := maps.NewClient(maps.WithAPIKey(apiKey))
	require.NoError(t, err)
	matrix, err := MatrixFromPlaces(ctx, client, MatrixPlacesInput{
		Origins:      []string{"ChIJ87a7BJ5YwokR4TLbUoQMB1s"},
		Destinations: []string{"ChIJVZfjPvZYwokR-sLEBmjjniw"},
	})
	require.NoError(t, err)
	require.Len(t, matrix.Rows, 1)
	require.Len(t, matrix.Rows[0].Elements, 1)
	require.Greater(t, matrix.Rows[0].Elements[0].Duration, time.Duration(0))
	require.Greater(t, matrix.Rows[0].Elements[0].Distance, 0)
	require.Equal(t, matrix.OriginAddresses[0], "96 St, New York, NY 10025, USA")
	require.Equal(t, matrix.DestinationAddresses[0], "59 St-Columbus Circle, Columbus Cir, New York, NY 10023, USA")
}
