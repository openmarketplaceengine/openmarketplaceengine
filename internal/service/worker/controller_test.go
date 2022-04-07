package worker

import (
	"context"
	"log"
	"net"
	"testing"

	workerV1beta1 "github.com/openmarketplaceengine/openmarketplaceengine/internal/omeapi/worker/v1beta1"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

func TestController(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(dialer()))
	defer func(conn *grpc.ClientConn) {
		innerErr := conn.Close()
		if innerErr != nil {
			log.Fatal(innerErr)
		}
	}(conn)

	require.NoError(t, err)
	client := workerV1beta1.NewWorkerServiceClient(conn)

	t.Run("testSetState", func(t *testing.T) {
		testSetState(t, client)
	})

	t.Run("testQueryByState", func(t *testing.T) {
		testQueryByState(t, client)
	})
}

func dialer() func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)

	server := grpc.NewServer()
	controller := newController()
	workerV1beta1.RegisterWorkerServiceServer(server, controller)

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func testSetState(t *testing.T, client workerV1beta1.WorkerServiceClient) {
	id := uuid.NewString()
	request := &workerV1beta1.SetStateRequest{
		WorkerId: id,
		State:    workerV1beta1.WorkerState_WORKER_STATE_ONLINE,
	}
	response, err := client.SetState(context.Background(), request)
	require.NoError(t, err)
	require.Equal(t, request.WorkerId, response.Worker.WorkerId)

	state, err := client.GetWorker(context.Background(), &workerV1beta1.GetWorkerRequest{
		WorkerId: id,
	})
	require.NoError(t, err)
	require.Equal(t, request.WorkerId, state.Worker.WorkerId)
	require.Equal(t, request.State, state.Worker.State, 0.001)
}

func testQueryByState(t *testing.T, client workerV1beta1.WorkerServiceClient) {
	id := uuid.NewString()
	request1 := &workerV1beta1.SetStateRequest{
		WorkerId: id,
		State:    workerV1beta1.WorkerState_WORKER_STATE_ONLINE,
	}
	request2 := &workerV1beta1.SetStateRequest{
		WorkerId: id,
		State:    workerV1beta1.WorkerState_WORKER_STATE_OFFLINE,
	}
	_, err := client.SetState(context.Background(), request1)
	require.NoError(t, err)

	_, err = client.SetState(context.Background(), request2)
	require.NoError(t, err)

	r1, err := client.QueryByState(context.Background(), &workerV1beta1.QueryByStateRequest{
		State: workerV1beta1.WorkerState_WORKER_STATE_ONLINE,
	})
	require.NoError(t, err)
	require.Len(t, r1.Workers, 1)

	r2, err := client.QueryByState(context.Background(), &workerV1beta1.QueryByStateRequest{
		State: workerV1beta1.WorkerState_WORKER_STATE_ON_JOB,
	})
	require.NoError(t, err)
	require.Len(t, r2.Workers, 0)
}
