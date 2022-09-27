// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package srv

import (
	"github.com/driverscooperative/geosrv/srv/htp"
	"github.com/driverscooperative/geosrv/srv/rpc"
)

var (
	Http = htp.NewHttpServer() //nolint
	Grpc = rpc.NewGrpcServer()
)
