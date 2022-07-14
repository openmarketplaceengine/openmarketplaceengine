// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package dao

import (
	"fmt"
	"strings"
)

type UpsertStatus int

const (
	UpsertUnknown UpsertStatus = 0
	UpsertCreated UpsertStatus = 1
	UpsertUpdated UpsertStatus = 2
)

func Upsert(ctx Context, insert, update func() Executable) (Result, UpsertStatus, error) {
	sql := insert()
	err := ExecTX(SkipErrorsContext(ctx, ErrUniqueViolation), sql)
	if err == nil && RowsAffected(sql.Result()) > 0 {
		return sql.Result(), UpsertCreated, nil
	}
	if err != nil && !ErrUniqueViolation.Is(err) {
		return nil, UpsertUnknown, err
	}
	sql = update()
	err = ExecTX(ctx, sql)
	if err == nil {
		return sql.Result(), UpsertUpdated, nil
	}
	return nil, UpsertUnknown, err
}

//-----------------------------------------------------------------------------
// Atomic Upsert
//-----------------------------------------------------------------------------

type colval struct {
	col string
	val interface{}
}

type AtomicUpsert struct {
	tbl string
	sql string
	key colval
	col []colval
	res Result
}

//-----------------------------------------------------------------------------

func UpsertAtomic(table string) *AtomicUpsert {
	return &AtomicUpsert{tbl: table, col: make([]colval, 0, 8)}
}

//-----------------------------------------------------------------------------

func (u *AtomicUpsert) Key(col string, val interface{}) *AtomicUpsert {
	u.key.col = col
	u.key.val = val
	return u
}

//-----------------------------------------------------------------------------

func (u *AtomicUpsert) Set(col string, val interface{}) *AtomicUpsert {
	u.col = append(u.col, colval{col, val})
	return u
}

//-----------------------------------------------------------------------------

func (u *AtomicUpsert) Put(col string, val interface{}) *AtomicUpsert {
	for i := 0; i < len(u.col); i++ {
		if c := &u.col[i]; c.col == col {
			c.val = val
			return u
		}
	}
	u.col = append(u.col, colval{col, val})
	return u
}

//-----------------------------------------------------------------------------

func (u *AtomicUpsert) build() string {
	var b strings.Builder
	b.WriteString("INSERT INTO ")
	b.WriteString(u.tbl)
	b.WriteString(" (")
	b.WriteString(u.key.col)
	n := len(u.col)
	for i := 0; i < n; i++ {
		b.WriteString(", ")
		b.WriteString(u.col[i].col)
	}
	b.WriteString(") VALUES ($1")
	for i := 0; i < n; i++ {
		b.WriteString(", ")
		_, _ = fmt.Fprintf(&b, "$%d", (i + 2))
	}
	b.WriteString(") ON CONFLICT (")
	b.WriteString(u.key.col)
	b.WriteString(") DO UPDATE SET ")
	for i := 0; i < n; i++ {
		if i > 0 {
			b.WriteString(", ")
		}
		_, _ = fmt.Fprintf(&b, "%s = $%d", u.col[i].col, (i + 2))
	}
	return b.String()
}

//-----------------------------------------------------------------------------

func (u *AtomicUpsert) Params() []interface{} {
	n := len(u.col) + 1
	p := make([]interface{}, n)
	p[0] = u.key.val
	for i := 1; i < n; i++ {
		p[i] = u.col[i-1].val
	}
	return p
}

//-----------------------------------------------------------------------------

func (u *AtomicUpsert) Clear() {
	u.sql = ""
	u.key.clear()
	u.res = nil
	if n := len(u.col); n > 0 {
		for i := 0; i < n; i++ {
			u.col[i].clear()
		}
		u.col = u.col[:0]
	}
}

//-----------------------------------------------------------------------------

func (u *AtomicUpsert) Table() string {
	return u.tbl
}

//-----------------------------------------------------------------------------

func (u *AtomicUpsert) String() string {
	if len(u.sql) == 0 {
		u.sql = u.build()
	}
	return u.sql
}

//-----------------------------------------------------------------------------

func (u *AtomicUpsert) Result() Result {
	return u.res
}

//-----------------------------------------------------------------------------

func (u *AtomicUpsert) Execute(ctx Context, exe Executor) (err error) {
	u.res, err = exe.ExecContext(ctx, u.String(), u.Params()...)
	return
}

//-----------------------------------------------------------------------------

func (u *AtomicUpsert) RowsAffected() int64 {
	return RowsAffected(u.res)
}

//-----------------------------------------------------------------------------

func (c *colval) clear() {
	c.col = ""
	c.val = nil
}
