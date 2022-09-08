package geohash

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestToGeoHash(t *testing.T) {
	lat := 37.656177
	lon := -122.473048
	require.Equal(t, "9q8yjp1g5", ToGeoHash(lat, lon, 9))
	require.Equal(t, "9q8yjp1g", ToGeoHash(lat, lon, 8))
	require.Equal(t, "9q8yjp1", ToGeoHash(lat, lon, 7))
	require.Equal(t, "9q8yjp", ToGeoHash(lat, lon, 6))
}
