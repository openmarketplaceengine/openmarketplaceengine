// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package dao

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
	"github.com/openmarketplaceengine/openmarketplaceengine/log"
)

type PgdbConn struct {
	cfg *pgx.ConnConfig
	sdb *sql.DB
	log log.Logger
}

var ErrNotStarted = errors.New("Pgdb not started")

var Pgdb = new(PgdbConn)

//-----------------------------------------------------------------------------

func DB() *sql.DB {
	return Pgdb.sdb
}

//-----------------------------------------------------------------------------

func (p *PgdbConn) Boot() (err error) {

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

	p.log = log.Named("PGDB")

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

	err = p.sdb.PingContext(cfg.Context())

	if err != nil {
		p.sdb = nil
		return
	}

	if len(pcfg.Schema) > 0 {
		err = p.SwitchSchema(pcfg.Schema)
		if err != nil {
			_ = p.sdb.Close()
			p.sdb = nil
			return
		}
		infof("using schema %q", pcfg.Schema)
	}

	return nil
}

//-----------------------------------------------------------------------------

func (p *PgdbConn) Stop() error {
	if p.sdb == nil {
		return cfg.CantStop("Pgdb")
	}
	return p.sdb.Close()
}

//-----------------------------------------------------------------------------

func (p *PgdbConn) SwitchSchema(name string) error {
	var rs RawSQL
	rs.Appendf("CREATE SCHEMA IF NOT EXISTS %q", name)
	rs.Appendf("SET search_path TO %q", name)
	return rs.Exec()
}

//-----------------------------------------------------------------------------

func (p *PgdbConn) Log(ctx context.Context, level pgx.LogLevel, msg string, data map[string]interface{}) {
	if level == pgx.LogLevelNone || p.log == nil {
		return
	}
	lev := matchLevel(level)
	if !p.log.IsLevel(lev) {
		return
	}
	if len(data) > 0 {
		p.log.Levelf(lev, "%s\n%s", msg, log.YAML(data))
		return
	}
	p.log.Levelf(lev, "%s", msg)
}

//-----------------------------------------------------------------------------

func matchLevel(level pgx.LogLevel) log.Level {
	switch level {
	case pgx.LogLevelTrace, pgx.LogLevelDebug:
		return log.LevelDebug
	case pgx.LogLevelInfo:
		return log.LevelInfo
	case pgx.LogLevelWarn:
		return log.LevelWarn
	}
	return log.LevelError
}

//-----------------------------------------------------------------------------

func failInit(err *error) bool {
	if Pgdb.sdb == nil {
		if err != nil {
			*err = ErrNotStarted
		}
		return true
	}
	return false
}

//-----------------------------------------------------------------------------

func infof(format string, args ...interface{}) {
	if Pgdb.log != nil {
		Pgdb.log.Infof(format, args...)
	}
}
