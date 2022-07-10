// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package dao

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDrop(t *testing.T) {
	const (
		tname = "drop_table"
		iname = "drop_index"
		vname = "drop_view"
	)

	WillTest(t, "test")

	var exec ListExec
	var drop Drop

	exec.Append(SQLExecf("create table if not exists %s (id int)", tname))
	exec.Append(SQLExecf("create index if not exists %s on %s (id)", iname, tname))
	exec.Append(SQLExecf("create or replace view %s as select * from %s", vname, tname))

	drop.AppendView(vname)
	drop.AppendIndex(iname)
	drop.AppendTable(tname)

	exec.Append(&drop)

	Pgdb.SetLogOpt(LogAll)

	err := ExecTX(context.Background(), &exec)

	require.NoError(t, err)
}
