// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package eta

import (
	"context"
	"fmt"
	"github.com/openmarketplaceengine/openmarketplaceengine/dom/job"
	rpc "github.com/openmarketplaceengine/openmarketplaceengine/internal/idl/api/eta/v1beta1"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/types/known/durationpb"
	"sync/atomic"

	"github.com/openmarketplaceengine/openmarketplaceengine/dao"
	"github.com/openmarketplaceengine/openmarketplaceengine/pkg/detector"
	"github.com/openmarketplaceengine/openmarketplaceengine/srv"
	svcJob "github.com/openmarketplaceengine/openmarketplaceengine/svc/job"
	"github.com/openmarketplaceengine/openmarketplaceengine/svc/location"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type controller struct {
	rpc.UnimplementedEstimatedJobServiceServer
	jobService *svcJob.Service
	nWorkers   int
}

func init() {
	srv.Grpc.Register(func(s *grpc.Server) error {
		srv.Grpc.Infof("registering: %s", rpc.EstimatedJobService_ServiceDesc.ServiceName)

		storeClient := dao.Reds.StoreClient
		noOp := detector.NewDetectorNoOp()
		storage := location.NewStorage(storeClient)
		service := svcJob.NewService(location.NewTracker(storage, noOp))
		rpc.RegisterEstimatedJobServiceServer(s, &controller{jobService: service, nWorkers: 3})
		return nil
	})
}

//-----------------------------------------------------------------------------

func (c *controller) GetEstimatedJobs(ctx context.Context, req *rpc.GetEstimatedJobsRequest) (*rpc.GetEstimatedJobsResponse, error) {

	if c.nWorkers < 1 {
		return nil, status.Errorf(codes.InvalidArgument, "at least one worker required")
	}

	err := req.ValidateAll()

	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	workerLocation, aJobs, err := c.jobService.QueryByPickupDistance(
		ctx, req.GetAreaKey(), req.GetWorkerId(), req.GetRadiusMeters(), req.GetLimit())

	if err != nil {
		return nil, fmt.Errorf("QueryByPickupDistance error: %w", err)
	}

	g, ctx := errgroup.WithContext(ctx)
	chunk := make(chan []*job.AvailableJob)

	g.Go(func() error {
		defer close(chunk)

		chunks := toChunks(aJobs, 10)

		for _, ch := range chunks {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case chunk <- ch:
			}
		}

		return nil
	})

	estimatedJobs := make(chan []*svcJob.EstimatedJob)

	workers := int32(c.nWorkers)
	for i := 0; i < c.nWorkers; i++ {
		g.Go(func() error {
			defer func() {
				// Last one closes
				if atomic.AddInt32(&workers, -1) == 0 {
					close(estimatedJobs)
				}
			}()

			for ch := range chunk {

				estimatedChunk, err := svcJob.EstimateJobs(ctx, workerLocation, ch)
				if err != nil {
					return fmt.Errorf("EstimateJobs error: %w", err)
				} else {
					select {
					case <-ctx.Done():
						return ctx.Err()
					case estimatedJobs <- estimatedChunk:
					}
				}
			}
			return nil
		})
	}

	eJobs := make([]*rpc.EstimatedJob, 0)
	g.Go(func() error {
		for ej := range estimatedJobs {
			jobs := transform(ej)
			eJobs = append(eJobs, jobs...)
		}
		return nil
	})

	return &rpc.GetEstimatedJobsResponse{Jobs: eJobs}, g.Wait()
}

func transform(eJobs []*svcJob.EstimatedJob) []*rpc.EstimatedJob {
	jobs := make([]*rpc.EstimatedJob, 0)
	for _, ej := range eJobs {
		j := &rpc.EstimatedJob{
			Id: ej.ID,
			WorkerToPickupEstimate: &rpc.Estimate{
				DistanceMeters: int32(ej.WorkerToPickup.DistanceMeters),
				Duration:       durationpb.New(ej.WorkerToPickup.Duration),
			},
			PickupToDropOffEstimate: &rpc.Estimate{
				DistanceMeters: int32(ej.PickupToDropOff.DistanceMeters),
				Duration:       durationpb.New(ej.PickupToDropOff.Duration),
			},
			WorkerLocation: &rpc.Location{
				Lat:     ej.WorkerLocation.Lat,
				Lng:     ej.WorkerLocation.Lng,
				Address: ej.WorkerLocation.Address,
			},
			PickupLocation: &rpc.Location{
				Lat:     ej.PickupLocation.Lat,
				Lng:     ej.PickupLocation.Lng,
				Address: ej.PickupLocation.Address,
			},
			DropOffLocation: &rpc.Location{
				Lat:     ej.DropOffLocation.Lat,
				Lng:     ej.DropOffLocation.Lng,
				Address: ej.DropOffLocation.Address,
			},
		}
		jobs = append(jobs, j)
	}
	return jobs
}

func toChunks(jobs []*job.AvailableJob, chunkSize int) [][]*job.AvailableJob {
	var r [][]*job.AvailableJob
	for i := 0; i < len(jobs); i += chunkSize {
		end := i + chunkSize
		if end > len(jobs) {
			end = len(jobs)
		}
		r = append(r, jobs[i:end])
	}

	return r
}
