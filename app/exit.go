// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package app

import "os"

type ExitCoder interface {
	ExitCode() int
}

func Exit(code int, errs ...error) {
	if n := len(errs); n > 0 {
		for i := 0; i < n; i++ {
			if err := errs[i]; err != nil {
				println(err.Error())
				if code == 0 {
					if ec, ok := err.(ExitCoder); ok {
						code = ec.ExitCode()
					}
				}
			}
		}
		if code == 0 {
			code = 1
		}
	}
	os.Exit(code)
}
