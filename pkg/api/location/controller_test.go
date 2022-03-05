package location

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
	v1 "github.com/openmarketplaceengine/openmarketplaceengine/pkg/api/location/proto/v1"
	redisClient "github.com/openmarketplaceengine/openmarketplaceengine/redis/client"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"net"
	"testing"
)

const port = 10123
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
		controller := New(redisClient.NewStoreClient(), redisClient.NewPubSubClient(), redisClient.NewPubSubClient(), areaKey)
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
	t.Run("testQueryLocationStreaming", func(t *testing.T) {
		testQueryLocationStreaming(t, client)
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

func testQueryLocationStreaming(t *testing.T, client v1.LocationServiceClient) {
	id := uuid.NewString()
	query := &v1.QueryLocationStreamingRequest{
		WorkerId: id,
	}

	ctx := context.Background()

	streaming, err := client.QueryLocationStreaming(ctx, query)
	require.NoError(t, err)
	done := make(chan string)

	go func() {
		for {
			response, err := streaming.Recv()
			if err == io.EOF {
				fmt.Printf("===> Recv EOF %s\n", err)
				return
			}
			if err != nil {
				fmt.Printf("===> Recv err %s\n", err)
				return
			}
			fmt.Printf("===> Recv %s\n", response)
			done <- response.WorkerId
		}
	}()

	for i := 0; i < 3; i++ {
		_, err = client.UpdateLocation(ctx, &v1.UpdateLocationRequest{
			WorkerId:  query.WorkerId,
			Longitude: 13,
			Latitude:  14,
		},
		)
		require.NoError(t, err)
	}

	require.Equal(t, query.WorkerId, <-done)
}
