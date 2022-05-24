// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package status

import (
	"context"
	"time"

	uptimev1 "github.com/openmarketplaceengine/openmarketplaceengine/internal/api/status/v1"
	"github.com/openmarketplaceengine/openmarketplaceengine/srv"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type controller struct {
	uptimev1.UnimplementedUptimeServiceServer
	start time.Time
}

func init() {
	srv.Grpc.Register(func(s *grpc.Server) error {
		srv.Grpc.Infof("registering: %s", uptimev1.UptimeService_ServiceDesc.ServiceName)
		uptimev1.RegisterUptimeServiceServer(s, &controller{start: time.Now()})
		return nil
	})
}

func (u *controller) GetUptime(context.Context, *uptimev1.GetUptimeRequest) (*uptimev1.GetUptimeResponse, error) {
	res := new(uptimev1.GetUptimeResponse)
	res.Uptime = int64(time.Since(u.start))
	res.Started = timestamppb.New(u.start)
	return res, nil
}
