// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package cfg

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

const maxPort = (1 << 16) - 1

type checker interface {
	Check(name ...string) error
}

type namedCheck struct {
	name  string
	check checker
}

type checkList struct {
	list []namedCheck
}

func (c *checkList) add(name string, check checker) *checkList { //nolint
	c.list = append(c.list, namedCheck{name, check})
	return c
}

func (c *checkList) run() error {
	for i := range c.list {
		nc := &c.list[i]
		if err := nc.check.Check(nc.name); err != nil {
			return err
		}
	}
	return nil
}

//-----------------------------------------------------------------------------

func checkPort(port int, min int, name []string) error {
	if port < min || port > maxPort {
		return fmt.Errorf("invalid %s network port: %d", field(name), port)
	}
	return nil
}

//-----------------------------------------------------------------------------

func checkAddr(addr string, zeroHost bool, minPort int, name []string) error {
	if len(addr) == 0 {
		return fmt.Errorf("empty %s network address", field(name))
	}
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return fmt.Errorf("invalid %s %w", field(name), err)
	}
	if len(host) == 0 && !zeroHost {
		return fmt.Errorf("missing %s host address", field(name))
	}
	var pnum int
	pnum, err = strconv.Atoi(port)
	if err != nil {
		return fmt.Errorf("invalid %s network port: %q", field(name), port)
	}
	return checkPort(pnum, minPort, name)
}

//-----------------------------------------------------------------------------

func checkRange(num, min, max int, name []string) error {
	if num < min || num > max {
		return fmt.Errorf("invalid %s value %d", field(name), num)
	}
	return nil
}

//-----------------------------------------------------------------------------

func matchString(s string, allow ...string) bool {
	for i := range allow {
		if s == allow[i] {
			return true
		}
	}
	return false
}

//-----------------------------------------------------------------------------

func field(name []string, more ...string) string {
	if len(more) > 0 {
		name = append(name, more...)
	}
	return strings.Join(name, ".")
}

//-----------------------------------------------------------------------------

func usec(n uint) time.Duration {
	return time.Duration(n) * time.Second
}
