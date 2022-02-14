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

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
	"github.com/openmarketplaceengine/openmarketplaceengine/log"
)

type HttpServer struct {
	chi.Router
	lsn net.Listener
	srv *http.Server
}

//-----------------------------------------------------------------------------

func NewHttpServer() *HttpServer {
	s := new(HttpServer)
	s.Router = chi.NewMux()
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
	log.Infof("HTTP listening on %s", addr)
	ctx := cfg.Context()
	s.srv = &http.Server{
		Addr: addr,
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
		ErrorLog:     log.NewStdLog(log.LevelError),
		Handler:      s.Router,
		IdleTimeout:  c.IdleTimeout(),
		ReadTimeout:  c.ReadTimeout(),
		WriteTimeout: c.WriteTimeout(),
	}
	s.srv.RegisterOnShutdown(func() {
		log.Infof("HTTP server shutdown")
	})
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
		log.Errorf("HTTP serving failed: %s", err)
	}
	cfg.Context().Stop()
}

//-----------------------------------------------------------------------------

func skipClosedError(err error) error {
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}
