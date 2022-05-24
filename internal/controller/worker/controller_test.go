package worker

import (
	"context"
	"log"
	"net"
	"testing"

	"github.com/openmarketplaceengine/openmarketplaceengine/dao"
	"github.com/openmarketplaceengine/openmarketplaceengine/dom/worker"

	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
	"github.com/openmarketplaceengine/openmarketplaceengine/dom"

	workerV1beta1 "github.com/openmarketplaceengine/openmarketplaceengine/internal/api/worker/v1beta1"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

func TestController(t *testing.T) {
	dom.WillTest(t, "test", true)
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

	t.Run("testUpdateWorkerState", func(t *testing.T) {
		testUpdateWorkerState(t, client)
	})
	t.Run("testUpdateWorkerStateBadRequest", func(t *testing.T) {
		testUpdateWorkerStateBadRequest(t, client)
	})
	t.Run("testListWorkersByState", func(t *testing.T) {
		testListWorkersByState(t, client)
	})
}

func dialer() func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)

	server := grpc.NewServer()
	workerV1beta1.RegisterWorkerServiceServer(server, &controller{})

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func testUpdateWorkerState(t *testing.T, client workerV1beta1.WorkerServiceClient) {
	ctx := cfg.Context()

	w := newWorker()
	require.NoError(t, w.Insert(ctx))

	request := &workerV1beta1.UpdateWorkerStatusRequest{
		WorkerId: w.ID,
		Status:   workerV1beta1.WorkerStatus_WORKER_STATUS_ON_JOB,
	}
	response, err := client.UpdateWorkerStatus(ctx, request)
	require.NoError(t, err)
	require.Equal(t, request.WorkerId, response.WorkerId)

	state, err := client.GetWorker(ctx, &workerV1beta1.GetWorkerRequest{
		WorkerId: w.ID,
	})
	require.NoError(t, err)
	require.Equal(t, request.WorkerId, state.Worker.WorkerId)
	require.Equal(t, request.Status, state.Worker.Status)
}

func testUpdateWorkerStateBadRequest(t *testing.T, client workerV1beta1.WorkerServiceClient) {
	ctx := cfg.Context()

	w := newWorker()
	require.NoError(t, w.Insert(ctx))

	_, err := client.UpdateWorkerStatus(ctx, &workerV1beta1.UpdateWorkerStatusRequest{
		WorkerId: w.ID,
	},
	)
	require.Error(t, err)
	require.Contains(t, err.Error(), "bad request")

	_, err = client.UpdateWorkerStatus(ctx, &workerV1beta1.UpdateWorkerStatusRequest{
		WorkerId: w.ID,
		Status:   workerV1beta1.WorkerStatus_WORKER_STATUS_OFFLINE,
	},
	)
	require.NoError(t, err)
}

func testListWorkersByState(t *testing.T, client workerV1beta1.WorkerServiceClient) {
	ctx := cfg.Context()

	w := newWorker()
	require.NoError(t, w.Insert(ctx))

	request1 := &workerV1beta1.UpdateWorkerStatusRequest{
		WorkerId: w.ID,
		Status:   workerV1beta1.WorkerStatus_WORKER_STATUS_PAUSED,
	}
	request2 := &workerV1beta1.UpdateWorkerStatusRequest{
		WorkerId: w.ID,
		Status:   workerV1beta1.WorkerStatus_WORKER_STATUS_DISABLED,
	}
	_, err := client.UpdateWorkerStatus(context.Background(), request1)
	require.NoError(t, err)

	_, err = client.UpdateWorkerStatus(context.Background(), request2)
	require.NoError(t, err)

	r1, err := client.ListWorkers(context.Background(), &workerV1beta1.ListWorkersRequest{
		Status:    workerV1beta1.WorkerStatus_WORKER_STATUS_PAUSED,
		PageSize:  10,
		PageToken: "",
	})
	require.NoError(t, err)
	require.Len(t, r1.Workers, 0)

	r2, err := client.ListWorkers(context.Background(), &workerV1beta1.ListWorkersRequest{
		Status:    workerV1beta1.WorkerStatus_WORKER_STATUS_DISABLED,
		PageSize:  10,
		PageToken: "",
	})
	require.NoError(t, err)
	require.Len(t, r2.Workers, 1)
}

func newWorker() *worker.Worker {
	stamp := dom.Time{}
	stamp.Now()
	return &worker.Worker{
		ID:        dao.MockUUID(),
		Status:    worker.Ready,
		Rating:    0,
		Jobs:      0,
		FirstName: "",
		LastName:  "",
		Vehicle:   "",
		Created:   stamp,
		Updated:   stamp,
	}
}
