package worker

import (
	"context"
	"fmt"

	"github.com/openmarketplaceengine/openmarketplaceengine/dom/worker"

	"github.com/openmarketplaceengine/openmarketplaceengine/internal/validate"

	"github.com/openmarketplaceengine/openmarketplaceengine/srv"
	"google.golang.org/grpc"

	workerV1beta1 "github.com/openmarketplaceengine/openmarketplaceengine/internal/api/worker/v1beta1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Controller struct {
	workerV1beta1.UnimplementedWorkerServiceServer
}

func newController() *Controller {
	return &Controller{}
}

func GrpcRegister() {
	srv.Grpc.Register(func(srv *grpc.Server) error {
		controller := newController()
		workerV1beta1.RegisterWorkerServiceServer(srv, controller)
		return nil
	})
}

func (c *Controller) GetWorker(ctx context.Context, request *workerV1beta1.GetWorkerRequest) (*workerV1beta1.GetWorkerResponse, error) {
	workerID := request.GetWorkerId()
	var v validate.Validator
	v.ValidateString("worker_id", workerID, validate.IsNotNull)

	errorInfo := v.ErrorInfo()
	if errorInfo != nil {
		st, err := status.New(codes.InvalidArgument, "bad request").
			WithDetails(errorInfo)
		if err != nil {
			panic(fmt.Errorf("enrich grpc status with details error: %w", err))
		}
		return nil, st.Err()
	}

	wrk, has, err := worker.QueryOne(ctx, request.WorkerId)

	if err != nil {
		st := status.Newf(codes.Internal, "get worker error: %s", err)
		st, err = st.WithDetails(request)
		if err != nil {
			panic(fmt.Errorf("enrich grpc status with details error: %w", err))
		}
		return nil, st.Err()
	}

	if !has {
		st := status.Newf(codes.NotFound, "worker not found")
		st, err := st.WithDetails(request)
		if err != nil {
			panic(fmt.Errorf("enrich grpc status with details error: %w", err))
		}
		return nil, st.Err()
	}

	w := transformWorker(wrk)
	return &workerV1beta1.GetWorkerResponse{
		Worker: w,
	}, nil
}

func (c *Controller) UpdateWorkerStatus(ctx context.Context, request *workerV1beta1.UpdateWorkerStatusRequest) (*workerV1beta1.UpdateWorkerStatusResponse, error) {
	workerID := request.GetWorkerId()
	workerStatus := request.GetStatus()
	var v validate.Validator
	v.ValidateString("worker_id", workerID, validate.IsNotNull)
	v.Validate("worker_status", workerStatus, ValidateStatus)

	errorInfo := v.ErrorInfo()
	if errorInfo != nil {
		st, err := status.New(codes.InvalidArgument, "bad request").
			WithDetails(errorInfo)
		if err != nil {
			panic(fmt.Errorf("enrich grpc status with details error: %w", err))
		}
		return nil, st.Err()
	}
	s := transformStatusFrom(workerStatus)
	err := worker.UpdateWorkerStatus(ctx, workerID, *s)

	if err != nil {
		st := status.Newf(codes.Internal, "update worker status error: %s", err)
		return nil, st.Err()
	}

	return &workerV1beta1.UpdateWorkerStatusResponse{
		WorkerId: workerID,
		Status:   workerStatus,
	}, nil
}

func (c *Controller) ListWorkers(ctx context.Context, request *workerV1beta1.ListWorkersRequest) (*workerV1beta1.ListWorkersResponse, error) {
	pageSize := request.GetPageSize()
	pageToken := request.GetPageToken()
	st := request.GetStatus()
	var v validate.Validator
	v.ValidateString("page_token", pageToken, validate.IsNotNull)
	v.Validate("page_size", pageSize, ValidatePageSize)

	var s *worker.Status
	if st != workerV1beta1.WorkerStatus_WORKER_STATUS_UNSPECIFIED {
		s = transformStatusFrom(st)
	}
	//todo make use of page_token as offset
	all, err := worker.QueryAll(ctx, s, int(pageSize), 0)

	if err != nil {
		st := status.Newf(codes.Internal, "query all error: %s", err)
		return nil, st.Err()
	}

	var workers = make([]*workerV1beta1.Worker, 0, int(pageSize))
	for _, v := range all {
		workers = append(workers, transformWorker(v))
	}

	return &workerV1beta1.ListWorkersResponse{
		Workers:       workers,
		NextPageToken: "todo",
	}, nil
}

var statusTo = map[worker.Status]workerV1beta1.WorkerStatus{
	worker.Offline:  workerV1beta1.WorkerStatus_WORKER_STATUS_OFFLINE,
	worker.Ready:    workerV1beta1.WorkerStatus_WORKER_STATUS_READY,
	worker.OnJob:    workerV1beta1.WorkerStatus_WORKER_STATUS_ON_JOB,
	worker.Paused:   workerV1beta1.WorkerStatus_WORKER_STATUS_PAUSED,
	worker.Disabled: workerV1beta1.WorkerStatus_WORKER_STATUS_DISABLED,
}

var statusFrom = map[workerV1beta1.WorkerStatus]worker.Status{
	workerV1beta1.WorkerStatus_WORKER_STATUS_OFFLINE:  worker.Offline,
	workerV1beta1.WorkerStatus_WORKER_STATUS_READY:    worker.Ready,
	workerV1beta1.WorkerStatus_WORKER_STATUS_ON_JOB:   worker.OnJob,
	workerV1beta1.WorkerStatus_WORKER_STATUS_PAUSED:   worker.Paused,
	workerV1beta1.WorkerStatus_WORKER_STATUS_DISABLED: worker.Disabled,
}

func transformWorker(wrk *worker.Worker) *workerV1beta1.Worker {
	s, ok := statusTo[wrk.Status]
	if !ok {
		s = workerV1beta1.WorkerStatus_WORKER_STATUS_UNSPECIFIED
	}
	return &workerV1beta1.Worker{
		WorkerId: wrk.ID,
		Status:   s,
	}
}

func transformStatusFrom(sta workerV1beta1.WorkerStatus) *worker.Status {
	s, ok := statusFrom[sta]
	if !ok {
		return nil
	}
	return &s
}

// ValidateStatus should not allow UNSPECIFIED status to pass through, as
// Controller API is explicit about status payload.
func ValidateStatus(value interface{}) error {
	v, ok := value.(workerV1beta1.WorkerStatus)
	if !ok || v == 0 {
		return fmt.Errorf("illegal status: %v", value)
	}
	return nil
}

func ValidatePageSize(value interface{}) error {
	v, ok := value.(int32)
	if ok && v > 0 {
		return nil
	}
	return fmt.Errorf("illegal page_size: %v", value)
}
