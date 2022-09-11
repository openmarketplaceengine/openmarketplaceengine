package estimate

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

var apiKey = os.Getenv("OME_GOOGLE_API_KEY")

func TestEstimates(t *testing.T) {
	if apiKey == "" {
		t.Skip("OME_GOOGLE_API_KEY env var is not set, skipping.")
	}
	ctx := context.Background()

	d1 := &Request{
		ID: uuid.NewString(),
		PickUp: LatLon{
			Lat: 40.636916,
			Lon: -74.195995,
		},
		DropOff: LatLon{
			Lat: 40.634408,
			Lon: -74.198356,
		},
	}
	d2 := &Request{
		ID: uuid.NewString(),
		PickUp: LatLon{
			Lat: 40.634408,
			Lon: -74.198356,
		},
		DropOff: LatLon{
			Lat: 40.636916,
			Lon: -74.195995,
		},
	}

	from := LatLon{
		Lon: -74.143650,
		Lat: 40.633650,
	}
	estimates, err := Estimates(ctx, apiKey, from, []*Request{d1, d2})
	require.NoError(t, err)

	require.NoError(t, err)
	require.NotNil(t, estimates)

	require.Equal(t, len(estimates), 2)
	require.Len(t, estimates, 2)
	require.Equal(t, estimates[0].PickUp.Address, "JRP3+QJ New York, NY, USA")
	require.Equal(t, estimates[1].PickUp.Address, "JRM2+QM New York, NY, USA")
}
