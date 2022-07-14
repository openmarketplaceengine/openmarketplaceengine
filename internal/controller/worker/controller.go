package worker

import (
	"context"
	"fmt"

	"github.com/openmarketplaceengine/openmarketplaceengine/dom/worker"
	"github.com/openmarketplaceengine/openmarketplaceengine/internal/idl/api/worker/v1beta1"
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

func (c *controller) GetWorker(ctx context.Context, req *v1beta1.GetWorkerRequest) (*v1beta1.GetWorkerResponse, error) {
	err := req.ValidateAll()

	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	wrk, has, err := worker.QueryOne(ctx, req.WorkerId)

	if err != nil {
		st := status.Newf(codes.Internal, "get worker error: %s", err)
		st, err = st.WithDetails(req)
		if err != nil {
			panic(fmt.Errorf("enrich grpc status with details error: %w", err))
		}
		return nil, st.Err()
	}

	if !has {
		st := status.Newf(codes.NotFound, "worker not found")
		st, err := st.WithDetails(req)
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

func (c *controller) UpdateWorkerStatus(ctx context.Context, req *v1beta1.UpdateWorkerStatusRequest) (*v1beta1.UpdateWorkerStatusResponse, error) {
	err := req.ValidateAll()

	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	workerID := req.GetWorkerId()
	workerStatus := req.GetStatus()

	s := transformStatusFrom(workerStatus)
	err = worker.UpdateWorkerStatus(ctx, workerID, *s)

	if err != nil {
		st := status.Newf(codes.Internal, "update worker status error: %s", err)
		return nil, st.Err()
	}

	return &v1beta1.UpdateWorkerStatusResponse{
		WorkerId: workerID,
		Status:   workerStatus,
	}, nil
}

func (c *controller) ListWorkers(ctx context.Context, req *v1beta1.ListWorkersRequest) (*v1beta1.ListWorkersResponse, error) {
	err := req.ValidateAll()

	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	pageSize := req.GetPageSize()
	st := req.GetStatus()

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
