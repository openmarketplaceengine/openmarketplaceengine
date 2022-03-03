package location

import (
	"context"
	"fmt"
	"github.com/openmarketplaceengine/openmarketplaceengine/core/model/location/protos"
	"io"
	"net"
	"testing"

	"github.com/google/uuid"
	v1 "github.com/openmarketplaceengine/openmarketplaceengine/core/model/location/protos/v1"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const port = 10123

func TestLocationGRPC(t *testing.T) {
	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
		if err != nil {
			panic(err)
		}
		grpcServer := grpc.NewServer()
		v1.RegisterLocationServiceServer(grpcServer, &protos.Server{})
		err = grpcServer.Serve(lis)
		if err != nil {
			panic(err)
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
		Longitude: 3,
		Latitude:  3,
	}
	response, err := client.UpdateLocation(context.Background(), request)
	require.NoError(t, err)
	require.Equal(t, request.WorkerId, response.WorkerId)
}

func testQueryLocationStreaming(t *testing.T, client v1.LocationServiceClient) {
	id := uuid.NewString()
	query := &v1.QueryLocationRequest{
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

	upd := &v1.UpdateLocationRequest{
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
