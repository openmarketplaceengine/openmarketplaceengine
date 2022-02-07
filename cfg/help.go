// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package cfg

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"unsafe"
)

//-----------------------------------------------------------------------------

func (c *ServerConfig) printHelp() {
	if c.flags == nil {
		return
	}
	c.exit = true
	checkConfFlag := true
	c.flags.VisitAll(func(f *flag.Flag) {
		if checkConfFlag && f.Name == confFlag {
			f.Usage = "config `file` path"
			checkConfFlag = false
			return
		}
		f.Usage = replByte(f.Usage, '|', '`')
	})
	c.flags.SetOutput(os.Stdout)
	fmt.Println("Usage:")
	c.flags.PrintDefaults()
}

//-----------------------------------------------------------------------------

func (c *ServerConfig) printEnviron() {
	if c.flags == nil {
		return
	}
	skip := map[string]bool{envFlag: true}
	c.exit = true
	fmt.Println("Environment Variables:")
	c.flags.VisitAll(func(f *flag.Flag) {
		if skip[f.Name] {
			return
		}
		name := replByte(f.Name, '.', '_')
		name = EnvPrefix + "_" + strings.ToUpper(name)
		if v, ok := os.LookupEnv(name); ok {
			if envIsPassword(name) {
				v = mockPass
			}
			fmt.Printf("  %s=%q\n", name, v)
			return
		}
		fmt.Println(" ", name)
	})
}

//-----------------------------------------------------------------------------

func replByte(s string, from byte, to byte) string {
	var b []byte
	for i := 0; i < len(s); i++ {
		if s[i] == from {
			if b == nil {
				b = []byte(s)
			}
			b[i] = to
		}
	}
	if b == nil {
		return s
	}
	return *(*string)(unsafe.Pointer(&b))
}

//-----------------------------------------------------------------------------

func envIsPassword(name string) bool {
	return strings.HasSuffix(name, "_PASS") ||
		strings.HasSuffix(name, "_PASSWORD")
}
