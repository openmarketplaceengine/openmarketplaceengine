package osrm

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/openmarketplaceengine/openmarketplaceengine/geo"
	"github.com/stretchr/testify/require"
)

func TestReverseGeocode(t *testing.T) {
	g := NewGeocoder(new(http.Client))
	ctx := context.Background()
	out, err := g.ReverseGeocode(ctx, geo.LatLng{Lat: -37.813611, Lng: 144.963056})
	require.NoError(t, err)
	fmt.Println(out)
}
