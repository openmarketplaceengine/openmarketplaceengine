package route

import (
	"net/http"
	"testing"

	"github.com/openmarketplaceengine/openmarketplaceengine/geo/osrm"
	"github.com/stretchr/testify/require"
)

func TestRoute(t *testing.T) {
	request := Request{
		Coordinates: []osrm.LngLat{
			{-73.980052, 40.751753},
			{-73.962662, 40.794156},
		},
		Alternatives:     false,
		Steps:            false,
		Annotations:      Distance,
		Geometries:       Geojson,
		Overview:         Simplified,
		ContinueStraight: "false",
		Waypoints:        nil,
	}

	response, err := Routes(&http.Client{}, request)
	require.NoError(t, err)
	require.Equal(t, "Ok", response.Code)

	// bytes, err := json.Marshal(response)
	// require.NoError(t, err)
	// fmt.Printf("%s", bytes)

	require.Len(t, response.Routes, 1)
	require.Len(t, response.Routes[0].Legs, 1)
	require.Len(t, response.Routes[0].Legs[0].Steps, 0)
	require.Greater(t, response.Routes[0].Legs[0].Weight, float64(0))
	require.Greater(t, response.Routes[0].Legs[0].Duration, float64(0))
	require.Greater(t, response.Routes[0].Legs[0].Distance, float64(0))
	require.Equal(t, response.Routes[0].Legs[0].Summary, "")
}

func TestRouteSteps(t *testing.T) {
	request := Request{
		Coordinates: []osrm.LngLat{
			{-73.980052, 40.751753},
			{-73.962662, 40.794156},
		},
		Alternatives:     false,
		Steps:            true,
		Annotations:      Distance,
		Geometries:       Geojson,
		Overview:         Simplified,
		ContinueStraight: "false",
		Waypoints:        nil,
	}

	response, err := Routes(&http.Client{}, request)
	require.NoError(t, err)
	require.Equal(t, "Ok", response.Code)

	// bytes, err := json.Marshal(response)
	// require.NoError(t, err)
	// fmt.Printf("%s", bytes)

	require.Len(t, response.Routes, 1)
	require.Len(t, response.Routes[0].Legs, 1)
	require.Len(t, response.Routes[0].Legs[0].Steps, 7)
}
