// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package dao

import (
	"database/sql"

	"github.com/driverscooperative/geosrv/cfg"
)

//-----------------------------------------------------------------------------

func Conn(ctx Context) (*sql.Conn, error) {
	if Running() {
		return Pgdb.sdb.Conn(ctx)
	}
	return nil, Pgdb.stateError()
}

//-----------------------------------------------------------------------------

func WithConn(ctx Context, run func(ctx Context, con *sql.Conn) error) (err error) {
	var con *sql.Conn
	con, err = Conn(ctx)
	if err != nil {
		return
	}
	defer FreeConn(con)
	defer cfg.CatchPanic(&err)
	return run(ctx, con)
}

//-----------------------------------------------------------------------------

func WithTran(ctx Context, opt *sql.TxOptions, run func(ctx Context, tx *sql.Tx) error) (err error) {
	var con *sql.Conn
	con, err = Conn(ctx)
	if err != nil {
		return
	}
	defer FreeConn(con)
	var tx *sql.Tx
	tx, err = con.BeginTx(ctx, opt)
	if err != nil {
		return
	}
	defer txfinish(tx, &err)
	defer cfg.CatchPanic(&err)
	return run(ctx, tx)
}

//-----------------------------------------------------------------------------

func FreeConn(con *sql.Conn) {
	logerr(con.Close())
}

//-----------------------------------------------------------------------------

func txfinish(tx *sql.Tx, err *error) {
	if *err != nil {
		logerr(txnoerr(tx.Rollback()))
		return
	}
	*err = tx.Commit()
}

func txnoerr(err error) error {
	if err != nil {
		switch err {
		case sql.ErrTxDone:
			return nil
		}
	}
	return err
}
