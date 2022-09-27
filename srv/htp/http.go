// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package htp

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/driverscooperative/geosrv/cfg"
	"github.com/driverscooperative/geosrv/log"
	"github.com/go-chi/chi/v5/middleware"
)

type HttpServer struct { //nolint
	Routes
	lsn net.Listener
	srv *http.Server
	log log.Logger
}

//-----------------------------------------------------------------------------

func NewHttpServer() *HttpServer { //nolint
	s := new(HttpServer)
	return s
}

//-----------------------------------------------------------------------------

func (s *HttpServer) Boot() (err error) {
	c := cfg.Http
	addr := c.Addr()
	s.lsn, err = net.Listen("tcp", addr)
	if err != nil {
		return
	}
	s.log = log.Named("HTTP")
	s.log.Infof("listening on %s", addr)
	ctx := cfg.Context()
	s.srv = &http.Server{
		Addr: addr,
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
		ErrorLog:     log.NewStdLog(log.LevelError),
		Handler:      s.Build(),
		IdleTimeout:  c.IdleTimeout(),
		ReadTimeout:  c.ReadTimeout(),
		WriteTimeout: c.WriteTimeout(),
	}
	go s.serve()
	return nil
}

//-----------------------------------------------------------------------------

func (s *HttpServer) Stop() error {
	if s.srv == nil {
		return fmt.Errorf("HTTP server not initialized")
	}
	wait := cfg.Http.StopTimeout()
	if wait == 0 {
		return s.srv.Close()
	}
	end, cancel := context.WithDeadline(context.Background(), time.Now().Add(wait))
	err := s.srv.Shutdown(end)
	cancel()
	if err != nil && err == context.DeadlineExceeded {
		log.Errorf("HTTP server shutdown deadline exceeded (%s). Force closing.", wait)
		return s.srv.Close()
	}
	return err
}

//-----------------------------------------------------------------------------

func (s *HttpServer) SetHeader(key, val string) {
	s.Use(middleware.SetHeader(key, val))
}

//-----------------------------------------------------------------------------

func (s *HttpServer) serve() {
	err := s.srv.Serve(s.lsn)
	if skipClosedError(err) != nil {
		s.log.Errorf("http.Serve failed: %s", err)
	}
}

//-----------------------------------------------------------------------------

func skipClosedError(err error) error {
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}
