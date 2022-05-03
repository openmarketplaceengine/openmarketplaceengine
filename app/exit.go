// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package app

import "os"

type ExitCoder interface {
	ExitCode() int
}

func Exit(err error) {
	code := 0
	if err != nil {
		code = 1
		if ec, ok := err.(ExitCoder); ok {
			code = ec.ExitCode()
		}
		if msg := err.Error(); len(msg) > 0 {
			println(msg)
		}
	}
	os.Exit(code)
}
