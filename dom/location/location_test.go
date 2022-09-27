package location

import (
	"math/rand"
	"testing"
	"time"

	"github.com/driverscooperative/geosrv/dao"
	"github.com/driverscooperative/geosrv/dom"
	"github.com/driverscooperative/geosrv/pkg/util"

	"github.com/driverscooperative/geosrv/cfg"
	"github.com/stretchr/testify/require"
)

func TestInsertLocation(t *testing.T) {
	dom.WillTest(t, "test", false)
	r := rand.New(rand.NewSource(time.Now().UnixMilli()))

	ctx := cfg.Context()
	for i := 0; i < 100; i++ {
		wloc := genLocation(r)
		require.NoError(t, wloc.Insert(ctx))
	}
}

func TestInsertLocations(t *testing.T) {
	dom.WillTest(t, "test", true)
	r := rand.New(rand.NewSource(time.Now().UnixMilli()))

	insert(t, r, 100)
}

func insert(t *testing.T, r *rand.Rand, max int) (wid dom.SUID, locations []Location) {
	wid = dao.MockUUID()
	ctx := cfg.Context()
	for i := 0; i < max; i++ {
		longitude := util.LongitudeInRange(r, -122.47304848490842, -122.43073395709482)
		latitude := util.LatitudeInRange(r, 37.65046887713942, 37.65617701286946)
		l := &Location{
			Worker:    wid,
			Longitude: longitude,
			Latitude:  latitude,
			Speed:     10,
		}
		err := l.Insert(ctx)
		locations = append(locations, Location{Longitude: longitude, Latitude: latitude})
		require.NoError(t, err)
	}
	return
}

func TestQueryLastLocation(t *testing.T) {
	dom.WillTest(t, "test", true)
	r := rand.New(rand.NewSource(time.Now().UnixMilli()))

	wid, loc := insert(t, r, 100)
	ctx := cfg.Context()
	location, _, err := QueryLast(ctx, wid)
	require.NoError(t, err)
	require.Equal(t, loc[len(loc)-1].Longitude, location.Longitude)
	require.Equal(t, loc[len(loc)-1].Latitude, location.Latitude)
}

func TestQueryAllLocations(t *testing.T) {
	dom.WillTest(t, "test", true)
	r := rand.New(rand.NewSource(time.Now().UnixMilli()))
	max := 100
	wid, loc := insert(t, r, 100)
	locations, err := QueryAll(cfg.Context(), wid, max)
	require.NoError(t, err)
	max = len(locations)
	for i := 0; i < max; i++ {
		require.Equal(t, loc[max-i-1].Longitude, locations[i].Longitude)
		require.Equal(t, loc[max-i-1].Latitude, locations[i].Latitude)
	}
}

func genLocation(r *rand.Rand) *Location {
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
