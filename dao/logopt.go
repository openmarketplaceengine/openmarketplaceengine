// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package dao

import (
	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
)

type (
	LogOpt uint
	EnvOpt = cfg.EnvOpt
)

const (
	LogErr LogOpt = 1 << iota
	LogSQL
)

const (
	EnvLogErr EnvOpt = "PGDB_LOG_ERR"
	EnvLogSQL EnvOpt = "PGDB_LOG_SQL"
)

//-----------------------------------------------------------------------------

func GetEnvLogOpt() (opt LogOpt) {
	if EnvLogErr.GetBool() {
		opt |= LogErr
	}
	if EnvLogSQL.GetBool() {
		opt |= LogSQL
	}
	return
}

func SetEnvLogOpt(opt LogOpt) {
	EnvLogErr.SetBool(opt.LogErr())
	EnvLogSQL.SetBool(opt.LogSQL())
}

//-----------------------------------------------------------------------------

func (opt LogOpt) LogSQL() bool {
	return opt&LogSQL != 0
}

func (opt LogOpt) LogErr() bool {
	return opt&LogErr != 0
}

//-----------------------------------------------------------------------------
// Environment
//-----------------------------------------------------------------------------

func (opt LogOpt) EnvOpt() (EnvOpt, bool) {
	switch opt {
	case LogErr:
		return EnvLogErr, true
	case LogSQL:
		return EnvLogSQL, true
	}
	return "", false
}
