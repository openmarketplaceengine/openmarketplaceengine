// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package cfg

import (
	"fmt"
	"time"
)

// PgdbConfig represents PostgreSQL connection properties.
type PgdbConfig struct {
	Addr Address `usage:"database connection #URL#, e.g. postgres://user:pass@localhost:5432/dbname"`
	Log  string  `default:"info" usage:"connection log #level# [trace|debug|info|warn|error|none]"`
	Ssl  string  `default:"ignore" usage:"connection SSL #mode# [ignore|disable|allow|prefer|require|verify-ca|verify-full]"`
	Max  struct {
		Idletime uint `usage:"sets the maximum amount of time a connection may be idle in #seconds#"`
		Lifetime uint `usage:"sets the maximum amount of time a connection may be reused #seconds#"`
		Idleconn uint `default:"2" usage:"sets the maximum #number# of connections in the idle connection pool"`
		Openconn uint `default:"0" usage:"sets the maximum #number# of open connections to the database (0 is unlimited)"`
	}
	Schema string `default:"ome" usage:"database schema #name#"`
	Simple bool   `default:"false" usage:"prefer simple protocol"`
}

// Check validates PgdbConfig params.
func (c *PgdbConfig) Check(name ...string) error {
	if len(c.Addr) == 0 {
		return fmt.Errorf("%s: empty connection address", field(name, "addr"))
	}
	if !c.levelValid(c.Log) {
		return fmt.Errorf("%s: invalid log level: %q", field(name, "log"), c.Log)
	}
	if !matchString(c.Ssl, "ignore", "disable", "allow", "prefer", "require", "verify-ca", "verify-full") {
		return fmt.Errorf("%s: invalid SSL mode: %q", field(name, "ssl"), c.Ssl)
	}
	return nil
}

// FullAddr appends PgdbConfig.Addr with optional parameters.
func (c *PgdbConfig) FullAddr() string {
	const sslmode = "sslmode"
	a := c.Addr
	if len(c.Ssl) > 0 && c.Ssl != "ignore" {
		a = a.AppendQuery(sslmode, c.Ssl, false)
	}
	return string(a)
}

// MaxIdleTime returns the maximum amount of time a connection may be idle.
func (c *PgdbConfig) MaxIdleTime() time.Duration {
	return usec(c.Max.Idletime)
}

// MaxLifeTime returns the maximum amount of time a connection may be reused.
func (c *PgdbConfig) MaxLifeTime() time.Duration {
	return usec(c.Max.Lifetime)
}

// MaxIdleConns returns the maximum number of connections in the idle connection pool.
//
// If n <= 0, no idle connections are retained.
//
// The default max idle connections is currently 2. This may change in a future release.
func (c *PgdbConfig) MaxIdleConns() int {
	return int(c.Max.Idleconn)
}

// MaxOpenConns returns the maximum number of open connections to the database.
//
// If n <= 0, then there is no limit on the number of open connections. The default is 0 (unlimited).
func (c *PgdbConfig) MaxOpenConns() int {
	return int(c.Max.Openconn)
}

//-----------------------------------------------------------------------------

func (c *PgdbConfig) levelValid(level string) bool {
	switch level {
	case "trace", "debug", "info", "warn", "error", "none":
		return true
	}
	return false
}
