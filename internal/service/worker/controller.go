package worker

import (
	"context"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"strings"

	workerV1beta1 "github.com/openmarketplaceengine/openmarketplaceengine/internal/omeapi/worker/v1beta1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Controller struct {
	workerV1beta1.UnimplementedWorkerServiceServer
	states map[string]workerV1beta1.WorkerState
}

func New() *Controller {
	return &Controller{
		states: make(map[string]workerV1beta1.WorkerState),
	}
}

func (c *Controller) GetWorker(ctx context.Context, request *workerV1beta1.GetWorkerRequest) (*workerV1beta1.GetWorkerResponse, error) {
	if strings.TrimSpace(request.GetWorkerId()) == "" {
		st, err := status.New(codes.NotFound, "Worker not found").WithDetails(
			request, &errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "worker_id",
						Description: "cannot be empty",
					},
				},
			},
		)
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
		return nil, st.Err()
	}
	v, ok := c.states[request.WorkerId]
	if !ok {
		st, err := status.New(codes.NotFound, "Worker not found").WithDetails(request)
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
		return nil, st.Err()
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
