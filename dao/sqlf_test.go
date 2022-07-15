// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package dao

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const queryTestTable = "query_test"

func TestQueryOne(t *testing.T) {
	ctx := initQueryTestTable(t)
	var count int
	sql := From(queryTestTable).Select("count(*)").To(&count).Where("id > ?", 0)
	Pgdb.SetLogOpt(LogAll)
	has, err := sql.QueryOne(ctx)
	Pgdb.SetLogOpt(LogErr)
	require.NoError(t, err)
	require.True(t, has)
	require.Equal(t, 3, count)
}

//-----------------------------------------------------------------------------

func TestSQL_IgnoreConflict(t *testing.T) {
	ctx := initQueryTestTable(t)
	sql := Insert(queryTestTable).Set("id", 1).IgnoreConflict()
	err := ExecTX(ctx, sql)
	require.NoError(t, err)
}

//-----------------------------------------------------------------------------

func TestSQL_Count(t *testing.T) {
	ctx := initQueryTestTable(t)
	var cnt int
	sql := From(queryTestTable).Count("id").To(&cnt)
	Pgdb.SetLogOpt(LogAll)
	_, err := sql.QueryOne(ctx)
	Pgdb.SetLogOpt(LogErr)
	require.NoError(t, err)
	require.Equal(t, 3, cnt)
}

//-----------------------------------------------------------------------------

func initQueryTestTable(t *testing.T) Context {
	ctx := WillTest(t, "test")
	err := ExecDB(ctx, CreateTable(queryTestTable, "id int primary key not null", "name text"))
	require.NoError(t, err)
	dropTestTable(t, ctx, queryTestTable)
	fillQueryTestTable(t, ctx)
	return ctx
}

//-----------------------------------------------------------------------------

func fillQueryTestTable(t *testing.T, ctx Context) {
	err := ExecTX(ctx,
		Insert(queryTestTable).Set("id", 1).Set("name", "one"),
		Insert(queryTestTable).Set("id", 2).Set("name", "two"),
		Insert(queryTestTable).Set("id", 3).Set("name", "three"),
	)
	require.NoError(t, err)
}

//-----------------------------------------------------------------------------

func TestCoalesce(t *testing.T) {
	const cstable = "coalesce_test"
	ctx := WillTest(t, "test")
	err := ExecDB(ctx, CreateTable(cstable,
		"id int primary key not null",
		"num int",
		"str text",
		"tme timestamptz",
	))
	require.NoError(t, err)
	dropTestTable(t, ctx, cstable)
	ins1 := Insert(cstable).Set("id", 0)
	ins2 := Insert(cstable).Set("id", 1).Set("num", 1).Set("str", "one").Set("tme", time.Now())
	require.NoError(t, ExecTX(ctx, ins1, ins2))
	var num int
	var str string
	var tme = time.Now()
	require.False(t, tme.IsZero())
	sql := From(cstable)
	sql.Select(Coalesce("num", "0")).To(&num)
	sql.Select(Coalesce("str", "")).To(&str)
	sql.Select(Coalesce("tme", MakeTimestamptz(1, 1, 1, 0, 0, 0))).To(&tme)
	sql.Where("id = ?", 0)
	Pgdb.SetLogOpt(LogSQL | LogErr)
	_, err = sql.QueryOne(ctx)
	Pgdb.SetLogOpt(LogErr)
	require.NoError(t, err)
	require.True(t, tme.IsZero(), "timestamp is not zero")
}

//-----------------------------------------------------------------------------

func dropTestTable(t *testing.T, ctx Context, tableName string) {
	t.Cleanup(func() {
		err := ExecDB(ctx, DropTable(tableName, true))
		if err != nil {
			t.Error(err)
		}
	})
}