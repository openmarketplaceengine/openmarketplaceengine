// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package status

import (
	"context"
	"time"

	uptimev1 "github.com/openmarketplaceengine/openmarketplaceengine/api/gen/status/v1"
	"github.com/openmarketplaceengine/openmarketplaceengine/srv"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type uptimeService struct {
	uptimev1.UnimplementedUptimeServiceServer
	start time.Time
}

//-----------------------------------------------------------------------------

func init() {
	srv.Grpc.Register(func(s *grpc.Server) error {
		srv.Grpc.Infof("registering: %s", uptimev1.UptimeService_ServiceDesc.ServiceName)
		uptimev1.RegisterUptimeServiceServer(s, &uptimeService{start: time.Now()})
		return nil
	})
}

//-----------------------------------------------------------------------------

func (u *uptimeService) GetUptime(context.Context, *emptypb.Empty) (*uptimev1.UptimeResponse, error) {
	res := new(uptimev1.UptimeResponse)
	res.Uptime = int64(time.Since(u.start))
	res.Started = timestamppb.New(u.start)
	return res, nil
}
