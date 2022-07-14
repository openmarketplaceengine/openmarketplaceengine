// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package dao

import (
	"encoding/json"
	"unsafe"
)

//-----------------------------------------------------------------------------

type loggingExecutor struct {
	exe Executor
	opt LogOpt
}

type Flog = func(format string, args ...interface{})

func (e *loggingExecutor) ExecContext(ctx Context, query string, args ...interface{}) (Result, error) {
	const verb = "EXEC"
	if e.opt.LogSQL() {
		logSQL(debugf, verb, query, args)
	}
	res, err := e.exe.ExecContext(ctx, query, args...)
	if err != nil && e.opt.LogErr() && !ShouldSkipError(ctx, err) {
		logErr(errorf, verb, query, args, err)
	}
	return res, err
}

func (e *loggingExecutor) QueryContext(ctx Context, query string, args ...interface{}) (*Rows, error) {
	const verb = "STMT"
	if e.opt.LogSQL() {
		logSQL(debugf, verb, query, args)
	}
	rows, err := e.exe.QueryContext(ctx, query, args...)
	if err != nil && e.opt.LogErr() && !ShouldSkipError(ctx, err) {
		logErr(errorf, verb, query, args, err)
	}
	return rows, err
}

func (e *loggingExecutor) QueryRowContext(ctx Context, query string, args ...interface{}) *Row {
	const verb = "STMT"
	if e.opt.LogSQL() {
		logSQL(debugf, verb, query, args)
	}
	return e.exe.QueryRowContext(ctx, query, args...)
}

//-----------------------------------------------------------------------------

func jsonArgs(args []interface{}) string {
	buf, err := json.Marshal(args)
	if err != nil {
		return "ARGS JSON ERROR: " + err.Error()
	}
	return *(*string)(unsafe.Pointer(&buf))
}

const _argsIndent = "                                     "

func logSQL(flog Flog, verb string, query string, args []interface{}) {
	if len(args) > 0 {
		flog("[%s] %s\n%s[ARGS] %s", verb, query, _argsIndent, jsonArgs(args))
		return
	}
	flog("[%s] %s", verb, query)
}

func logErr(flog Flog, verb string, query string, args []interface{}, err error) {
	if len(args) > 0 {
		flog("%v\n%s[%s] %s\n%s[ARGS] %s", err, _argsIndent, verb, query, _argsIndent, jsonArgs(args))
		return
	}
	flog("%v\n%s[%s] %s", err, _argsIndent, verb, query)
}

//-----------------------------------------------------------------------------

type flogExec struct {
	flog Flog
}

func (e *flogExec) ExecContext(ctx Context, query string, args ...interface{}) (Result, error) {
	logSQL(e.flog, "EXEC", query, args)
	return nil, nil
}

func (e *flogExec) QueryContext(ctx Context, query string, args ...interface{}) (*Rows, error) {
	logSQL(e.flog, "STMT", query, args)
	return nil, nil
}

func (e *flogExec) QueryRowContext(ctx Context, query string, args ...interface{}) *Row {
	logSQL(e.flog, "STMT", query, args)
	return nil
}

//-----------------------------------------------------------------------------

func ExecFlog(ctx Context, flog Flog, execs ...Executable) (err error) {
	if flog == nil {
		flog = debugf
	}
	exe := flogExec{flog}
	for i := 0; i < len(execs) && err == nil; i++ {
		err = execs[i].Execute(ctx, &exe)
	}
	return
}
