// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package dao

import (
	"database/sql"
	"fmt"
)

// Executor interface depicts ExecContext from
// sql.DB, sql.Conn, sql.Tx, sql.Stmt.
type Executor interface {
	ExecContext(ctx Context, query string, args ...interface{}) (Result, error)
}

// Executable implementers perform actual Executor.ExecContext calls.
type Executable interface {
	Execute(ctx Context, exe Executor) error
}

//-----------------------------------------------------------------------------

// SQLExec represents plain SQL statement without any parameter bindings.
type SQLExec string

func SQLExecf(format string, args ...interface{}) SQLExec {
	return SQLExec(fmt.Sprintf(format, args...))
}

func (e SQLExec) Execute(ctx Context, exe Executor) error {
	_, err := exe.ExecContext(ctx, string(e))
	return err
}

func (e SQLExec) String() string {
	return string(e)
}

//-----------------------------------------------------------------------------
// Executable Runners
//-----------------------------------------------------------------------------

// ExecDB runs executables with sql.DB.
func ExecDB(ctx Context, execs ...Executable) error {
	return WithConn(ctx, func(ctx Context, con *sql.Conn) (err error) {
		for i := 0; i < len(execs) && err == nil; i++ {
			err = execs[i].Execute(ctx, con)
		}
		return
	})
}

// ExecTX runs executables with sql.Tx.
func ExecTX(ctx Context, execs ...Executable) error {
	return WithTran(ctx, func(ctx Context, tx *sql.Tx) (err error) {
		for i := 0; i < len(execs) && err == nil; i++ {
			err = execs[i].Execute(ctx, tx)
		}
		return
	})
}
