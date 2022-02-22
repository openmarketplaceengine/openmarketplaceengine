// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package cfg

import (
	"fmt"
	"time"
)

// GrpcConfig defines gRPC server configuration params.
type GrpcConfig struct {
	Bind    string `default:"0.0.0.0" usage:"gRPC server local bind #address#"`
	Port    int    `default:"8090" usage:"gRPC server #port#"`
	Timeout struct {
		Connect uint `default:"120" usage:"timeout in #seconds# to establish a new connection (including HTTP/2 handshaking)"`
		Stop    uint `default:"20" usage:"gRPC server graceful shutdown timeout in #seconds#"`
	}
	Workers uint `default:"0" usage:"#number# of worker goroutines that should be used to process incoming streams"`
}

// Check validates GrpcConfig params.
func (c *GrpcConfig) Check(name ...string) (err error) {
	if err = checkPort(c.Port, minHttpPort, append(name, "port")); err == nil {
		if c.Timeout.Connect < 1 {
			err = fmt.Errorf("%s is too small", field(name, "timeout", "connect"))
		}
	}
	return
}

func (c *GrpcConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Bind, c.Port)
}

func (c *GrpcConfig) ConnectTimeout() time.Duration {
	return usec(c.Timeout.Connect)
}

func (c *GrpcConfig) StopTimeout() time.Duration {
	return usec(c.Timeout.Stop)
}

func (c *GrpcConfig) StreamWorkers() uint32 {
	return uint32(c.Workers)
}
