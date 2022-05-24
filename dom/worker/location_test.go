package worker

import (
	"math/rand"
	"testing"
	"time"

	"github.com/openmarketplaceengine/openmarketplaceengine/dao"
	"github.com/openmarketplaceengine/openmarketplaceengine/dom"
	"github.com/openmarketplaceengine/openmarketplaceengine/pkg/util"

	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
	"github.com/stretchr/testify/require"
)

func TestWorkerLocation_Persist(t *testing.T) {
	dom.WillTest(t, "test", false)
	r := rand.New(rand.NewSource(time.Now().UnixMilli()))

	ctx := cfg.Context()
	for i := 0; i < 100; i++ {
		wloc := genWorkerLocation(r)
		require.NoError(t, wloc.Persist(ctx))
	}
}

//-----------------------------------------------------------------------------

func TestAddWorkerLocation(t *testing.T) {
	dom.WillTest(t, "test", true)
	r := rand.New(rand.NewSource(time.Now().UnixMilli()))

	insWorkerLocations(t, r, 100)
}

//-----------------------------------------------------------------------------

func insWorkerLocations(t *testing.T, r *rand.Rand, max int) (wid dom.SUID, locations []Location) {
	wid = dao.MockUUID()
	ctx := cfg.Context()
	for i := 0; i < max; i++ {
		longitude := util.LongitudeInRange(r, -122.47304848490842, -122.43073395709482)
		latitude := util.LatitudeInRange(r, 37.65046887713942, 37.65617701286946)
		err := AddLocation(ctx, wid, longitude, latitude, time.Now(), randRange(10, 80))
		locations = append(locations, Location{Longitude: longitude, Latitude: latitude})
		require.NoError(t, err)
	}
	return
}

//-----------------------------------------------------------------------------

func TestLastWorkerLocation(t *testing.T) {
	dom.WillTest(t, "test", true)
	r := rand.New(rand.NewSource(time.Now().UnixMilli()))

	wid, loc := insWorkerLocations(t, r, 100)
	ctx := cfg.Context()
	location, _, err := LastLocation(ctx, wid)
	require.NoError(t, err)
	require.Equal(t, loc[len(loc)-1].Longitude, location.Longitude)
	require.Equal(t, loc[len(loc)-1].Latitude, location.Latitude)
}

//-----------------------------------------------------------------------------

func TestListWorkerLocation(t *testing.T) {
	dom.WillTest(t, "test", true)
	r := rand.New(rand.NewSource(time.Now().UnixMilli()))
	max := 100
	wid, loc := insWorkerLocations(t, r, 100)
	locations, err := ListLocations(cfg.Context(), wid, max)
	require.NoError(t, err)
	max = len(locations)
	for i := 0; i < max; i++ {
		require.Equal(t, loc[max-i-1].Longitude, locations[i].Longitude)
		require.Equal(t, loc[max-i-1].Latitude, locations[i].Latitude)
	}
}

//-----------------------------------------------------------------------------

func genWorkerLocation(r *rand.Rand) *Location {
	stamp := dom.Time{}
	stamp.Now()
	return &Location{
		Worker:    dao.MockUUID(),
		Stamp:     stamp,
		Longitude: util.LongitudeInRange(r, -122.47304848490842, -122.43073395709482),
		Latitude:  util.LatitudeInRange(r, 37.65046887713942, 37.65617701286946),
		Speed:     10,
	}
}
