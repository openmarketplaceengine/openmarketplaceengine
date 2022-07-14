// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package job

import (
	"context"

	"github.com/openmarketplaceengine/openmarketplaceengine/dom/job"
	rpc "github.com/openmarketplaceengine/openmarketplaceengine/internal/idl/api/job/v1beta1"
	"github.com/openmarketplaceengine/openmarketplaceengine/srv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type jobStateServer struct {
	rpc.UnimplementedJobStateServiceServer
}

//-----------------------------------------------------------------------------

func init() {
	srv.Grpc.Register(func(s *grpc.Server) error {
		srv.Grpc.Infof("registering: %s", rpc.JobStateService_ServiceDesc.ServiceName)
		rpc.RegisterJobStateServiceServer(s, new(jobStateServer))
		return nil
	})
}

//-----------------------------------------------------------------------------

func (s *jobStateServer) GetJobState(ctx context.Context, req *rpc.GetJobStateRequest) (*rpc.GetJobStateResponse, error) {
	state, found, err := job.GetState(ctx, req.JobId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	if !found {
		return nil, status.Errorf(codes.NotFound, "job not found: %q", req.JobId)
	}
	res := &rpc.GetJobStateResponse{
		JobId: req.JobId,
		State: rpc.JobState(state),
	}
	return res, nil
}

//-----------------------------------------------------------------------------

func (s *jobStateServer) UpdateJobState(ctx context.Context, req *rpc.UpdateJobStateRequest) (*rpc.UpdateJobStateResponse, error) {
	state, found := job.StateFromNumber(int(req.State))
	if !found || state == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid job state argument: %q", req.State.String())
	}
	set, err := job.SetState(ctx, req.JobId, state)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	if !set {
		return nil, status.Errorf(codes.NotFound, "job not found: %q", req.JobId)
	}
	res := &rpc.UpdateJobStateResponse{
		JobId: req.JobId,
		State: rpc.JobState(state),
	}
	return res, nil
}
