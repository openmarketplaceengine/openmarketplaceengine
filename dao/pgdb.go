// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package dao

import (
	"context"
	"database/sql"
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
	"github.com/openmarketplaceengine/openmarketplaceengine/log"
)

type (
	Result  = sql.Result
	Rows    = sql.Rows
	Row     = sql.Row
	Context = context.Context
)

type PgdbConn struct {
	state cfg.State64
	cfg   *pgx.ConnConfig
	sdb   *sql.DB
	lopt  LogOpt
	drop  ListExec
	auto  ListExec
	upgr  upgradeManager
}

const (
	pfxErr = "pgdb" // error prefix
	pfxLog = "PGDB" // log prefix
)

var Pgdb = new(PgdbConn)
var plog = log.Log()

//-----------------------------------------------------------------------------

func DB() *sql.DB {
	return Pgdb.sdb
}

//-----------------------------------------------------------------------------

func (p *PgdbConn) Boot() (err error) {
	//
	if !p.state.TryBoot() {
		return p.stateError()
	}

	defer p.clearAutos()

	defer p.state.BootOrFail(&err)

	p.lopt = GetEnvLogOpt()

	if !EnvLogErr.Has() {
		p.lopt |= LogErr
	}

	plog = log.Named(pfxLog)

	pcfg := cfg.Pgdb()

	var addr string

	addr, err = pcfg.FullAddr()

	if err != nil {
		return
	}

	p.cfg, err = pgx.ParseConfig(addr)

	if err != nil {
		return
	}

	p.cfg.PreferSimpleProtocol = pcfg.Simple

	p.cfg.LogLevel = pgx.LogLevelNone

	p.sdb = stdlib.OpenDB(*p.cfg)

	if d := pcfg.MaxIdleTime(); d > 0 {
		p.sdb.SetConnMaxIdleTime(d)
	}

	if d := pcfg.MaxLifeTime(); d > 0 {
		p.sdb.SetConnMaxLifetime(d)
	}

	if n := pcfg.MaxIdleConns(); n != 0 {
		p.sdb.SetMaxIdleConns(n)
	}

	if n := pcfg.MaxOpenConns(); n != 0 {
		p.sdb.SetMaxOpenConns(n)
	}

	ctx := cfg.Context()

	err = p.sdb.PingContext(ctx)

	if err != nil {
		p.abort()
		return
	}

	p.state.SetRunning()

	if schema := pcfg.Schema; len(schema) > 0 {
		err = p.SwitchSchema(ctx, schema)
		if err != nil {
			p.abort()
			return
		}
		infof("using schema %q", schema)
	}

	if err = p.upgr.start(ctx); err != nil {
		p.abort()
		return
	}

	if err = p.autoExec(ctx); err != nil {
		p.abort()
		return
	}

	return nil
}

//-----------------------------------------------------------------------------

func (p *PgdbConn) autoExec(ctx Context) error {
	exec := p.drop.Join(p.auto.Slice())
	if len(exec) > 0 {
		return ExecTX(ctx, exec...)
	}
	return nil
}

//-----------------------------------------------------------------------------

func (p *PgdbConn) clearAutos() {
	p.drop.Clear()
	p.auto.Clear()
	p.upgr.clear()
}

//-----------------------------------------------------------------------------

func (p *PgdbConn) Stop() error {
	if p.state.TryStop() {
		return p.state.StopOrFail(p.sdb.Close)
	}
	return p.stateError()
}

//-----------------------------------------------------------------------------

func Running() bool {
	return Pgdb.state.Running()
}

//-----------------------------------------------------------------------------

func Invalid() bool {
	return Pgdb.state.Invalid()
}

//-----------------------------------------------------------------------------

func (p *PgdbConn) SwitchSchema(ctx Context, name string) error {
	return ExecDB(ctx,
		SQLExecf("CREATE SCHEMA IF NOT EXISTS %s", name),
		SQLExecf("SET search_path TO %s, public", name),
	)
}

//-----------------------------------------------------------------------------

func (p *PgdbConn) stateError() error {
	return p.state.StateError(pfxErr)
}

//-----------------------------------------------------------------------------

func (p *PgdbConn) abort() {
	if p.sdb != nil {
		_ = p.sdb.Close()
	}
}

//-----------------------------------------------------------------------------

func (p *PgdbConn) SetLogOpt(opt LogOpt) (old LogOpt) {
	old = p.lopt
	p.lopt = opt
	return
}

//-----------------------------------------------------------------------------

func (p *PgdbConn) LogOpt() LogOpt {
	return p.lopt
}

//-----------------------------------------------------------------------------

func isdebug() bool {
	return plog.IsDebug()
}

//-----------------------------------------------------------------------------

func debugf(format string, args ...interface{}) { //nolint:deadcode
	if isdebug() {
		plog.Debugf(format, args...)
	}
}

//-----------------------------------------------------------------------------

func infof(format string, args ...interface{}) {
	plog.Infof(format, args...)
}

//-----------------------------------------------------------------------------

func warnf(format string, args ...interface{}) { //nolint:deadcode
	plog.Warnf(format, args...)
}

//-----------------------------------------------------------------------------

func errorf(format string, args ...interface{}) {
	plog.Errorf(format, args...)
}

//-----------------------------------------------------------------------------

func logerr(err error, prefix ...string) {
	if err != nil {
		if len(prefix) > 0 {
			errorf("%s %s", strings.Join(prefix, " "), err)
			return
		}
		errorf("%s", err)
	}
}

//-----------------------------------------------------------------------------

func (p *PgdbConn) resetForTests(t Tester) {
	if Running() {
		infof("stopping...")
		err := Pgdb.Stop()
		if err != nil {
			t.Errorf("%s", err)
		}
	}
	p.state.SetUnused()
	p.cfg = nil
	p.sdb = nil
}
