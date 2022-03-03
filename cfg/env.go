// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package cfg

import "os"

const trimChar = ' '

//-----------------------------------------------------------------------------

func init() {
	trimEnv()
}

//-----------------------------------------------------------------------------

func trimEnv() {
	env := os.Environ()
	for _, s := range env {
		if skipTrim(s) {
			continue
		}
		k, v := envKeyVal(fastTrim(s))
		_ = os.Setenv(k, v)
	}
}

//-----------------------------------------------------------------------------

func envKeyVal(env string) (key, val string) {
	for i := 0; i < len(env); i++ {
		if env[i] == '=' {
			key = env[:i]
			val = env[i+1:]
			return
		}
	}
	key = env
	return
}

//-----------------------------------------------------------------------------

func skipTrim(s string) bool {
	n := len(s)
	return n == 0 || (s[0] > trimChar && s[n-1] > trimChar)
}

//-----------------------------------------------------------------------------

func fastTrim(s string) string {
	end := -1
	for i := len(s); i > 0; i-- {
		if s[i-1] <= trimChar {
			end = i - 1
			continue
		}
		break
	}
	if end > -1 {
		s = s[:end]
	}
	for len(s) > 0 {
		if s[0] <= trimChar {
			s = s[1:]
			continue
		}
		break
	}
	return s
}
