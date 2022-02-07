// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package cfg

import (
	"fmt"

	"google.golang.org/grpc"
)

// GrpcConfig defines gRPC server configuration params.
type GrpcConfig struct {
	Addr    string `default:":9090" usage:"gRPC server |host:port| listening address"`
	Timeout struct {
		Connect uint `default:"120" usage:"timeout in |seconds| to establish a new connection (including HTTP/2 handshaking)"`
	}
	Workers uint `yaml:",omitempty" default:"0" usage:"|number| of worker goroutines that should be used to process incoming streams"`
}

// Check validates GrpcConfig params.
func (c *GrpcConfig) Check(name ...string) (err error) {
	if err = checkAddr(c.Addr, true, minHttpPort, name); err == nil {
		if c.Timeout.Connect < 1 {
			err = fmt.Errorf("%s is too small", field(append(name, "timeout", "connect")))
		}
	}
	return
}

// ServerOptions returns an array of gRPC server configuration options.
func (c *GrpcConfig) ServerOptions() []grpc.ServerOption {
	t := c.Timeout
	opts := []grpc.ServerOption{
		grpc.ConnectionTimeout(usec(t.Connect)),
		grpc.NumStreamWorkers(uint32(c.Workers)),
	}
	return opts
}
