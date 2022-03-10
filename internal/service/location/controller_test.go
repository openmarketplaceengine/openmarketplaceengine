package location

import (
	"context"
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	locationV1beta1 "github.com/openmarketplaceengine/openmarketplaceengine/internal/omeapi/location/v1beta1"

	"github.com/google/uuid"
	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
	redisClient "github.com/openmarketplaceengine/openmarketplaceengine/redis/client"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const areaKey = "san_fran"

func TestController(t *testing.T) {
	err := cfg.Load()
	require.NoError(t, err)

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(dialer()))
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
}

func dialer() func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)

	server := grpc.NewServer()
	controller := New(redisClient.NewStoreClient(), redisClient.NewPubSubClient(), areaKey)
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
		WorkerId:  id,
		Longitude: 12.000001966953278,
		Latitude:  13.000001966953278,
	}
	response, err := client.UpdateLocation(context.Background(), request)
	require.NoError(t, err)
	require.Equal(t, request.WorkerId, response.WorkerId)

	location, err := client.QueryLocation(context.Background(), &locationV1beta1.QueryLocationRequest{
		WorkerId: id,
	})
	require.NoError(t, err)
	require.Equal(t, request.WorkerId, location.WorkerId)
	require.InDelta(t, request.Longitude, location.Longitude, 0.001)
	require.InDelta(t, request.Latitude, location.Latitude, 0.001)
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
		WorkerId:  id,
		Longitude: 12,
		Latitude:  13,
	},
	)
	require.NoError(t, err)
	require.Equal(t, request.WorkerId, response.WorkerId)

	location, err := client.QueryLocation(ctx, request)
	require.NoError(t, err)
	require.Less(t, location.LastSeenTime.AsTime().UnixMilli(), time.Now().UnixMilli())
}
