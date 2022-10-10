// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package rpc

import (
	"fmt"
	"net"

	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
	"github.com/openmarketplaceengine/openmarketplaceengine/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type RegFun = func(srv *grpc.Server) error

type GrpcServer struct {
	lsn net.Listener
	srv *grpc.Server
	reg []RegFun
}

var glog = log.Log()

//-----------------------------------------------------------------------------

func NewGrpcServer() *GrpcServer {
	s := new(GrpcServer)
	return s
}

//-----------------------------------------------------------------------------

func (s *GrpcServer) Boot() (err error) {
	glog = log.Named("GRPC")

	c := cfg.Grpc

	addr := c.Addr()

	s.lsn, err = net.Listen("tcp", addr)

	if err != nil {
		return
	}

	glog.Infof("listening on %s", addr)

	s.srv = grpc.NewServer(s.configOptions()...)

	if err = s.runreg(s.reg); err != nil {
		_ = s.lsn.Close()
		return
	}

	reflection.Register(s.srv)
	go s.serve()
	return nil
}

//-----------------------------------------------------------------------------

func (s *GrpcServer) Stop() error {
	if s.srv == nil {
		return fmt.Errorf("GRPC server not initialized")
	}
	s.srv.GracefulStop()
	return nil
}

//-----------------------------------------------------------------------------

// Register adds grpc.Server registration callback that is invoked
// after server creation.
func (s *GrpcServer) Register(fn RegFun) {
	if fn == nil {
		panic("Register function argument is nil")
	}
	s.reg = append(s.reg, fn)
}

//-----------------------------------------------------------------------------

func (s *GrpcServer) Debugf(format string, args ...interface{}) {
	glog.Debugf(format, args...)
}

func (s *GrpcServer) Infof(format string, args ...interface{}) {
	glog.Infof(format, args...)
}

func (s *GrpcServer) Warnf(format string, args ...interface{}) {
	glog.Warnf(format, args...)
}

func (s *GrpcServer) Errorf(format string, args ...interface{}) {
	glog.Errorf(format, args...)
}

//-----------------------------------------------------------------------------

func (s *GrpcServer) runreg(reg []RegFun) (err error) {
	for i := 0; i < len(reg) && err == nil; i++ {
		err = reg[i](s.srv)
	}
	return
}

//-----------------------------------------------------------------------------

func (s *GrpcServer) serve() {
	err := s.srv.Serve(s.lsn)
	if err != nil {
		glog.Errorf("serve error: %s", err)
	}
	cfg.Context().Stop()
}

//-----------------------------------------------------------------------------

// ServerOptions returns an array of gRPC server configuration options.
func (s *GrpcServer) configOptions() []grpc.ServerOption {
	c := cfg.Grpc
	opts := []grpc.ServerOption{
		grpc.ConnectionTimeout(c.ConnectTimeout()),
		grpc.NumStreamWorkers(c.StreamWorkers()),
	}
	return opts
}
