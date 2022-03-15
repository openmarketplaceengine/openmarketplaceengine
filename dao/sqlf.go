// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package dao

import (
	"github.com/leporo/sqlf"
)

type SQL struct {
	*sqlf.Stmt
}

func init() {
	sqlf.SetDialect(sqlf.PostgreSQL)
}

//-----------------------------------------------------------------------------

func NewSQL(verb string, args ...interface{}) *SQL {
	return &SQL{sqlf.New(verb, args...)}
}

//-----------------------------------------------------------------------------

func (s *SQL) Execute(ctx Context, exe Executor) error {
	_, err := s.ExecAndClose(ctx, exe)
	return err
}

//-----------------------------------------------------------------------------

func Insert(table string) *SQL {
	return &SQL{sqlf.InsertInto(table)}
}

//-----------------------------------------------------------------------------

func Update(table string) *SQL {
	return &SQL{sqlf.Update(table)}
}

//-----------------------------------------------------------------------------

func Delete(table string) *SQL {
	return &SQL{sqlf.DeleteFrom(table)}
}

//-----------------------------------------------------------------------------

func Select(expr string, args ...interface{}) *SQL {
	return &SQL{sqlf.Select(expr, args...)}
}

//-----------------------------------------------------------------------------

func From(expr string, args ...interface{}) *SQL {
	return &SQL{sqlf.From(expr, args...)}
}
