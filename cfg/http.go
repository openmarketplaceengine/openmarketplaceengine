// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package cfg

import (
	"fmt"
	"time"
)

const minHttpPort = 80 //nolint

// HttpConfig contains HTTP server configuration params.
type HttpConfig struct { //nolint
	Bind string `default:"0.0.0.0" usage:"HTTP server local bind #address#"`
	Name string `default:"OME/1.0" usage:"HTTP #server name# response header"`
	Port int    `default:"8080" usage:"HTTP server #port#"`
	Path string `usage:"HTTP public routable #path#"`
	// HTTP server timeouts
	//
	// See: https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
	Timeout struct {
		Idle  uint `default:"30" usage:"HTTP request idle keep-alive timeout in #seconds#"`
		Read  uint `default:"10" usage:"HTTP request read timeout in #seconds#"`
		Write uint `default:"10" usage:"HTTP response write timeout in #seconds#"`
		Stop  uint `default:"20" usage:"HTTP server graceful shutdown timeout in #seconds#"`
	}
}

func (c *HttpConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Bind, c.Port)
}

func (c *HttpConfig) IdleTimeout() time.Duration {
	return usec(c.Timeout.Idle)
}

func (c *HttpConfig) ReadTimeout() time.Duration {
	return usec(c.Timeout.Read)
}

func (c *HttpConfig) WriteTimeout() time.Duration {
	return usec(c.Timeout.Write)
}

func (c *HttpConfig) StopTimeout() time.Duration {
	return usec(c.Timeout.Stop)
}

// Check validates HttpConfig params.
func (c *HttpConfig) Check(name ...string) (err error) {
	err = checkPort(c.Port, minHttpPort, append(name, "port"))
	return
}
