package worker

import (
	"context"

	"github.com/openmarketplaceengine/openmarketplaceengine/srv"
	"google.golang.org/grpc"

	workerV1beta1 "github.com/openmarketplaceengine/openmarketplaceengine/internal/omeapi/worker/v1beta1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Controller struct {
	workerV1beta1.UnimplementedWorkerServiceServer
	states map[string]workerV1beta1.WorkerState
}

func newController() *Controller {
	return &Controller{
		states: make(map[string]workerV1beta1.WorkerState),
	}
}

func GrpcRegister() {
	srv.Grpc.Register(func(srv *grpc.Server) error {
		controller := newController()
		workerV1beta1.RegisterWorkerServiceServer(srv, controller)
		return nil
	})
}

func (c *Controller) GetWorker(ctx context.Context, request *workerV1beta1.GetWorkerRequest) (*workerV1beta1.GetWorkerResponse, error) {
	v, ok := c.states[request.WorkerId]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "WorkerId=%s", request.WorkerId)
	}
	return &workerV1beta1.GetWorkerResponse{
		Worker: &workerV1beta1.Worker{
			WorkerId: request.WorkerId,
			State:    v,
		},
	}, nil
}
func (c *Controller) SetState(ctx context.Context, request *workerV1beta1.SetStateRequest) (*workerV1beta1.SetStateResponse, error) {
	c.states[request.WorkerId] = request.GetState()
	return &workerV1beta1.SetStateResponse{
		Worker: &workerV1beta1.Worker{
			WorkerId: request.GetWorkerId(),
			State:    request.State,
		},
	}, nil
}

func (c *Controller) QueryByState(ctx context.Context, request *workerV1beta1.QueryByStateRequest) (*workerV1beta1.QueryByStateResponse, error) {
	var workers []*workerV1beta1.Worker
	for k, v := range c.states {
		if v == request.GetState() {
			workers = append(workers, &workerV1beta1.Worker{
				WorkerId: k,
				State:    v,
			})
		}
	}

	return &workerV1beta1.QueryByStateResponse{
		Workers: workers,
	}, nil
}
