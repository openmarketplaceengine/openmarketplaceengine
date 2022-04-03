// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package rpc

import (
	"fmt"
	"net"

	"github.com/openmarketplaceengine/openmarketplaceengine/internal/service/tollgate/bbox"
	"github.com/openmarketplaceengine/openmarketplaceengine/internal/service/tollgate/detector"

	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
	locationV1beta1 "github.com/openmarketplaceengine/openmarketplaceengine/internal/omeapi/location/v1beta1"
	"github.com/openmarketplaceengine/openmarketplaceengine/internal/service/location"
	"github.com/openmarketplaceengine/openmarketplaceengine/log"
	redisClient "github.com/openmarketplaceengine/openmarketplaceengine/redis/client"
	"google.golang.org/grpc"
)

type GrpcServer struct {
	lsn net.Listener
	srv *grpc.Server
}

//-----------------------------------------------------------------------------

func NewGrpcServer() *GrpcServer {
	s := new(GrpcServer)
	return s
}

//-----------------------------------------------------------------------------

func (s *GrpcServer) Boot() (err error) {
	c := cfg.Grpc
	addr := c.Addr()
	s.lsn, err = net.Listen("tcp", addr)
	if err != nil {
		return
	}
	log.Infof("GRPC listening on %s", addr)
	s.srv = grpc.NewServer(s.configOptions()...)

	storeClient := redisClient.NewStoreClient()
	d, err := detector.NewDetector(cfg.Context(), bbox.NewStorage(storeClient))
	if err != nil {
		return
	}
	controller := location.New(storeClient, redisClient.NewPubSubClient(), "global", d)
	locationV1beta1.RegisterLocationServiceServer(s.srv, controller)

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

func (s *GrpcServer) serve() {
	err := s.srv.Serve(s.lsn)
	if err != nil {
		log.Errorf("GRPC serve error: %s", err)
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
