package worker

import (
	"context"
	"fmt"

	"github.com/openmarketplaceengine/openmarketplaceengine/dom/worker"
	v1beta1 "github.com/openmarketplaceengine/openmarketplaceengine/internal/api/worker/v1beta1"
	"github.com/openmarketplaceengine/openmarketplaceengine/internal/validate"
	"github.com/openmarketplaceengine/openmarketplaceengine/srv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type controller struct {
	v1beta1.UnimplementedWorkerServiceServer
}

func init() {
	srv.Grpc.Register(func(s *grpc.Server) error {
		srv.Grpc.Infof("registering: %s", v1beta1.WorkerService_ServiceDesc.ServiceName)
		v1beta1.RegisterWorkerServiceServer(s, &controller{})
		return nil
	})
}

func (c *controller) GetWorker(ctx context.Context, request *v1beta1.GetWorkerRequest) (*v1beta1.GetWorkerResponse, error) {
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
	return &v1beta1.GetWorkerResponse{
		Worker: w,
	}, nil
}

func (c *controller) UpdateWorkerStatus(ctx context.Context, request *v1beta1.UpdateWorkerStatusRequest) (*v1beta1.UpdateWorkerStatusResponse, error) {
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

	return &v1beta1.UpdateWorkerStatusResponse{
		WorkerId: workerID,
		Status:   workerStatus,
	}, nil
}

func (c *controller) ListWorkers(ctx context.Context, request *v1beta1.ListWorkersRequest) (*v1beta1.ListWorkersResponse, error) {
	pageSize := request.GetPageSize()
	pageToken := request.GetPageToken()
	st := request.GetStatus()
	var v validate.Validator
	v.ValidateString("page_token", pageToken, validate.IsNotNull)
	v.Validate("page_size", pageSize, ValidatePageSize)

	var s *worker.Status
	if st != v1beta1.WorkerStatus_WORKER_STATUS_UNSPECIFIED {
		s = transformStatusFrom(st)
	}
	//todo make use of page_token as offset
	all, err := worker.QueryAll(ctx, s, int(pageSize), 0)

	if err != nil {
		st := status.Newf(codes.Internal, "query all error: %s", err)
		return nil, st.Err()
	}

	var workers = make([]*v1beta1.Worker, 0, int(pageSize))
	for _, v := range all {
		workers = append(workers, transformWorker(v))
	}

	return &v1beta1.ListWorkersResponse{
		Workers:       workers,
		NextPageToken: "todo",
	}, nil
}

var statusTo = map[worker.Status]v1beta1.WorkerStatus{
	worker.Offline:  v1beta1.WorkerStatus_WORKER_STATUS_OFFLINE,
	worker.Ready:    v1beta1.WorkerStatus_WORKER_STATUS_READY,
	worker.OnJob:    v1beta1.WorkerStatus_WORKER_STATUS_ON_JOB,
	worker.Paused:   v1beta1.WorkerStatus_WORKER_STATUS_PAUSED,
	worker.Disabled: v1beta1.WorkerStatus_WORKER_STATUS_DISABLED,
}

var statusFrom = map[v1beta1.WorkerStatus]worker.Status{
	v1beta1.WorkerStatus_WORKER_STATUS_OFFLINE:  worker.Offline,
	v1beta1.WorkerStatus_WORKER_STATUS_READY:    worker.Ready,
	v1beta1.WorkerStatus_WORKER_STATUS_ON_JOB:   worker.OnJob,
	v1beta1.WorkerStatus_WORKER_STATUS_PAUSED:   worker.Paused,
	v1beta1.WorkerStatus_WORKER_STATUS_DISABLED: worker.Disabled,
}

func transformWorker(wrk *worker.Worker) *v1beta1.Worker {
	s, ok := statusTo[wrk.Status]
	if !ok {
		s = v1beta1.WorkerStatus_WORKER_STATUS_UNSPECIFIED
	}
	return &v1beta1.Worker{
		WorkerId: wrk.ID,
		Status:   s,
	}
}

func transformStatusFrom(sta v1beta1.WorkerStatus) *worker.Status {
	s, ok := statusFrom[sta]
	if !ok {
		return nil
	}
	return &s
}

// ValidateStatus should not allow UNSPECIFIED status to pass through, as
// controller API is explicit about status payload.
func ValidateStatus(value interface{}) error {
	v, ok := value.(v1beta1.WorkerStatus)
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
