// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package dao

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestUpsert(t *testing.T) {
	const table = "upsert_test"
	ctx := WillTest(t, "test")
	dropTestTable(t, ctx, table)
	require.NoError(t, ExecTX(ctx, CreateTable(table, "id int primary key not null", "info text")))
	ins := func() Executable { return Insert(table).Set("id", 1).Set("info", "insert").IgnoreConflict() }
	upd := func() Executable { return Update(table).Set("info", "update").Where("id = ?", 1) }
	_, ups, err := Upsert(ctx, ins, upd)
	require.NoError(t, err)
	require.Equal(t, UpsertCreated, ups)
	_, ups, err = Upsert(ctx, ins, upd)
	require.NoError(t, err)
	require.Equal(t, UpsertUpdated, ups)
}

//-----------------------------------------------------------------------------

func TestAtomicUpsert(t *testing.T) {
	const table = "upsert_atomic"
	ctx := WillTest(t, "test")
	dropTestTable(t, ctx, table)
	require.NoError(t, ExecTX(ctx, CreateTable(table, "id int primary key not null", "info text", "stamp timestamptz")))
	logopt := Pgdb.SetLogOpt(LogAll)
	ups := UpsertAtomic(table)
	ups.Key("id", 1).Set("info", "info").Set("stamp", time.Now())
	require.NoError(t, ExecTX(ctx, ups))
	require.Equal(t, 1, int(ups.RowsAffected()))
	ups.Put("info", "upsert")
	require.NoError(t, ExecTX(ctx, ups))
	require.Equal(t, 1, int(ups.RowsAffected()))
	var info string
	sql := From(table).Select("info").To(&info).Where("id = ?", 1)
	has, err := sql.QueryOne(ctx)
	require.NoError(t, err)
	require.True(t, has)
	require.Equal(t, "upsert", info)
	Pgdb.SetLogOpt(logopt)
}

//-----------------------------------------------------------------------------
// Upsert statement building benchmarks
//-----------------------------------------------------------------------------

func benchUpsertBuild(b *testing.B, f func() string) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if s := f(); len(s) == 0 {
			b.Fatal("empty string")
		}
	}
}

func BenchmarkUpsertBuild(b *testing.B) {
	ups := UpsertAtomic("upsert_atomic")
	ups.Key("id", 1).Set("info", "info").Set("stamp", time.Now()).Set("test", nil)
	b.Run("Std", func(b *testing.B) {
		benchUpsertBuild(b, ups.buildStd)
	})
	b.Run("Buf", func(b *testing.B) {
		benchUpsertBuild(b, ups.build)
	})
}

func (u *AtomicUpsert) buildStd() string {
	var b strings.Builder
	b.Grow(256)
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
