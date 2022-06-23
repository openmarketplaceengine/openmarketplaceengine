// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package dao

import (
	"context"
	"embed"
	"testing"

	"github.com/stretchr/testify/require"
)

//go:embed testdata/fsys/upgrade/*.upgrade.yaml
var testUpfs embed.FS

//go:embed testdata/fsys/upgrade/dummy.yaml
var testUpfsDummy embed.FS

func TestUpgradeCRUD(t *testing.T) {
	const ver = -1

	WillTest(t, "test")

	ctx := context.Background()

	require.NoError(t, upgradeDelete(ctx, ver))

	has, err := upgradeSelect(ctx, ver)
	require.NoError(t, err)
	require.False(t, has)

	var upg Upgrade

	upg.Version = -1
	upg.Details = "Upgrade CRUD testing"
	upg.success = true
	require.NoError(t, upg.Insert(ctx))

	has, err = upgradeSelect(ctx, ver)
	require.NoError(t, err)
	require.True(t, has)

	require.NoError(t, upgradeDelete(ctx, ver))

	upg.success = false
	upg.errtext = "test error message"
	require.NoError(t, upg.Insert(ctx))

	has, err = upgradeSelect(ctx, ver)
	require.Error(t, err)
	require.False(t, has)

	require.NoError(t, upgradeDelete(ctx, ver))
}

//-----------------------------------------------------------------------------

func TestRegisterUpgrade(t *testing.T) {
	RegisterUpgrade(testUpfs)
	RegisterUpgrade(testUpfsDummy)
	upgr := &Pgdb.upgr
	err := upgr.readFsys()
	require.NoError(t, err)
	list := &upgr.list
	if list.Len() < 3 {
		t.Fatalf("Upgrade manager must have read minimum 3 scripts, but have %d", list.Len())
	}
	require.Equal(t, "-01.upgrade.yaml", list.Path(0).Name)
	require.Equal(t, "-02.upgrade.yaml", list.Path(1).Name)
	require.Equal(t, "dummy.yaml", list.Path(2).Name)
}