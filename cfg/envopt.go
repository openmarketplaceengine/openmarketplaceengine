// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package cfg

import (
	"os"
	"strings"
)

type EnvOpt string

//-----------------------------------------------------------------------------

func (opt EnvOpt) Key() string {
	const pfx = EnvPrefix + "_"
	s := string(opt)
	if strings.HasPrefix(s, pfx) {
		return s
	}
	return pfx + s
}

func (opt EnvOpt) String() string {
	return opt.Key()
}

//-----------------------------------------------------------------------------
// Generic
//-----------------------------------------------------------------------------

func (opt EnvOpt) Get() string {
	return os.Getenv(opt.Key())
}

func (opt EnvOpt) Set(val string) error {
	return os.Setenv(opt.Key(), val)
}

func (opt EnvOpt) Del() error {
	return os.Unsetenv(opt.Key())
}

func (opt EnvOpt) Has() bool {
	_, ok := os.LookupEnv(opt.Key())
	return ok
}

//-----------------------------------------------------------------------------
// Bool
//-----------------------------------------------------------------------------

func (opt EnvOpt) GetBool() bool {
	switch opt.Get() {
	case "1", "true", "yes":
		return true
	default:
		return false
	}
}

func (opt EnvOpt) SetBool(value bool) bool {
	if value {
		return opt.Set("1") == nil
	}
	return opt.Del() == nil
}
