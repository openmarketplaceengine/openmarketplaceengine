// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package dao

import (
	"context"
	"testing"

	"github.com/jackc/pgconn"
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

//-----------------------------------------------------------------------------

func TestShouldSkipError(t *testing.T) {
	ctx := context.Background()
	err := &pgconn.PgError{Code: string(ErrUniqueViolation)}
	require.False(t, ShouldSkipError(ctx, err))
	ctx = SkipErrorsContext(ctx, ErrUndefinedTable, ErrUniqueViolation)
	require.True(t, ShouldSkipError(ctx, err))
}