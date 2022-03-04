package v1

import (
	"context"
	"fmt"
	"io"
	"net"
	"testing"

	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
	redisClient "github.com/openmarketplaceengine/openmarketplaceengine/redis/client"

	"github.com/google/uuid"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const port = 10123
const areaKey = "san_fran"

func TestController(t *testing.T) {
	err := cfg.Load()
	require.NoError(t, err)

	storeClient := redisClient.NewStoreClient()
	require.NotNil(t, storeClient)

	go func() {
		lis, innerErr := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
		if innerErr != nil {
			panic(innerErr)
		}
		grpcServer := grpc.NewServer()
		RegisterLocationServiceServer(grpcServer, New(storeClient, areaKey))
		innerErr = grpcServer.Serve(lis)
		if innerErr != nil {
			panic(innerErr)
		}
	}()

	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", port), grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials()))

	require.NoError(t, err)

	defer func(conn *grpc.ClientConn) {
		innerErr := conn.Close()
		if innerErr != nil {
			require.NoError(t, innerErr)
		}
	}(conn)

	client := NewLocationServiceClient(conn)

	t.Run("testUpdateLocation", func(t *testing.T) {
		testUpdateLocation(t, client)
	})
	t.Run("testQueryLocationStreaming", func(t *testing.T) {
		testQueryLocationStreaming(t, client)
	})
}

func testUpdateLocation(t *testing.T, client LocationServiceClient) {
	id := uuid.NewString()
	request := &UpdateLocationRequest{
		WorkerId:  id,
		Longitude: 12.000001966953278,
		Latitude:  13.000001966953278,
	}
	response, err := client.UpdateLocation(context.Background(), request)
	require.NoError(t, err)
	require.Equal(t, request.WorkerId, response.WorkerId)

	location, err := client.QueryLocation(context.Background(), &QueryLocationRequest{
		WorkerId: id,
	})
	require.NoError(t, err)
	require.Equal(t, request.WorkerId, location.WorkerId)
	require.InDelta(t, request.Longitude, location.Longitude, 0.001)
	require.InDelta(t, request.Latitude, location.Latitude, 0.001)
}

func testQueryLocationStreaming(t *testing.T, client LocationServiceClient) {
	id := uuid.NewString()
	query := &QueryLocationStreamingRequest{
		WorkerId: id,
	}

	streaming, err := client.QueryLocationStreaming(context.Background(), query)
	require.NoError(t, err)
	sync := make(chan string, 1)
	go func() {
		for i := 0; i < 3; i++ {
			response, innerErr := streaming.Recv()
			if innerErr == io.EOF {
				break
			}
			require.NoError(t, innerErr)
			require.NotNil(t, response)
			if i <= 2 {
				sync <- response.WorkerId
			}
		}
	}()

	upd := &UpdateLocationRequest{
		WorkerId:  id,
		Longitude: 3,
		Latitude:  3,
	}
	for i := 0; i < 3; i++ {
		_, err = client.UpdateLocation(context.Background(), upd)
		require.NoError(t, err)
	}
	require.Equal(t, query.WorkerId, <-sync)
}
