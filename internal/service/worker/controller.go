package worker

import (
	"context"
	"fmt"

	"github.com/openmarketplaceengine/openmarketplaceengine/internal/validate"

	"github.com/openmarketplaceengine/openmarketplaceengine/srv"
	"google.golang.org/grpc"

	workerV1beta1 "github.com/openmarketplaceengine/openmarketplaceengine/internal/omeapi/worker/v1beta1"
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

	worker, has, err := GetWorker(ctx, request.WorkerId)

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

	w := transformWorker(worker)
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
	err := SetWorkerStatus(ctx, workerID, s)

	if err != nil {
		st := status.Newf(codes.Internal, "update worker status error: %s", err)
		st, err = st.WithDetails(request)
		if err != nil {
			panic(fmt.Errorf("enrich grpc status with details error: %w", err))
		}
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
	v.Validate("status", st, ValidateStatus)

	var s Status
	if st != workerV1beta1.WorkerStatus_WORKER_STATUS_UNSPECIFIED {
		s = transformStatusFrom(st)
	}
	//todo make use of page_token as offset
	all, err := QueryAll(ctx, &s, int(pageSize), 0)

	if err != nil {
		st := status.Newf(codes.Internal, "query all error: %s", err)
		st, err = st.WithDetails(request)
		if err != nil {
			panic(fmt.Errorf("enrich grpc status with details error: %w", err))
		}
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

var statusTo = map[Status]workerV1beta1.WorkerStatus{
	Offline:  workerV1beta1.WorkerStatus_WORKER_STATUS_OFFLINE,
	Ready:    workerV1beta1.WorkerStatus_WORKER_STATUS_READY,
	OnJob:    workerV1beta1.WorkerStatus_WORKER_STATUS_ON_JOB,
	Paused:   workerV1beta1.WorkerStatus_WORKER_STATUS_PAUSED,
	Disabled: workerV1beta1.WorkerStatus_WORKER_STATUS_DISABLED,
}

var statusFrom = map[workerV1beta1.WorkerStatus]Status{
	workerV1beta1.WorkerStatus_WORKER_STATUS_OFFLINE:  Offline,
	workerV1beta1.WorkerStatus_WORKER_STATUS_READY:    Ready,
	workerV1beta1.WorkerStatus_WORKER_STATUS_ON_JOB:   OnJob,
	workerV1beta1.WorkerStatus_WORKER_STATUS_PAUSED:   Paused,
	workerV1beta1.WorkerStatus_WORKER_STATUS_DISABLED: Disabled,
}

func transformWorker(wrk *Worker) *workerV1beta1.Worker {
	s, ok := statusTo[wrk.Status]
	if !ok {
		s = workerV1beta1.WorkerStatus_WORKER_STATUS_UNSPECIFIED
	}
	return &workerV1beta1.Worker{
		WorkerId: wrk.ID,
		Status:   s,
	}
}

func transformStatusFrom(sta workerV1beta1.WorkerStatus) Status {
	s, ok := statusFrom[sta]
	if !ok {
		s = Offline
	}
	return s
}

func ValidateStatus(value interface{}) error {
	v, ok := value.(workerV1beta1.WorkerStatus)
	if !ok {
		return fmt.Errorf("not a status value: %v", value)
	}
	switch v {
	case workerV1beta1.WorkerStatus_WORKER_STATUS_OFFLINE:
	case workerV1beta1.WorkerStatus_WORKER_STATUS_READY:
	case workerV1beta1.WorkerStatus_WORKER_STATUS_ON_JOB:
	case workerV1beta1.WorkerStatus_WORKER_STATUS_PAUSED:
	case workerV1beta1.WorkerStatus_WORKER_STATUS_DISABLED:
		break
	default:
		return fmt.Errorf("unknown status value: %v", value)
	}
	return nil
}

func ValidatePageSize(value interface{}) error {
	v, ok := value.(int32)
	if ok && v > 0 {
		return nil
	}
	return fmt.Errorf("wrong page_size value: %v", value)
}
