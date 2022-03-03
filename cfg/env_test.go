// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package cfg

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnvSkipTrim(t *testing.T) {
	m := map[string]bool{
		"":      true,
		" ":     false,
		"x=y":   true,
		"x=y\n": false,
		"\tx=y": false,
	}
	for k, v := range m {
		require.Equal(t, v, skipTrim(k))
	}
}

func TestEnvFastTrim(t *testing.T) {
	m := map[string]string{
		"":                "",
		" ":               "",
		"x=y":             "x=y",
		"x=y\n":           "x=y",
		"x=y\n\r\t\f":     "x=y",
		"\t\n\rx=y":       "x=y",
		"\t\n\rx=y\v\b\r": "x=y",
	}
	for k, v := range m {
		require.Equal(t, v, fastTrim(k))
	}
}

func TestEnvKeyVal(t *testing.T) {
	ary := [][3]string{
		{"", "", ""},
		{"=", "", ""},
		{"x=", "x", ""},
		{"x=y", "x", "y"},
		{"=y", "", "y"},
	}
	for i := range ary {
		kva := &ary[i]
		key, val := envKeyVal(kva[0])
		require.Equal(t, kva[1], key)
		require.Equal(t, kva[2], val)
	}
}

//-----------------------------------------------------------------------------

func TestTrimEnv(t *testing.T) {
	const key = "\n\t\rTRIM"
	require.NoError(t, os.Setenv(key, "TRUE\v\n\r   "))
	trimEnv()
	require.NoError(t, os.Unsetenv(key))
	env := os.Environ()
	for i := range env {
		if skipTrim(env[i]) {
			continue
		}
		t.Fatalf("env var needs trimming: %q", env[i])
	}
	require.Equal(t, "TRUE", os.Getenv("TRIM"))
}
