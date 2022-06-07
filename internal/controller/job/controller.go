// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package job

import (
	"context"

	"github.com/openmarketplaceengine/openmarketplaceengine/dao"
	"github.com/openmarketplaceengine/openmarketplaceengine/dom/job"
	rpc "github.com/openmarketplaceengine/openmarketplaceengine/internal/api/job/v1beta1"
	typ "github.com/openmarketplaceengine/openmarketplaceengine/internal/api/type/v1beta1"
	"github.com/openmarketplaceengine/openmarketplaceengine/pkg/detector"
	"github.com/openmarketplaceengine/openmarketplaceengine/srv"
	svcJob "github.com/openmarketplaceengine/openmarketplaceengine/svc/job"
	"github.com/openmarketplaceengine/openmarketplaceengine/svc/location"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type controller struct {
	rpc.UnimplementedJobServiceServer
	jobService *svcJob.Service
}

func init() {
	srv.Grpc.Register(func(s *grpc.Server) error {
		srv.Grpc.Infof("registering: %s", rpc.JobService_ServiceDesc.ServiceName)

		storeClient := dao.Reds.StoreClient
		noOp := detector.NewDetectorNoOp()
		storage := location.NewStorage(storeClient)
		service := svcJob.NewService(location.NewTracker(storage, noOp))
		rpc.RegisterJobServiceServer(s, &controller{jobService: service})
		return nil
	})
}

//-----------------------------------------------------------------------------

func (c *controller) ImportJob(ctx context.Context, req *rpc.ImportJobRequest) (*rpc.ImportJobResponse, error) {
	var act = rpc.JobAction_JOB_ACTION_CREATED
	var j job.Job
	c.setJob(&j, req.Job)
	_, ups, err := j.Upsert(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if ups == dao.UpsertUpdated {
		act = rpc.JobAction_JOB_ACTION_UPDATED
	}
	res := &rpc.ImportJobResponse{Action: act}
	return res, nil
}

//-----------------------------------------------------------------------------

func (c *controller) ExportJob(ctx context.Context, req *rpc.ExportJobRequest) (*rpc.ExportJobResponse, error) {
	ids := req.Ids
	cnt := len(ids)
	if cnt == 0 {
		return nil, status.Error(codes.InvalidArgument, "empty job ids array")
	}
	jobs := make([]*rpc.ExportJobItem, cnt)
	for i := 0; i < cnt; i++ {
		jobID := ids[i]
		if len(jobID) == 0 {
			jobs[i] = &rpc.ExportJobItem{Id: "", Job: nil}
			continue
		}
		val, _, err := job.QueryOne(ctx, jobID)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed querying job %q: %v", jobID, err)
		}
		jobs[i] = &rpc.ExportJobItem{Id: jobID, Job: c.getJobInfo(val)}
	}
	return &rpc.ExportJobResponse{Jobs: jobs}, nil
}

func (c *controller) GetAvailableJobs(ctx context.Context, req *rpc.GetAvailableJobsRequest) (*rpc.GetAvailableJobsResponse, error) {
	// todo add validation using proto validation extension from Kevin

	availableJobs, err := c.jobService.GetAvailableJobs(ctx, req.GetAreaKey(), req.GetWorkerId(), req.GetRadiusMeters(), req.GetLimit())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get available jobs error: %v", err)
	}

	jobs := make([]*rpc.AvailableJob, 0)
	for _, aj := range availableJobs {
		j := &rpc.AvailableJob{
			Job: &rpc.JobInfo{
				Id:          aj.ID,
				WorkerId:    aj.WorkerID,
				Created:     timestamppb.New(aj.Created),
				Updated:     timestamppb.New(aj.Updated),
				State:       aj.State,
				PickupDate:  timestamppb.New(aj.PickupDate),
				PickupAddr:  aj.PickupAddr,
				PickupLoc:   &typ.Location{Latitude: aj.PickupLat, Longitude: aj.PickupLon},
				DropoffAddr: aj.DropoffAddr,
				DropoffLoc:  &typ.Location{Latitude: aj.DropoffLat, Longitude: aj.DropoffLon},
				TripType:    aj.TripType,
				Category:    aj.Category,
			},
			DistanceMeters: aj.DistanceMeters,
		}
		jobs = append(jobs, j)
	}

	return &rpc.GetAvailableJobsResponse{Jobs: jobs}, nil
}

//-----------------------------------------------------------------------------

func (c *controller) setJob(job *job.Job, req *rpc.JobInfo) {
	job.ID = req.Id
	job.WorkerID = req.WorkerId
	job.Created = req.Created.AsTime()
	job.Updated = req.Updated.AsTime()
	job.State = req.State
	job.PickupDate = req.PickupDate.AsTime()
	job.PickupAddr = req.PickupAddr
	job.PickupLat = req.PickupLoc.GetLatitude()
	job.PickupLon = req.PickupLoc.GetLongitude()
	job.DropoffAddr = req.DropoffAddr
	job.DropoffLat = req.DropoffLoc.GetLatitude()
	job.DropoffLon = req.DropoffLoc.GetLongitude()
	job.TripType = req.TripType
	job.Category = req.Category
}

//-----------------------------------------------------------------------------

func (c *controller) getJobInfo(job *job.Job) *rpc.JobInfo {
	if job == nil {
		return nil
	}
	inf := new(rpc.JobInfo)
	inf.Id = job.ID
	inf.WorkerId = job.WorkerID
	inf.Created = timestamppb.New(job.Created)
	inf.Updated = timestamppb.New(job.Updated)
	inf.State = job.State
	inf.PickupDate = timestamppb.New(job.PickupDate)
	inf.PickupAddr = job.PickupAddr
	inf.PickupLoc = &typ.Location{Latitude: job.PickupLat, Longitude: job.PickupLon}
	inf.DropoffAddr = job.DropoffAddr
	inf.DropoffLoc = &typ.Location{Latitude: job.DropoffLat, Longitude: job.DropoffLon}
	inf.TripType = job.TripType
	inf.Category = job.Category
	return inf
}
