package location

import (
	"context"
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
	"github.com/openmarketplaceengine/openmarketplaceengine/dao"
	"github.com/openmarketplaceengine/openmarketplaceengine/dom"
	"github.com/openmarketplaceengine/openmarketplaceengine/dom/tollgate"
	locationV1beta1 "github.com/openmarketplaceengine/openmarketplaceengine/internal/api/location/v1beta1"
	typeV1beta1 "github.com/openmarketplaceengine/openmarketplaceengine/internal/api/type/v1beta1"
	"github.com/openmarketplaceengine/openmarketplaceengine/pkg/detector"
	"github.com/openmarketplaceengine/openmarketplaceengine/svc/location"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/timestamppb"
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
	t.Run("testUpdateLocationBadRequest", func(t *testing.T) {
		testUpdateLocationBadRequest(t, client)
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
	controller, err := newController()
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
		Value: &locationV1beta1.LocationUpdate{
			WorkerId: id,
			Location: &typeV1beta1.Location{
				Longitude: 12.000001966953278,
				Latitude:  13.000001966953278,
			},
			UpdateTime: timestamppb.Now(),
		},
		AreaKey:      "a",
		ValidateOnly: false,
	}
	response, err := client.UpdateLocation(context.Background(), request)
	require.NoError(t, err)
	require.Equal(t, request.GetValue().WorkerId, response.WorkerId)

	l, err := client.GetLocation(context.Background(), &locationV1beta1.GetLocationRequest{
		WorkerId: id,
		AreaKey:  "a",
	})
	require.NoError(t, err)
	require.Equal(t, request.GetValue().GetWorkerId(), l.WorkerId)
	require.Equal(t, request.GetAreaKey(), "a")
}

func testUpdateLocationBadRequest(t *testing.T, client locationV1beta1.LocationServiceClient) {
	id := uuid.NewString()

	_, err := client.UpdateLocation(context.Background(), &locationV1beta1.UpdateLocationRequest{
		Value: &locationV1beta1.LocationUpdate{
			WorkerId: id,
			Location: &typeV1beta1.Location{
				Longitude: 1200,
				Latitude:  1300,
			},
			UpdateTime: timestamppb.Now(),
		},
		AreaKey:      "a",
		ValidateOnly: false,
	},
	)
	require.Error(t, err)
	require.EqualError(t, err, "rpc error: code = InvalidArgument desc = bad request")

	_, err = client.UpdateLocation(context.Background(), &locationV1beta1.UpdateLocationRequest{
		Value: &locationV1beta1.LocationUpdate{
			WorkerId: id,
			Location: &typeV1beta1.Location{
				Longitude: 12,
				Latitude:  13,
			},
			UpdateTime: timestamppb.Now(),
		},
		AreaKey:      "",
		ValidateOnly: false,
	},
	)
	require.Error(t, err)
	require.EqualError(t, err, "rpc error: code = InvalidArgument desc = bad request")

	_, err = client.GetLocation(context.Background(), &locationV1beta1.GetLocationRequest{
		WorkerId: id,
	})
	require.Error(t, err)
	require.EqualError(t, err, "rpc error: code = InvalidArgument desc = bad request")

	_, err = client.GetLocation(context.Background(), &locationV1beta1.GetLocationRequest{
		WorkerId: id,
		AreaKey:  "a",
	})
	require.Error(t, err)
	require.EqualError(t, err, "rpc error: code = NotFound desc = location not found")
}

func testQueryLocation(t *testing.T, client locationV1beta1.LocationServiceClient) {
	id := uuid.NewString()
	request := &locationV1beta1.GetLocationRequest{
		WorkerId: id,
		AreaKey:  "a",
	}

	ctx := context.Background()

	_, err := client.GetLocation(ctx, request)
	require.Error(t, err)
	require.Contains(t, err.Error(), "NotFound")

	response, err := client.UpdateLocation(ctx, &locationV1beta1.UpdateLocationRequest{
		Value: &locationV1beta1.LocationUpdate{
			WorkerId: id,
			Location: &typeV1beta1.Location{
				Longitude: 12.000001966953278,
				Latitude:  13.000001966953278,
			},
			UpdateTime: timestamppb.Now(),
		},
		AreaKey:      "a",
		ValidateOnly: false,
	},
	)
	require.NoError(t, err)
	require.Equal(t, request.WorkerId, response.WorkerId)

	l, err := client.GetLocation(ctx, request)
	require.NoError(t, err)
	require.Less(t, l.LastSeenTime.AsTime().UnixNano(), time.Now().UnixNano())
}

func testTollgateCrossing(t *testing.T, client locationV1beta1.LocationServiceClient) {
	s := location.NewStorage(dao.Reds.StoreClient)
	ctx := context.Background()
	id := uuid.NewString()
	err := s.ForgetLocation(ctx, "a", id)
	require.NoError(t, err)
	from := &locationV1beta1.UpdateLocationRequest{

		Value: &locationV1beta1.LocationUpdate{
			WorkerId: id,
			Location: &typeV1beta1.Location{
				Longitude: -74.195995,
				Latitude:  40.636916,
			},
			UpdateTime: timestamppb.Now(),
		},
		AreaKey:      "a",
		ValidateOnly: false,
	}
	to := &locationV1beta1.UpdateLocationRequest{

		Value: &locationV1beta1.LocationUpdate{
			WorkerId: id,
			Location: &typeV1beta1.Location{
				Longitude: -74.198356,
				Latitude:  40.634408,
			},
			UpdateTime: timestamppb.Now(),
		},
		AreaKey:      "a",
		ValidateOnly: false,
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
	require.Equal(t, from.GetValue().WorkerId, response1.WorkerId)
	require.Nil(t, response1.Crossing)

	response2, err := client.UpdateLocation(ctx, to)
	require.NoError(t, err)
	require.Equal(t, from.GetValue().WorkerId, response2.WorkerId)
	require.NotNil(t, response2.Crossing)

	c := <-crossings
	require.Equal(t, tollgateID, c.TollgateID)
	require.Equal(t, detector.Direction("SW"), c.Direction)
	require.Equal(t, id, c.WorkerID)
	require.InDelta(t, to.GetValue().GetLocation().GetLatitude(), c.Movement.To.Latitude, 0.003)
	require.InDelta(t, to.GetValue().GetLocation().GetLongitude(), c.Movement.To.Longitude, 0.003)
}

func crossingChannel(tollgateID string) string {
	return fmt.Sprintf("channel:crossing:%s", tollgateID)
}
