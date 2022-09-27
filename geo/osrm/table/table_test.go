package table

import (
	"net/http"
	"testing"

	"github.com/driverscooperative/geosrv/geo/osrm"
	"github.com/stretchr/testify/require"
)

func TestTable_1_origin_2_destinations(t *testing.T) {
	request := Request{
		Coordinates: []osrm.LngLat{
			{-73.980052, 40.751753},
			{-73.962662, 40.794156},
		},
		Origins:      []int{0},
		Destinations: []int{0, 1},
		Annotations:  DurationDistance,
	}

	response, err := Table(&http.Client{}, request)
	require.NoError(t, err)
	require.Equal(t, "Ok", response.Code)

	require.Len(t, response.Distances, 1)
	require.Len(t, response.Distances[0], 2)
	require.Len(t, response.Durations, 1)
	require.Len(t, response.Durations[0], 2)

	require.Len(t, response.Destinations, 2)
	require.Len(t, response.Sources, 1)

	require.Greater(t, response.Distances[0][1], float64(0))
	require.Greater(t, response.Durations[0][1], float64(0))
}

func TestTable_1_origin_1_destination(t *testing.T) {
	request := Request{
		Coordinates: []osrm.LngLat{
			{-73.9650115649754, 40.791680675548136},
			{-73.98145413365043, 40.76866089218841},
		},
		Origins:      []int{0},
		Destinations: []int{1},
		Annotations:  DurationDistance,
	}

	response, err := Table(&http.Client{}, request)
	require.NoError(t, err)
	require.Equal(t, "Ok", response.Code)

	require.Len(t, response.Distances, 1)
	require.Len(t, response.Distances[0], 1)
	require.Len(t, response.Durations, 1)
	require.Len(t, response.Durations[0], 1)

	require.Len(t, response.Destinations, 1)
	require.Len(t, response.Sources, 1)

	require.Greater(t, response.Distances[0][0], float64(0))
	require.Greater(t, response.Durations[0][0], float64(0))
}

func TestTable_2_origins_3_destinations(t *testing.T) {
	request := Request{
		Coordinates: []osrm.LngLat{
			{-73.980052, 40.751753},
			{-73.962662, 40.794156},
			{-73.952794, 40.785965},
		},
		Origins:      []int{0, 1},
		Destinations: []int{0, 1, 2},
		Annotations:  DurationDistance,
	}

	response, err := Table(&http.Client{}, request)
	require.NoError(t, err)
	require.Equal(t, "Ok", response.Code)

	require.Len(t, response.Destinations, 3)
	require.Len(t, response.Sources, 2)

	require.Len(t, response.Distances, 2)
	require.Len(t, response.Distances[0], 3)
	require.Len(t, response.Distances[1], 3)
	require.Len(t, response.Durations, 2)
	require.Len(t, response.Durations[0], 3)
	require.Len(t, response.Durations[1], 3)
}
