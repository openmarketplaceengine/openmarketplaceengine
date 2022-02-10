// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package log

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDevelConfig(t *testing.T) {
	callInit(t, DevelConfig())
	sampleEntries()
}

func TestJsonConfig(t *testing.T) {
	callInit(t, DevelConfig().WithStyle("json"))
	sampleEntries()
}

//-----------------------------------------------------------------------------
// Helpers
//-----------------------------------------------------------------------------

func sampleEntries() {
	Debugf("sample log text with: %q", "debug")
	Infof("sample log text with: %q", "info")
	Warnf("sample log text with: %q", "warn")
	Errorf("sample log text with: %q", "error")
}

//-----------------------------------------------------------------------------

func callInit(t testing.TB, c ConfigHolder) {
	require.NoError(t, Init(c))
	t.Cleanup(SafeSync)
}
