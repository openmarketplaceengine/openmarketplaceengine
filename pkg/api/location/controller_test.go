package location

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
	v1 "github.com/openmarketplaceengine/openmarketplaceengine/pkg/api/location/proto/v1"
	redisClient "github.com/openmarketplaceengine/openmarketplaceengine/redis/client"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const port = 8090
const areaKey = "san_fran"

func TestController(t *testing.T) {
	err := cfg.Load()
	require.NoError(t, err)

	address := fmt.Sprintf("localhost:%d", port)

	go func() {
		lis, innerErr := net.Listen("tcp", address)
		if innerErr != nil {
			panic(innerErr)
		}
		grpcServer := grpc.NewServer()
		controller := New(redisClient.NewStoreClient(), redisClient.NewPubSubClient(), areaKey)
		v1.RegisterLocationServiceServer(grpcServer, controller)
		innerErr = grpcServer.Serve(lis)
		if innerErr != nil {
			panic(innerErr)
		}
	}()

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	client := v1.NewLocationServiceClient(conn)

	t.Run("testUpdateLocation", func(t *testing.T) {
		testUpdateLocation(t, client)
	})
	t.Run("testQueryLocation", func(t *testing.T) {
		testQueryLocation(t, client)
	})
}

func testUpdateLocation(t *testing.T, client v1.LocationServiceClient) {
	id := uuid.NewString()
	request := &v1.UpdateLocationRequest{
		WorkerId:  id,
		Longitude: 12.000001966953278,
		Latitude:  13.000001966953278,
	}
	response, err := client.UpdateLocation(context.Background(), request)
	require.NoError(t, err)
	require.Equal(t, request.WorkerId, response.WorkerId)

	location, err := client.QueryLocation(context.Background(), &v1.QueryLocationRequest{
		WorkerId: id,
	})
	require.NoError(t, err)
	require.Equal(t, request.WorkerId, location.WorkerId)
	require.InDelta(t, request.Longitude, location.Longitude, 0.001)
	require.InDelta(t, request.Latitude, location.Latitude, 0.001)
}

func testQueryLocation(t *testing.T, client v1.LocationServiceClient) {
	id := uuid.NewString()
	request := &v1.QueryLocationRequest{
		WorkerId: id,
	}

	ctx := context.Background()

	_, err := client.QueryLocation(ctx, request)
	require.Error(t, err)
	require.Contains(t, err.Error(), fmt.Sprintf("location not found for WorkerId=%s", request.WorkerId))

	response, err := client.UpdateLocation(ctx, &v1.UpdateLocationRequest{
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
