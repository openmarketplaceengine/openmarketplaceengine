// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package job

import (
	"context"

	"github.com/openmarketplaceengine/openmarketplaceengine/dao"
	"github.com/openmarketplaceengine/openmarketplaceengine/dom/job"
	jobV1 "github.com/openmarketplaceengine/openmarketplaceengine/internal/api/job/v1"
	"github.com/openmarketplaceengine/openmarketplaceengine/srv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type controller struct {
	jobV1.UnimplementedJobServiceServer
}

func init() {
	srv.Grpc.Register(func(s *grpc.Server) error {
		srv.Grpc.Infof("registering: %s", jobV1.JobService_ServiceDesc.ServiceName)
		jobV1.RegisterJobServiceServer(s, &controller{})
		return nil
	})
}

func (s *controller) ImportJob(ctx context.Context, req *jobV1.ImportJobRequest) (*jobV1.ImportJobResponse, error) {
	var act = jobV1.JobAction_JOB_ACTION_CREATED
	var j job.Job
	s.setJob(&j, req)
	_, ups, err := j.Upsert(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if ups == dao.UpsertUpdated {
		act = jobV1.JobAction_JOB_ACTION_UPDATED
	}
	res := &jobV1.ImportJobResponse{Action: act}
	return res, nil
}

func (s *controller) setJob(job *job.Job, req *jobV1.ImportJobRequest) {
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
