// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package job

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStateFromString(t *testing.T) {
	for i := range stateString {
		_, found := StateFromString(stateString[i])
		require.True(t, found)
	}
	_, found := StateFromString("dummy")
	require.False(t, found)
}

func TestState_String(t *testing.T) {
	for k, v := range stringState {
		require.Equal(t, k, v.String())
	}
}