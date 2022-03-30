// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package dao

import (
	"database/sql"
	"fmt"
)

var txExecOpts = sql.TxOptions{
	Isolation: sql.LevelReadCommitted,
	ReadOnly:  false,
}

//-----------------------------------------------------------------------------
// Executable Interfaces
//-----------------------------------------------------------------------------

// Executor interface depicts ExecContext, QueryContext, and QueryRowContext
// from sql.DB, sql.Conn, sql.Tx, sql.Stmt.
type Executor interface {
	ExecContext(ctx Context, query string, args ...interface{}) (Result, error)
	QueryContext(ctx Context, query string, args ...interface{}) (*Rows, error)
	QueryRowContext(ctx Context, query string, args ...interface{}) *Row
}

// Executable implementers perform actual Executor.ExecContext calls.
type Executable interface {
	Execute(ctx Context, exe Executor) error
	Result() Result
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
	return WithTran(ctx, &txExecOpts, func(ctx Context, tx *sql.Tx) (err error) {
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
	Stmt   string
	Args   []interface{}
	result Result
}

func NewArgsExec(stmt string, args ...interface{}) *ArgsExec {
	return &ArgsExec{stmt, args, nil}
}

func (a *ArgsExec) Execute(ctx Context, exe Executor) (err error) {
	if len(a.Args) > 0 {
		a.result, err = exe.ExecContext(ctx, a.Stmt, a.Args...)
	} else {
		a.result, err = exe.ExecContext(ctx, a.Stmt)
	}
	return
}

func (a *ArgsExec) Result() Result {
	return a.result
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

func (e SQLExec) Result() Result {
	return nil
}

func (e SQLExec) String() string {
	return string(e)
}

//-----------------------------------------------------------------------------

type ListExec struct {
	list []Executable
}

func (l *ListExec) Append(execs ...Executable) {
	l.list = append(l.list, execs...)
}

func (l *ListExec) Execute(ctx Context, exe Executor) (err error) {
	for i := 0; i < len(l.list) && err == nil; i++ {
		err = l.list[i].Execute(ctx, exe)
	}
	return
}

func (l *ListExec) Result() Result {
	return nil
}

func (l *ListExec) Slice() []Executable {
	return l.list
}

func (l *ListExec) Clear() {
	if n := len(l.list); n > 0 {
		for i := 0; i < n; i++ {
			l.list[i] = nil
		}
		l.list = l.list[:0]
	}
}

func (l *ListExec) Join(execs []Executable) []Executable {
	n1 := len(l.list)
	n2 := len(execs)
	if n1 == 0 {
		return execs
	}
	if n2 == 0 {
		return l.list
	}
	join := make([]Executable, 0, n1+n2)
	join = append(join, l.list...)
	join = append(join, execs...)
	return join
}
