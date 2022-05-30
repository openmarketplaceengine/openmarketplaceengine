// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package job

import (
	"context"

	"github.com/openmarketplaceengine/openmarketplaceengine/dao"
	"github.com/openmarketplaceengine/openmarketplaceengine/dom/job"
	"github.com/openmarketplaceengine/openmarketplaceengine/internal/api/job/v1beta1"
	"github.com/openmarketplaceengine/openmarketplaceengine/srv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type controller struct {
	v1beta1.UnimplementedJobServiceServer
}

func init() {
	srv.Grpc.Register(func(s *grpc.Server) error {
		srv.Grpc.Infof("registering: %s", v1beta1.JobService_ServiceDesc.ServiceName)
		v1beta1.RegisterJobServiceServer(s, &controller{})
		return nil
	})
}

func (s *controller) ImportJob(ctx context.Context, req *v1beta1.ImportJobRequest) (*v1beta1.ImportJobResponse, error) {
	var act = v1beta1.JobAction_JOB_ACTION_CREATED
	var j job.Job
	s.setJob(&j, req.Job)
	_, ups, err := j.Upsert(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if ups == dao.UpsertUpdated {
		act = v1beta1.JobAction_JOB_ACTION_UPDATED
	}
	res := &v1beta1.ImportJobResponse{Action: act}
	return res, nil
}

func (s *controller) setJob(job *job.Job, req *v1beta1.JobInfo) {
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
