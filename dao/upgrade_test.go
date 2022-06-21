// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package dao

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

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
