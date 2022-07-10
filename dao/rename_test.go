// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package dao

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRenameTable(t *testing.T) {
	const (
		oldName = "rename_old"
		newName = "rename_new"
	)
	WillTest(t, "test")
	var exec ListExec
	exec.Append(SQLExecf("create table if not exists %s (id text)", oldName))
	exec.Append(RenameTable(oldName, newName))
	exec.Append(DropTable(newName, true))
	err := ExecTX(context.Background(), &exec)
	require.NoError(t, err)
}
