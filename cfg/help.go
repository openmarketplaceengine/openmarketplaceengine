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

	"github.com/cristalhq/aconfig"
)

const envTag = "env"

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
		f.Usage = replByte(f.Usage, '#', '`')
	})
	c.flags.SetOutput(os.Stdout)
	fmt.Println("Usage:")
	c.flags.PrintDefaults()
}

//-----------------------------------------------------------------------------

func (c *ServerConfig) printEnviron() {
	if c.field == nil {
		return
	}
	// skip := map[string]bool{envFlag: true}
	c.exit = true
	temp := make([]string, 0, 8)
	fmt.Println("Environment Variables:")
	for _, f := range c.field {
		name := fieldEnvName(f, temp)
		fmt.Println(name)
	}
	// c.flags.VisitAll(func(f *flag.Flag) {
	// 	if skip[f.Name] {
	// 		return
	// 	}
	// 	name := replByte(f.Name, '.', '_')
	// 	name = EnvPrefix + "_" + strings.ToUpper(name)
	// 	if v, ok := os.LookupEnv(name); ok {
	// 		if envIsPassword(name) {
	// 			v = mockPass
	// 		}
	// 		fmt.Printf("  %s=%q\n", name, v)
	// 		return
	// 	}
	// 	fmt.Println(" ", name)
	// })
}

//-----------------------------------------------------------------------------

func fieldEnvName(f aconfig.Field, n []string) string {
	const envPfx = EnvPrefix + "_"
	p, ok := f.Parent()
	if !ok { // top level field
		return envPfx + f.Tag(envTag)
	}
	n = append(n[:0], f.Tag(envTag))
	n = fieldParentTag(p, envTag, n)
	n = append(n, EnvPrefix)
	for i, j := 0, len(n)-1; i < j; i, j = i+1, j-1 {
		n[i], n[j] = n[j], n[i]
	}
	return strings.Join(n, "_")
}

//-----------------------------------------------------------------------------

func fieldParentTag(f aconfig.Field, t string, n []string) []string {
	for ok := (f != nil); ok; {
		n = append(n, f.Tag(t))
		f, ok = f.Parent()
	}
	return n
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

func envIsPassword(name string) bool { //nolint
	return strings.HasSuffix(name, "_PASS") ||
		strings.HasSuffix(name, "_PASSWORD")
}
