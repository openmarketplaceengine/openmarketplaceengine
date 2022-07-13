// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package dao

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDrop(t *testing.T) {
	const (
		tname = "drop_table"
		iname = "drop_index"
		vname = "drop_view"
	)

	ctx := WillTest(t, "test")

	var exec ListExec
	var drop Drop

	exec.Append(CreateTable(tname, "id int"))
	exec.Append(SQLExecf("create index if not exists %s on %s (id)", iname, tname))
	exec.Append(SQLExecf("create or replace view %s as select * from %s", vname, tname))

	drop.AppendView(vname)
	drop.AppendIndex(iname)
	drop.AppendTable(tname)

	exec.Append(&drop)

	Pgdb.SetLogOpt(LogAll)

	err := ExecTX(ctx, &exec)

	require.NoError(t, err)
}

func TestDropColumn(t *testing.T) {
	const table = "drop_column_test"
	ctx := WillTest(t, "test")
	var exe ListExec
	exe.Append(CreateTable(table, "c1 int", "c2 int", "c3 int"))
	exe.Append(DropColumn(table, "c1"))
	exe.Append(DropColumn(table, "c2"))
	exe.Append(DropColumn(table, "c3"))
	exe.Append(DropTable(table, true))
	err := ExecTX(ctx, &exe)
	require.NoError(t, err)
}