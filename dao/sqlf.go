// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package dao

import (
	"github.com/leporo/sqlf"
	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
)

type (
	SQL = sqlf.Stmt
)

func init() {
	sqlf.SetDialect(sqlf.PostgreSQL)
}

//-----------------------------------------------------------------------------

func NewSQL(verb string, args ...interface{}) *SQL {
	return sqlf.New(verb, args...)
}

//-----------------------------------------------------------------------------

func Exec(sql ...*SQL) (err error) {
	if failInit(&err) {
		return
	}
	ctx := cfg.Context()
	for i := 0; i < len(sql) && err == nil; i++ {
		_, err = sql[i].ExecAndClose(ctx, DB())
	}
	return
}
