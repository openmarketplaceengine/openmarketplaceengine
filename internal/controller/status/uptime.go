// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package status

import (
	"context"
	"time"

	"github.com/openmarketplaceengine/openmarketplaceengine/internal/idl/api/status/v1beta1"
	"github.com/openmarketplaceengine/openmarketplaceengine/srv"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type controller struct {
	v1beta1.UnimplementedUptimeServiceServer
	start time.Time
}

func init() {
	srv.Grpc.Register(func(s *grpc.Server) error {
		srv.Grpc.Infof("registering: %s", v1beta1.UptimeService_ServiceDesc.ServiceName)
		v1beta1.RegisterUptimeServiceServer(s, &controller{start: time.Now()})
		return nil
	})
}

func (u *controller) GetUptime(context.Context, *v1beta1.GetUptimeRequest) (*v1beta1.GetUptimeResponse, error) {
	res := new(v1beta1.GetUptimeResponse)
	res.Uptime = int64(time.Since(u.start))
	res.Started = timestamppb.New(u.start)
	return res, nil
}
