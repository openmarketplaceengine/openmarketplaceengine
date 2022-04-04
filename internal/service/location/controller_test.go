package location

import (
	"context"
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	"github.com/openmarketplaceengine/openmarketplaceengine/dom"
	"github.com/openmarketplaceengine/openmarketplaceengine/internal/service/tollgate/bbox"
	"github.com/openmarketplaceengine/openmarketplaceengine/internal/service/tollgate/model"

	"github.com/openmarketplaceengine/openmarketplaceengine/internal/service/tollgate/detector"

	"github.com/google/uuid"
	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
	locationV1beta1 "github.com/openmarketplaceengine/openmarketplaceengine/internal/omeapi/location/v1beta1"
	"github.com/openmarketplaceengine/openmarketplaceengine/internal/service/tollgate"
	redisClient "github.com/openmarketplaceengine/openmarketplaceengine/redis/client"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const areaKey = "san_fran"
const tollgateID = "tollgate-123"

func TestController(t *testing.T) {
	err := cfg.Load()
	require.NoError(t, err)

	dom.WillTest(t, "test", false)
	ctx := context.Background()

	err = model.CreateIfNotExists(ctx, &model.Tollgate{
		ID:     tollgateID,
		Name:   "TestController",
		BBoxes: nil,
		GateLine: &model.GateLine{
			Line: tollgate.Line{
				Lon1: -79.870262,
				Lat1: 41.198497,
				Lon2: -79.870218,
				Lat2: 41.200268,
			},
		},
	})
	require.NoError(t, err)

	d, err := detector.NewDetector(ctx, bbox.NewStorage(redisClient.NewStoreClient()))
	require.NoError(t, err)

	conn, err := grpc.DialContext(ctx, "", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(dialer(d)))
	defer func(conn *grpc.ClientConn) {
		innerErr := conn.Close()
		if innerErr != nil {
			log.Fatal(innerErr)
		}
	}(conn)

	require.NoError(t, err)
	client := locationV1beta1.NewLocationServiceClient(conn)

	t.Run("testUpdateLocation", func(t *testing.T) {
		testUpdateLocation(t, client)
	})
	t.Run("testQueryLocation", func(t *testing.T) {
		testQueryLocation(t, client)
	})
	t.Run("testTollgateCrossing", func(t *testing.T) {
		testTollgateCrossing(t, client)
	})
}

func dialer(detector tollgate.Detector) func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)

	server := grpc.NewServer()
	controller := New(redisClient.NewStoreClient(), redisClient.NewPubSubClient(), areaKey, detector)
	locationV1beta1.RegisterLocationServiceServer(server, controller)

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func testUpdateLocation(t *testing.T, client locationV1beta1.LocationServiceClient) {
	id := uuid.NewString()
	request := &locationV1beta1.UpdateLocationRequest{
		WorkerId: id,
		Lon:      12.000001966953278,
		Lat:      13.000001966953278,
	}
	response, err := client.UpdateLocation(context.Background(), request)
	require.NoError(t, err)
	require.Equal(t, request.WorkerId, response.WorkerId)

	location, err := client.QueryLocation(context.Background(), &locationV1beta1.QueryLocationRequest{
		WorkerId: id,
	})
	require.NoError(t, err)
	require.Equal(t, request.WorkerId, location.WorkerId)
	require.InDelta(t, request.Lon, location.Lon, 0.001)
	require.InDelta(t, request.Lat, location.Lat, 0.001)
}

func testQueryLocation(t *testing.T, client locationV1beta1.LocationServiceClient) {
	id := uuid.NewString()
	request := &locationV1beta1.QueryLocationRequest{
		WorkerId: id,
	}

	ctx := context.Background()

	_, err := client.QueryLocation(ctx, request)
	require.Error(t, err)
	require.Contains(t, err.Error(), fmt.Sprintf("location not found for WorkerId=%s", request.WorkerId))

	response, err := client.UpdateLocation(ctx, &locationV1beta1.UpdateLocationRequest{
		WorkerId: id,
		Lon:      12,
		Lat:      13,
	},
	)
	require.NoError(t, err)
	require.Equal(t, request.WorkerId, response.WorkerId)

	location, err := client.QueryLocation(ctx, request)
	require.NoError(t, err)
	require.Less(t, location.LastSeenTime.AsTime().UnixMilli(), time.Now().UnixMilli())
}

func testTollgateCrossing(t *testing.T, client locationV1beta1.LocationServiceClient) {
	ctx := context.Background()
	id := uuid.NewString()
	from := &locationV1beta1.UpdateLocationRequest{
		WorkerId: id,
		Lon:      -79.871248,
		Lat:      41.199493,
	}
	to := &locationV1beta1.UpdateLocationRequest{
		WorkerId: id,
		Lon:      -79.867927,
		Lat:      41.199329,
	}

	sync := make(chan string)
	var crossings <-chan tollgate.Crossing
	go func() {
		crossings = subscribe(crossingChannel("*"))
		sync <- "done"
	}()

	select {
	case <-sync:
		break
	case <-time.After(5 * time.Second):
		require.Fail(t, "sync timeout")
	}

	response1, err := client.UpdateLocation(ctx, from)
	require.NoError(t, err)
	require.Equal(t, from.WorkerId, response1.WorkerId)
	require.Nil(t, response1.TollgateCrossing)

	response2, err := client.UpdateLocation(ctx, to)
	require.NoError(t, err)
	require.Equal(t, from.WorkerId, response2.WorkerId)
	require.NotNil(t, response2.TollgateCrossing)

	c := <-crossings
	require.Equal(t, tollgateID, c.TollgateID)
	require.Equal(t, tollgate.Direction("SE"), c.Direction)
	require.Equal(t, id, c.WorkerID)
	require.InDelta(t, to.Lat, c.Movement.To.Lat, 0.003)
	require.InDelta(t, to.Lon, c.Movement.To.Lon, 0.003)
}
