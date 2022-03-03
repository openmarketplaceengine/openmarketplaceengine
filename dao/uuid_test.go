// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package dao

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewXid(t *testing.T) {
	const count = 128
	xids := make([]XUID, count)
	for i := 0; i < count; i++ {
		xids[i] = NewXid()
	}
	head := xids[0]
	tail := xids[count-1]
	half := xids[count/2]
	sort.Strings(xids)
	require.Equal(t, head, xids[0])
	require.Equal(t, tail, xids[count-1])
	require.Equal(t, half, xids[count/2])
}
