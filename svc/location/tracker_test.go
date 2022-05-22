package location

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
	"github.com/openmarketplaceengine/openmarketplaceengine/dao"
	"github.com/openmarketplaceengine/openmarketplaceengine/dom"
	"github.com/openmarketplaceengine/openmarketplaceengine/dom/crossing"
	"github.com/openmarketplaceengine/openmarketplaceengine/dom/tollgate"
	"github.com/openmarketplaceengine/openmarketplaceengine/dom/worker"
	"github.com/openmarketplaceengine/openmarketplaceengine/pkg/detector"
	svcTollgate "github.com/openmarketplaceengine/openmarketplaceengine/svc/tollgate"
	"github.com/stretchr/testify/require"
)

const tollgateID = "tollgate-123"

func TestTracker(t *testing.T) {
	err := cfg.Load()
	require.NoError(t, err)

	dom.WillTest(t, "test", false)
	if !dao.Reds.State.Running() {
		require.NoError(t, dao.Reds.Boot())
	}

	ctx := context.Background()

	err = svcTollgate.Load(ctx)
	require.NoError(t, err)

	_, err = tollgate.CreateIfNotExists(ctx, &tollgate.Tollgate{
		ID:     tollgateID,
		Name:   "TestController2",
		BBoxes: nil,
		GateLine: &tollgate.GateLine{
			Line: &detector.Line{
				Lon1: -74.195995,
				Lat1: 40.636916,
				Lon2: -74.198356,
				Lat2: 40.634408,
			},
		},
	})
	require.NoError(t, err)

	tollgates, err := tollgate.QueryAll(cfg.Context(), 100)
	require.NoError(t, err)

	storeClient := dao.Reds.StoreClient
	//pubSubClient := dao.Reds.PubSubClient
	_d := detector.NewDetector(transformTollgates(tollgates), NewBBoxStorage(storeClient))
	tracker := NewTracker(NewStorage(storeClient), _d)
	require.NoError(t, err)

	t.Run("testDetectCrossing", func(t *testing.T) {
		testDetectCrossing(t, tracker)
	})
}

func testDetectCrossing(t *testing.T, tracker *Tracker) {
	ctx := context.Background()

	//target := make(chan detector.Crossing)
	//channel := crossingChannel(tollgateID)
	//subscribe(t, client, channel, target)

	id := uuid.NewString()
	lon0 := -74.195995
	lat0 := 40.636916
	lon1 := -74.198356
	lat1 := 40.634408

	areaKey := "test-tracker"
	c0, err := tracker.TrackLocation(ctx, areaKey, id, lon0, lat0)

	require.NoError(t, err)
	require.Nil(t, c0)

	c1, err := tracker.TrackLocation(ctx, areaKey, id, lon1, lat1)

	require.NoError(t, err)
	require.NotNil(t, c1)
	require.Equal(t, id, c1.WorkerID)

	//select {
	//case c := <-target:
	//	require.Equal(t, tollgateID, c.TollgateID)
	//	require.Equal(t, detector.Direction("SW"), c.Direction)
	//	require.Equal(t, id, c.WorkerID)
	//	require.InDelta(t, lon1, c.Movement.To.Lon, 0.003)
	//	require.InDelta(t, lat1, c.Movement.To.Lat, 0.003)
	//	break
	//case <-time.After(5 * time.Second):
	//	require.Fail(t, "timeout")
	//}

	wheres := []crossing.Where{{Expr: "worker_id = ?", Args: []interface{}{id}}, {Expr: "tollgate_id = ?", Args: []interface{}{tollgateID}}}
	orderBy := []string{"created desc"}
	retrieved, err := crossing.QueryBy(ctx, wheres, orderBy, 100)
	require.NoError(t, err)
	require.Len(t, retrieved, 1)
	require.Equal(t, tollgateID, retrieved[0].TollgateID)
	require.Equal(t, id, retrieved[0].WorkerID)

	ls, err := worker.ListLocations(ctx, id, 100)
	require.NoError(t, err)
	require.Len(t, ls, 2)
}

//func subscribe(t *testing.T, client *redis.Client, channel string, target chan detector.Crossing) {
//	pubSub := client.PSubscribe(context.Background(), channel)
//	go func() {
//		for m := range pubSub.Channel() {
//			var c detector.Crossing
//			err := json.Unmarshal([]byte(m.Payload), &c)
//			require.NoError(t, err)
//			target <- c
//		}
//	}()
//	//subscribe takes a while
//	time.Sleep(50 * time.Millisecond)
//}

func transformTollgates(tollgates []*tollgate.Tollgate) (result []*detector.Tollgate) {
	for _, t := range tollgates {
		var line *detector.Line
		var bBoxes []*detector.BBox
		var bBoxesRequired int32

		if t.GateLine != nil {
			line = t.GateLine.Line
		}

		if t.BBoxes != nil {
			bBoxes = t.BBoxes.BBoxes
			bBoxesRequired = t.BBoxes.Required
		}

		result = append(result, &detector.Tollgate{
			ID:             t.ID,
			Line:           line,
			BBoxes:         bBoxes,
			BBoxesRequired: bBoxesRequired,
		})
	}
	return
}
