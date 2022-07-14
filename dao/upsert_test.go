// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package dao

import (
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