// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package dao

import (
	"database/sql"
	"time"

	"github.com/leporo/sqlf"
)

type SQL struct {
	stmt *sqlf.Stmt
}

func init() {
	sqlf.SetDialect(sqlf.PostgreSQL)
}

//-----------------------------------------------------------------------------

func NewSQL(verb string, args ...interface{}) SQL {
	return SQL{sqlf.New(verb, args...)}
}

//-----------------------------------------------------------------------------

func (s SQL) Execute(ctx Context, exe Executor) error {
	_, err := s.stmt.ExecAndClose(ctx, exe)
	return err
}

//-----------------------------------------------------------------------------

func Insert(table string) SQL {
	return SQL{sqlf.InsertInto(table)}
}

//-----------------------------------------------------------------------------

func Update(table string) SQL {
	return SQL{sqlf.Update(table)}
}

//-----------------------------------------------------------------------------

func Delete(table string) SQL {
	return SQL{sqlf.DeleteFrom(table)}
}

//-----------------------------------------------------------------------------

func Select(expr string, args ...interface{}) SQL {
	return SQL{sqlf.Select(expr, args...)}
}

//-----------------------------------------------------------------------------

func From(expr string, args ...interface{}) SQL {
	return SQL{sqlf.From(expr, args...)}
}

//-----------------------------------------------------------------------------

func (s SQL) Set(field string, value interface{}) SQL {
	s.stmt.Set(field, value)
	return s
}

//-----------------------------------------------------------------------------

// SetNonZero binds an INSERT field if value is not zero.
func (s SQL) SetNonZero(field string, value interface{}) (self SQL) {
	self = s
	if value == nil {
		return
	}
	switch v := value.(type) {
	case string:
		if len(v) == 0 {
			return
		}
	case []byte:
		if len(v) == 0 {
			return
		}
	case int, int8, int16, int32, int64:
		if v == 0 {
			return
		}
	case uint, uint8, uint16, uint32, uint64, uintptr:
		if v == 0 {
			return
		}
	case time.Time:
		if v.IsZero() {
			return
		}
	}
	s.stmt.Set(field, value)
	return
}

//-----------------------------------------------------------------------------

func (s SQL) Bind(v interface{}) SQL {
	s.stmt.Bind(v)
	return s
}

//-----------------------------------------------------------------------------

func (s SQL) Select(expr string, args ...interface{}) SQL {
	s.stmt.Select(expr, args...)
	return s
}

//-----------------------------------------------------------------------------

func (s SQL) To(dest ...interface{}) SQL {
	s.stmt.To(dest...)
	return s
}

//-----------------------------------------------------------------------------

func (s SQL) Where(expr string, args ...interface{}) SQL {
	s.stmt.Where(expr, args...)
	return s
}

//-----------------------------------------------------------------------------

func (s SQL) OrderBy(expr ...string) SQL {
	s.stmt.OrderBy(expr...)
	return s
}

//-----------------------------------------------------------------------------

func (s SQL) Limit(limit interface{}) SQL {
	s.stmt.Limit(limit)
	return s
}

//-----------------------------------------------------------------------------
// Query
//-----------------------------------------------------------------------------

func (s SQL) QueryOne(ctx Context) error {
	return WithConn(ctx, func(ctx Context, con *sql.Conn) error {
		return s.stmt.QueryRowAndClose(ctx, con)
	})
}

//-----------------------------------------------------------------------------

func (s SQL) QueryEach(ctx Context, eachFunc func(rows *Rows)) error {
	return WithConn(ctx, func(ctx Context, con *sql.Conn) error {
		return s.stmt.QueryAndClose(ctx, con, eachFunc)
	})
}

//-----------------------------------------------------------------------------

func (s SQL) QueryRows(ctx Context, rowsFunc func(rows *Rows) error) error {
	return WithConn(ctx, func(ctx Context, con *sql.Conn) error {
		rows, err := con.QueryContext(ctx, s.stmt.String(), s.stmt.Args()...)
		s.stmt.Close()
		if err != nil {
			return err
		}
		err = rowsFunc(rows)
		if err != nil {
			_ = rows.Close()
			return err
		}
		err = rows.Close()
		if err != nil {
			return err
		}
		return rows.Err()
	})
}
