package crossing

import (
	"math/rand"
	"testing"
	"time"

	"github.com/openmarketplaceengine/openmarketplaceengine/internal/service/tollgate/model"

	"github.com/openmarketplaceengine/openmarketplaceengine/internal/service/tollgate"

	"github.com/google/uuid"
	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
	"github.com/openmarketplaceengine/openmarketplaceengine/dom"
	"github.com/openmarketplaceengine/openmarketplaceengine/internal/util"
	"github.com/stretchr/testify/require"
)

func TestTollgateCrossing(t *testing.T) {
	dom.WillTest(t, "test", false)
	ctx := cfg.Context()

	r := rand.New(rand.NewSource(time.Now().UnixMilli()))
	err := deleteAll(ctx)
	require.NoError(t, err)
	t.Run("testCreate", func(t *testing.T) {
		testCreate(ctx, t, r)
	})
}

func testCreate(ctx dom.Context, t *testing.T, r *rand.Rand) {
	toll := newRandomTollgate(r, "testCreate")

	err := toll.Insert(ctx)
	require.NoError(t, err)

	driverID := uuid.NewString()
	tollgateID := toll.ID
	x := newRandomCrossing(r, tollgateID, driverID)
	err = x.Insert(ctx)
	require.NoError(t, err)

	wheres := []Where{{Expr: "driver_id = ?", Args: []interface{}{driverID}}, {Expr: "tollgate_id = ?", Args: []interface{}{tollgateID}}}
	orderBy := []string{"created_at desc"}
	crossings, err := QueryBy(ctx, wheres, orderBy, 100)
	require.NoError(t, err)
	require.Len(t, crossings, 1)
	require.Equal(t, tollgateID, crossings[0].TollgateID)
	require.Equal(t, driverID, crossings[0].DriverID)
	require.Less(t, crossings[0].CreatedAt.UnixMilli(), time.Now().UnixMilli())
}

func newRandomTollgate(r *rand.Rand, name string) *model.Tollgate {
	id := uuid.NewString()

	return &model.Tollgate{
		ID:   id,
		Name: name,
		BBoxes: &model.BBoxes{
			BBoxes: []*tollgate.BBox{{
				LonMin: util.LongitudeInRange(r, -122.47304848490842, -122.43073395709482),
				LatMin: util.LatitudeInRange(r, 37.65046887713942, 37.65617701286946),
				LonMax: util.LongitudeInRange(r, -122.47304848490842, -122.43073395709482),
				LatMax: util.LatitudeInRange(r, 37.65046887713942, 37.65617701286946),
			}},
			Required: 2,
		},
		GateLine: &model.GateLine{
			Line: tollgate.Line{
				Lon1: util.LongitudeInRange(r, -122.47304848490842, -122.43073395709482),
				Lat1: util.LatitudeInRange(r, 37.65046887713942, 37.65617701286946),
				Lon2: util.LongitudeInRange(r, -122.47304848490842, -122.43073395709482),
				Lat2: util.LatitudeInRange(r, 37.65046887713942, 37.65617701286946),
			},
		},
	}
}

func newRandomCrossing(r *rand.Rand, tollgateID dom.SUID, driverID dom.SUID) *TollgateCrossing {
	return &TollgateCrossing{
		ID:         uuid.NewString(),
		TollgateID: tollgateID,
		DriverID:   driverID,
		Crossing: Crossing{
			Crossing: tollgate.Crossing{
				TollgateID: tollgateID,
				DriverID:   driverID,
				Location: &tollgate.Location{
					Lon: util.LongitudeInRange(r, -122.47304848490842, -122.43073395709482),
					Lat: util.LatitudeInRange(r, 37.65046887713942, 37.65617701286946),
				},
				Direction: "N",
				Alg:       0,
			},
		},
	}
}
