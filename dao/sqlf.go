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

func NewSQL(verb string, args ...interface{}) *SQL {
	return &SQL{sqlf.New(verb, args...)}
}

//-----------------------------------------------------------------------------

func (s *SQL) Execute(ctx Context, exe Executor) error {
	_, err := s.ExecAndClose(ctx, exe)
	return err
}
