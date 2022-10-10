package osrm

import (
	"net/http"
	"testing"
	"time"

	"github.com/openmarketplaceengine/openmarketplaceengine/geo"
	"github.com/openmarketplaceengine/openmarketplaceengine/geo/distance"
	"github.com/stretchr/testify/require"
)

func TestGetMatrix(t *testing.T) {
	a := geo.LatLng{
		Lat: 40.791680675548136,
		Lng: -73.9650115649754,
	}
	b := geo.LatLng{
		Lat: 40.76866089218841,
		Lng: -73.98145413365043,
	}

	matrix, err := Matrix(&http.Client{}, distance.MatrixPointsInput{
		Origins:      []geo.LatLng{a},
		Destinations: []geo.LatLng{b},
	})
	require.NoError(t, err)
	require.Len(t, matrix.Rows, 1)
	require.Len(t, matrix.Rows[0].Elements, 1)
	require.Greater(t, matrix.Rows[0].Elements[0].Duration, time.Duration(0))
	require.Greater(t, matrix.Rows[0].Elements[0].Distance, 0)
}
