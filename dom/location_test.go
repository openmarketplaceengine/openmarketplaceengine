package dom

import (
	"testing"

	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
	"github.com/stretchr/testify/require"
)

func TestWorkerLocation_Persist(t *testing.T) {
	WillTest(t, "test", false)
	for i := 0; i < 100; i++ {
		wloc := genWorkerLocation()
		require.NoError(t, wloc.Persist(cfg.Context()))
	}
}

//-----------------------------------------------------------------------------

func genWorkerLocation() *WorkerLocation {
	return &WorkerLocation{
		Worker:    mockUUID(),
		Stamp:     mockStamp(),
		Longitude: mockCoord(),
		Latitude:  mockCoord(),
		Speed:     mockSpeed(),
	}
}
