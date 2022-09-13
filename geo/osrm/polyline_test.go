package osrm

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/twpayne/go-polyline"
)

func TestToPolyline(t *testing.T) {
	coords := [][]float64{
		{40.76568, 73.97624},
		{40.76588, 73.97672},
		{40.76704, 73.97947},
	}
	bytes := polyline.EncodeCoords(coords)
	require.NotEmpty(t, bytes)
	require.Equal(t, "o`ywFonobMg@_BgFeP", string(bytes))
}
