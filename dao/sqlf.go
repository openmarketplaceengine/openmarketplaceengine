// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package dao

import (
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
