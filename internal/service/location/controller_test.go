package location

import (
	"context"
	"log"
	"net"
	"testing"
	"time"

	"github.com/openmarketplaceengine/openmarketplaceengine/internal/service/location/storage"

	"github.com/openmarketplaceengine/openmarketplaceengine/dao"

	"github.com/openmarketplaceengine/openmarketplaceengine/internal/service/tollgate"

	"github.com/openmarketplaceengine/openmarketplaceengine/dom"
	"github.com/openmarketplaceengine/openmarketplaceengine/internal/service/tollgate/detector"

	"github.com/google/uuid"
	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
	locationV1beta1 "github.com/openmarketplaceengine/openmarketplaceengine/internal/omeapi/location/v1beta1"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const tollgateID = "tollgate-123"

func TestController(t *testing.T) {
	err := cfg.Load()
	require.NoError(t, err)

	dom.WillTest(t, "test", false)
	if !dao.Reds.State.Running() {
		require.NoError(t, dao.Reds.Boot())
	}

	ctx := context.Background()

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

	require.NoError(t, err)

	conn, err := grpc.DialContext(ctx, "", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(dialer(t)))
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

func dialer(t *testing.T) func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)

	server := grpc.NewServer()
	controller, err := newController(dao.Reds.StoreClient, dao.Reds.PubSubClient)
	require.NoError(t, err)
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
	require.Contains(t, err.Error(), "NotFound")

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
	s := storage.New(dao.Reds.StoreClient)
	ctx := context.Background()
	id := uuid.NewString()
	err := s.ForgetLocation(ctx, areaKey, id)
	require.NoError(t, err)
	from := &locationV1beta1.UpdateLocationRequest{
		WorkerId: id,
		Lon:      -74.195995,
		Lat:      40.636916,
	}
	to := &locationV1beta1.UpdateLocationRequest{
		WorkerId: id,
		Lon:      -74.198356,
		Lat:      40.634408,
	}

	sync := make(chan string)
	var crossings <-chan detector.Crossing
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
	require.Equal(t, detector.Direction("SW"), c.Direction)
	require.Equal(t, id, c.WorkerID)
	require.InDelta(t, to.Lat, c.Movement.To.Lat, 0.003)
	require.InDelta(t, to.Lon, c.Movement.To.Lon, 0.003)
}
