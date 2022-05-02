package tollgate

import (
	"math/rand"
	"testing"
	"time"

	"github.com/openmarketplaceengine/openmarketplaceengine/pkg/detector"

	"github.com/google/uuid"
	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
	"github.com/openmarketplaceengine/openmarketplaceengine/dao"
	"github.com/openmarketplaceengine/openmarketplaceengine/dom"
	"github.com/openmarketplaceengine/openmarketplaceengine/pkg/util"
	"github.com/stretchr/testify/require"
)

func TestTollgateModel(t *testing.T) {
	dom.WillTest(t, "test", false)
	ctx := cfg.Context()

	r := rand.New(rand.NewSource(time.Now().UnixMilli()))

	t.Run("testCreate", func(t *testing.T) {
		testCreate(ctx, t, r)
	})

	t.Run("testCreateIfNotExists", func(t *testing.T) {
		testCreateIfNotExists(ctx, t, r)
	})

	t.Run("testUpdate", func(t *testing.T) {
		testUpdate(ctx, t, r)
	})

	t.Run("testQuery", func(t *testing.T) {
		testQuery(ctx, t, r)
	})
}

func testCreate(ctx dom.Context, t *testing.T, r *rand.Rand) {
	toll := newRandomTollgate(r, "testCreate")

	err := toll.Insert(ctx)
	require.NoError(t, err)

	one, err := QueryOne(ctx, toll.ID)
	require.NoError(t, err)
	require.Equal(t, toll.Name, one.Name)
	require.Equal(t, toll.BBoxes, one.BBoxes)
	require.Equal(t, toll.GateLine, one.GateLine)
	require.Less(t, one.Created.UnixMilli(), time.Now().UnixMilli())
	require.Greater(t, one.Created.UnixMilli(), int64(0))
	require.Equal(t, dao.Time{}, one.Updated)
}

func testCreateIfNotExists(ctx dom.Context, t *testing.T, r *rand.Rand) {
	toll := newRandomTollgate(r, "testCreate")

	created, err := CreateIfNotExists(ctx, toll)
	require.NoError(t, err)
	require.True(t, created)

	created, err = CreateIfNotExists(ctx, toll)
	require.NoError(t, err)
	require.False(t, created)
}

func testUpdate(ctx dom.Context, t *testing.T, r *rand.Rand) {
	toll := newRandomTollgate(r, "testUpdate")

	err := toll.Insert(ctx)
	require.NoError(t, err)

	newName := uuid.NewString()
	toll.Name = newName
	err = toll.Update(ctx)
	require.NoError(t, err)
	one, err := QueryOne(ctx, toll.ID)
	require.NoError(t, err)
	require.Equal(t, newName, one.Name)
	require.Less(t, one.Updated.UnixMilli(), time.Now().UnixMilli())
	require.Greater(t, one.Updated.UnixMilli(), int64(0))
}

func testQuery(ctx dom.Context, t *testing.T, r *rand.Rand) {
	err := deleteAll(ctx)
	require.NoError(t, err)

	toll1 := newRandomTollgate(r, "testQuery")
	toll2 := newRandomTollgate(r, "testQuery")

	err = toll1.Insert(ctx)
	require.NoError(t, err)

	err = toll2.Insert(ctx)
	require.NoError(t, err)

	all, err := QueryAll(ctx, 2)
	require.NoError(t, err)
	require.Len(t, all, 2)
	require.Equal(t, toll1.ID, all[1].ID)
	require.Equal(t, toll1.BBoxes, all[1].BBoxes)
}

func newRandomTollgate(r *rand.Rand, name string) *Tollgate {
	id := uuid.NewString()

	return &Tollgate{
		ID:   id,
		Name: name,
		BBoxes: &BBoxes{
			BBoxes: []*detector.BBox{{
				LonMin: util.LongitudeInRange(r, -122.47304848490842, -122.43073395709482),
				LatMin: util.LatitudeInRange(r, 37.65046887713942, 37.65617701286946),
				LonMax: util.LongitudeInRange(r, -122.47304848490842, -122.43073395709482),
				LatMax: util.LatitudeInRange(r, 37.65046887713942, 37.65617701286946),
			}},
			Required: 2,
		},
		GateLine: &GateLine{
			Line: &detector.Line{
				Lon1: util.LongitudeInRange(r, -122.47304848490842, -122.43073395709482),
				Lat1: util.LatitudeInRange(r, 37.65046887713942, 37.65617701286946),
				Lon2: util.LongitudeInRange(r, -122.47304848490842, -122.43073395709482),
				Lat2: util.LatitudeInRange(r, 37.65046887713942, 37.65617701286946),
			},
		},
	}
}
