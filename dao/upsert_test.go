// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package dao

import (
	"testing"

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
