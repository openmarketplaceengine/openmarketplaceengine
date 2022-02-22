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
	sampleEntries(Log())
}

func TestDevelLevel(t *testing.T) {
	callInit(t, DevelConfig().WithTrace(false).WithLevel("info"))
	sampleEntries(Log())
}

func TestJsonConfig(t *testing.T) {
	callInit(t, DevelConfig().WithStyle("json"))
	sampleEntries(Log())
}

func TestNamedLogger(t *testing.T) {
	callInit(t, DevelConfig().WithTrace(false))
	sampleEntries(Log().Named("TEST"))
}

func TestNamedLevel(t *testing.T) {
	callInit(t, DevelConfig().WithTrace(false))
	sampleEntries(Log().NamedLevel("TEST", LevelWarn))
}

func TestLevelSetter(t *testing.T) {
	callInit(t, DevelConfig())
	if _, ok := zlog.c.(levelSetter); ok {
		t.Logf("zlog.Core is a levelSetter")
	}
	if _, ok := zlog.c.(*levelCore); ok {
		t.Logf("zlog.Core is a levelCore")
	}
}

//-----------------------------------------------------------------------------
// Helpers
//-----------------------------------------------------------------------------

func sampleEntries(log Logger) {
	if !testing.Verbose() {
		return
	}
	log.Debugf("sample log text with: %q", "debug")
	log.Infof("sample log text with: %q", "info")
	log.Warnf("sample log text with: %q", "warn")
	log.Errorf("sample log text with: %q", "error")
}

//-----------------------------------------------------------------------------

func callInit(t testing.TB, c ConfigHolder) {
	require.NoError(t, Init(c))
	t.Cleanup(SafeSync)
}
