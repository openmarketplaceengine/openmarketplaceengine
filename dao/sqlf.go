// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package dao

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/leporo/sqlf"
)

type SQL struct {
	stmt   *sqlf.Stmt
	result Result
}

func init() {
	sqlf.SetDialect(sqlf.PostgreSQL)
}

//-----------------------------------------------------------------------------

func NewSQL(verb string, args ...interface{}) *SQL {
	return &SQL{sqlf.New(verb, args...), nil}
}

//-----------------------------------------------------------------------------

func (s *SQL) Execute(ctx Context, exe Executor) (err error) {
	s.result, err = s.stmt.ExecAndClose(ctx, exe)
	return
}

//-----------------------------------------------------------------------------

func (s *SQL) Result() Result {
	return s.result
}

//-----------------------------------------------------------------------------

func (s *SQL) RowsAffected() int64 {
	return RowsAffected(s.result)
}

//-----------------------------------------------------------------------------

func Insert(table string) *SQL {
	return &SQL{sqlf.InsertInto(table), nil}
}

//-----------------------------------------------------------------------------

func Update(table string) *SQL {
	return &SQL{sqlf.Update(table), nil}
}

//-----------------------------------------------------------------------------

func Delete(table string) *SQL {
	return &SQL{sqlf.DeleteFrom(table), nil}
}

//-----------------------------------------------------------------------------

func Select(expr string, args ...interface{}) *SQL {
	return &SQL{sqlf.Select(expr, args...), nil}
}

//-----------------------------------------------------------------------------

func From(expr string, args ...interface{}) *SQL {
	return &SQL{sqlf.From(expr, args...), nil}
}

//-----------------------------------------------------------------------------

func (s *SQL) Set(field string, value interface{}) *SQL {
	s.stmt.Set(field, value)
	return s
}

//-----------------------------------------------------------------------------

// SetNonZero binds an INSERT field if value is not zero.
func (s *SQL) SetNonZero(field string, value interface{}) (self *SQL) {
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

func (s *SQL) Bind(v interface{}) *SQL {
	s.stmt.Bind(v)
	return s
}

//-----------------------------------------------------------------------------

func (s *SQL) Select(expr string, args ...interface{}) *SQL {
	s.stmt.Select(expr, args...)
	return s
}

//-----------------------------------------------------------------------------

func (s *SQL) To(dest ...interface{}) *SQL {
	s.stmt.To(dest...)
	return s
}

//-----------------------------------------------------------------------------

func (s *SQL) Where(expr string, args ...interface{}) *SQL {
	s.stmt.Where(expr, args...)
	return s
}

//-----------------------------------------------------------------------------

func (s *SQL) OrderBy(expr ...string) *SQL {
	s.stmt.OrderBy(expr...)
	return s
}

//-----------------------------------------------------------------------------

func (s *SQL) Limit(limit interface{}) *SQL {
	s.stmt.Limit(limit)
	return s
}

func (s *SQL) Offset(offset interface{}) *SQL {
	s.stmt.Offset(offset)
	return s
}

//-----------------------------------------------------------------------------
// Query
//-----------------------------------------------------------------------------

func (s *SQL) QueryOne(ctx Context) (bool, error) {
	has := false
	err := s.QueryRows(ctx, func(rows *Rows) error {
		if rows.Next() {
			has = true
			if dest := s.stmt.Dest(); len(dest) > 0 {
				return rows.Scan(dest...)
			}
		}
		return nil
	})
	return has, err
}

//-----------------------------------------------------------------------------

func (s *SQL) QueryEach(ctx Context, eachFunc func(rows *Rows)) error {
	return WithConn(ctx, func(ctx Context, con *sql.Conn) error {
		exe := loggingExecutor{con, Pgdb.LogOpt()}
		return s.stmt.QueryAndClose(ctx, &exe, eachFunc)
	})
}

//-----------------------------------------------------------------------------

func (s *SQL) QueryRows(ctx Context, rowsFunc func(rows *Rows) error) error {
	return WithConn(ctx, func(ctx Context, con *sql.Conn) error {
		exe := loggingExecutor{con, Pgdb.LogOpt()}
		rows, err := exe.QueryContext(ctx, s.stmt.String(), s.stmt.Args()...)
		if err != nil {
			s.stmt.Close()
			return err
		}
		err = rowsFunc(rows)
		s.stmt.Close()
		if err != nil {
			_ = rows.Close()
			return err
		}
		err = rows.Err()
		if err != nil {
			_ = rows.Close()
			return err
		}
		return rows.Close()
	})
}

//-----------------------------------------------------------------------------
// SQL Functions
//-----------------------------------------------------------------------------

func Coalesce(column string, ifnull string) string {
	if len(ifnull) == 0 {
		ifnull = "''"
	}
	return fmt.Sprintf("COALESCE(%s, %s)", column, ifnull)
}
