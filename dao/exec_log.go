// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package dao

//-----------------------------------------------------------------------------

type loggingExecutor struct {
	exe Executor
	opt LogOpt
}

func (e *loggingExecutor) logSQL(verb string, query string, args []interface{}) {
	if len(args) > 0 {
		infof("[%s] %s\n[ARGS] %#v", verb, query, args)
		return
	}
	infof("[%s] %s", verb, query)
}

func (e *loggingExecutor) logErr(verb string, query string, args []interface{}, err error) {
	if len(args) > 0 {
		errorf("%v\n%[s] %s\n[ARGS] %#v", err, verb, query, args)
		return
	}
	errorf("%v\n%[s] %s", err, verb, query)
}

func (e *loggingExecutor) ExecContext(ctx Context, query string, args ...interface{}) (Result, error) {
	const verb = "EXEC"
	if (e.opt & LogSQL) != 0 {
		e.logSQL(verb, query, args)
	}
	res, err := e.exe.ExecContext(ctx, query, args...)
	if (err != nil) && ((e.opt & LogErr) != 0) {
		e.logErr(verb, query, args, err)
	}
	return res, err
}

func (e *loggingExecutor) QueryContext(ctx Context, query string, args ...interface{}) (*Rows, error) {
	const verb = "QUERY"
	if (e.opt & LogSQL) != 0 {
		e.logSQL(verb, query, args)
	}
	rows, err := e.exe.QueryContext(ctx, query, args...)
	if (err != nil) && ((e.opt & LogErr) != 0) {
		e.logErr(verb, query, args, err)
	}
	return rows, err
}

func (e *loggingExecutor) QueryRowContext(ctx Context, query string, args ...interface{}) *Row {
	const verb = "QUERY"
	if (e.opt & LogSQL) != 0 {
		e.logSQL(verb, query, args)
	}
	return e.exe.QueryRowContext(ctx, query, args...)
}
