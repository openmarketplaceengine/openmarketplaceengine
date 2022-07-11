// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package cfg

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnvOpt(t *testing.T) {
	const opt EnvOpt = "envopt"
	require.Equal(t, "", opt.Get())
	require.False(t, opt.GetBool())
	opt.SetBool(true)
	require.Equal(t, "1", opt.Get())
	require.True(t, opt.GetBool())
}
