// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package dao

import (
	"testing"

	"github.com/stretchr/testify/require"
)

//-----------------------------------------------------------------------------

func TestErrUndefinedTable(t *testing.T) {
	ctx := WillTest(t, "test")
	sql := From("undefined").Select("count(*)")
	has, err := sql.QueryOne(ctx)
	require.False(t, has)
	require.Error(t, err)
	require.True(t, ErrUndefinedTable.Is(err))
}
