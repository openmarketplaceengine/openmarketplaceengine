// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package cfg

import (
	"context"
	"net"
	"net/http"
	"time"
)

const minHttpPort = 80 //nolint

// HttpConfig contains HTTP server configuration params.
type HttpConfig struct { //nolint
	Addr string `default:":9080" usage:"HTTP server |host:port| listening address"`
	Name string `default:"OME/1.0" usage:"HTTP server name response header"`
	// HTTP server timeouts
	//
	// See: https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
	Timeout struct {
		Idle  uint `default:"30" usage:"HTTP request idle keep-alive timeout in |seconds|"`
		Read  uint `default:"10" usage:"HTTP request read timeout in |seconds|"`
		Write uint `default:"10" usage:"HTTP response write timeout in |seconds|"`
		Stop  uint `default:"20" usage:"HTTP server graceful shutdown timeout in |seconds|"`
	}
}

// Apply sets appropriate http.Server properties.
func (c *HttpConfig) Apply(s *http.Server) *http.Server {
	s.Addr = c.Addr
	s.IdleTimeout = usec(c.Timeout.Idle)
	s.ReadTimeout = usec(c.Timeout.Read)
	s.WriteTimeout = usec(c.Timeout.Write)
	return s
}

// CreateServer creates and initializes new HTTP server instance.
func (c *HttpConfig) CreateServer(ctx context.Context, mux http.Handler) *http.Server {
	s := &http.Server{
		Handler: mux,
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
	}
	return c.Apply(s)
}

// Shutdown attempts to gracefully shutdown http.Server waiting for
// HttpConfig.Timeout.Stop seconds.
//
// Closes the http.Server immediately if HttpConfig.Timeout.Stop is zero.
func (c *HttpConfig) Shutdown(s *http.Server) error {
	if c.Timeout.Stop == 0 {
		return s.Close()
	}
	end, cancel := context.WithDeadline(context.Background(), time.Now().Add(usec(c.Timeout.Stop)))
	err := s.Shutdown(end)
	cancel()
	if err != nil && err == context.DeadlineExceeded {
		return s.Close()
	}
	return err
}

// Check validates HttpConfig params.
func (c *HttpConfig) Check(name ...string) (err error) {
	err = checkAddr(c.Addr, true, minHttpPort, name)
	return
}
