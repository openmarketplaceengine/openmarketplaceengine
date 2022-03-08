// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package dao

import (
	"database/sql"
	"fmt"
)

//-----------------------------------------------------------------------------
// Executable Interfaces
//-----------------------------------------------------------------------------

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

//-----------------------------------------------------------------------------
// Executors
//-----------------------------------------------------------------------------

// ArgsExec executes SQL CRUD statements with parameters binding.
type ArgsExec struct {
	Stmt string
	Args []interface{}
}

func NewArgsExec(stmt string, args ...interface{}) *ArgsExec {
	return &ArgsExec{stmt, args}
}

func (a *ArgsExec) Execute(ctx Context, exe Executor) (err error) {
	if len(a.Args) > 0 {
		_, err = exe.ExecContext(ctx, a.Stmt, a.Args...)
	} else {
		_, err = exe.ExecContext(ctx, a.Stmt)
	}
	return
}

//-----------------------------------------------------------------------------

// SQLExec represents plain SQL statement without any parameter bindings.
//
// Attention! Do not use this executor for regular CRUD operations.
// It is only intended to be used for special statements where
// parameters binding do not work.
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
