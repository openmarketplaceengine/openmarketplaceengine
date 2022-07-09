// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package dao

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

const queryTestTable = "query_test"

func TestQueryOne(t *testing.T) {
	WillTest(t, "test")
	ctx := initQueryTestTable(t)
	Pgdb.SetLogOpt(LogSQL | LogErr)
	var count int
	sql := From(queryTestTable).Select("count(*)").To(&count).Where("id > ?", 0)
	has, err := sql.QueryOne(ctx)
	Pgdb.SetLogOpt(LogErr)
	require.NoError(t, err)
	require.True(t, has)
	require.Equal(t, 3, count)
}

//-----------------------------------------------------------------------------

func initQueryTestTable(t *testing.T) Context {
	ctx := context.Background()
	err := ExecDB(ctx, SQLExecf("create table if not exists %q (id int primary key not null, name text)", queryTestTable))
	require.NoError(t, err)
	t.Cleanup(func() {
		err := ExecDB(ctx, DropTable(queryTestTable, true))
		if err != nil {
			t.Error(err)
		}
	})
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
