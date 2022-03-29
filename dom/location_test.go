package dom

import (
	"testing"
	"time"

	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
	"github.com/stretchr/testify/require"
)

func TestWorkerLocation_Persist(t *testing.T) {
	WillTest(t, "test", false)
	ctx := cfg.Context()
	for i := 0; i < 100; i++ {
		wloc := genWorkerLocation()
		require.NoError(t, wloc.Persist(ctx))
	}
}

//-----------------------------------------------------------------------------

func TestAddWorkerLocation(t *testing.T) {
	WillTest(t, "test", true)
	insWorkerLocations(t, 100)
}

//-----------------------------------------------------------------------------

func insWorkerLocations(t *testing.T, max int) (wid SUID, loc []Coord) {
	wid = numUUID(1, "drv")
	loc = genCoord(max)
	ctx := cfg.Context()
	for i := 0; i < len(loc); i++ {
		c := loc[i]
		err := AddWorkerLocation(ctx, wid, c.Longitude, c.Latitude, time.Now(), mockRange(10, 100))
		require.NoError(t, err)
	}
	return
}

//-----------------------------------------------------------------------------

func TestLastWorkerLocation(t *testing.T) {
	WillTest(t, "test", true)
	wid, loc := insWorkerLocations(t, 100)
	ctx := cfg.Context()
	cor, err := LastWorkerLocation(ctx, wid)
	require.NoError(t, err)
	c := loc[len(loc)-1]
	require.Equal(t, c, cor)
}

//-----------------------------------------------------------------------------

func TestListWorkerLocation(t *testing.T) {
	WillTest(t, "test", true)
	max := 100
	wid, loc := insWorkerLocations(t, max)
	cor, err := ListWorkerLocation(cfg.Context(), wid, max)
	require.NoError(t, err)
	max = len(cor)
	for i := 0; i < max; i++ {
		require.Equal(t, loc[max-i-1], cor[i])
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

//-----------------------------------------------------------------------------

func genCoord(n int) []Coord {
	ary := make([]Coord, n)
	for i := 0; i < n; i++ {
		c := &ary[i]
		c.Longitude = mockCoord()
		c.Latitude = mockCoord()
	}
	return ary
}
