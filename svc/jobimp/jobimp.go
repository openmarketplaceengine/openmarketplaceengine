// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package jobimp

import (
	"context"

	"github.com/openmarketplaceengine/openmarketplaceengine/dao"
	"github.com/openmarketplaceengine/openmarketplaceengine/dom"
	svc "github.com/openmarketplaceengine/openmarketplaceengine/internal/api/job/v1"
	"github.com/openmarketplaceengine/openmarketplaceengine/srv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type jobimpService struct {
	svc.UnimplementedJobServiceServer
}

//-----------------------------------------------------------------------------

func init() {
	srv.Grpc.Register(func(s *grpc.Server) error {
		srv.Grpc.Infof("registering: %s", svc.JobService_ServiceDesc.ServiceName)
		svc.RegisterJobServiceServer(s, &jobimpService{})
		return nil
	})
}

//-----------------------------------------------------------------------------

func (j *jobimpService) ImportJob(ctx context.Context, req *svc.ImportJobRequest) (*svc.ImportJobResponse, error) {
	var act = svc.JobAction_JOB_ACTION_CREATED
	var job dom.Jobimp
	j.setJob(&job, req)
	_, ups, err := dao.Upsert(ctx, job.Insert, job.Update)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if ups == dao.UpsertUpdated {
		act = svc.JobAction_JOB_ACTION_UPDATED
	}
	res := &svc.ImportJobResponse{Action: act}
	return res, nil
}

func (j *jobimpService) setJob(job *dom.Jobimp, req *svc.ImportJobRequest) {
	job.ID = req.Id
	job.WorkerID = req.WorkerId
	job.Created = req.Created.AsTime()
	job.Updated = req.Updated.AsTime()
	job.State = req.State
	job.PickupDate = req.PickupDate.AsTime()
	job.PickupAddr = req.PickupAddr
	job.PickupLat = req.PickupLoc.Lat
	job.PickupLon = req.PickupLoc.Lon
	job.DropoffAddr = req.DropoffAddr
	job.DropoffLat = req.DropoffLoc.Lat
	job.DropoffLon = req.DropoffLoc.Lon
	job.TripType = req.TripType
	job.Category = req.Category
}
