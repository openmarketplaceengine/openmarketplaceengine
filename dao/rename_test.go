// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package dao

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRenameTable(t *testing.T) {
	const (
		oldName = "rename_old"
		newName = "rename_new"
	)
	ctx := WillTest(t, "test")
	var exec ListExec
	exec.Append(CreateTable(oldName, "id text"))
	exec.Append(RenameTable(oldName, newName))
	exec.Append(DropTable(newName, true))
	Pgdb.SetLogOpt(LogAll)
	err := ExecTX(ctx, &exec)
	require.NoError(t, err)
}
