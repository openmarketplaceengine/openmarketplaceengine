package crossing

import (
	"math/rand"
	"testing"
	"time"

	"github.com/openmarketplaceengine/openmarketplaceengine/internal/detector"

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

	workerID := uuid.NewString()
	tollgateID := toll.ID
	x := newRandomCrossing(r, tollgateID, workerID)
	err = x.Insert(ctx)
	require.NoError(t, err)

	wheres := []Where{{Expr: "worker_id = ?", Args: []interface{}{workerID}}, {Expr: "tollgate_id = ?", Args: []interface{}{tollgateID}}}
	orderBy := []string{"created desc"}
	crossings, err := QueryBy(ctx, wheres, orderBy, 100)
	require.NoError(t, err)
	require.Len(t, crossings, 1)
	require.Equal(t, tollgateID, crossings[0].TollgateID)
	require.Equal(t, workerID, crossings[0].WorkerID)
	require.Less(t, crossings[0].Created.UnixMilli(), time.Now().UnixMilli())
}

func newRandomTollgate(r *rand.Rand, name string) *tollgate.Tollgate {
	id := uuid.NewString()

	return &tollgate.Tollgate{
		ID:   id,
		Name: name,
		BBoxes: &tollgate.BBoxes{
			BBoxes: []*detector.BBox{{
				LonMin: util.LongitudeInRange(r, -122.473048, -122.430733),
				LatMin: util.LatitudeInRange(r, 37.656177, 37.656177),
				LonMax: util.LongitudeInRange(r, -122.473048, -122.430733),
				LatMax: util.LatitudeInRange(r, 37.656177, 37.656177),
			}},
			Required: 2,
		},
		GateLine: &tollgate.GateLine{
			Line: &detector.Line{
				Lon1: util.LongitudeInRange(r, -122.473048, -122.430733),
				Lat1: util.LatitudeInRange(r, 37.656177, 37.656177),
				Lon2: util.LongitudeInRange(r, -122.473048, -122.430733),
				Lat2: util.LatitudeInRange(r, 37.656177, 37.656177),
			},
		},
	}
}

func newRandomCrossing(r *rand.Rand, tollgateID dom.SUID, workerID dom.SUID) *TollgateCrossing {
	return &TollgateCrossing{
		ID:         uuid.NewString(),
		TollgateID: tollgateID,
		WorkerID:   workerID,
		Crossing: Crossing{
			Crossing: detector.Crossing{
				TollgateID: tollgateID,
				WorkerID:   workerID,
				Movement: &detector.Movement{
					SubjectID: "",
					From: &detector.Location{
						Lon: util.LongitudeInRange(r, -122.473048, -122.430733),
						Lat: util.LatitudeInRange(r, 37.656177, 37.656177),
					},
					To: &detector.Location{
						Lon: util.LongitudeInRange(r, -122.473048, -122.430733),
						Lat: util.LatitudeInRange(r, 37.656177, 37.656177),
					},
				},
				Direction: "N",
				Alg:       detector.LineAlg,
			},
		},
	}
}
