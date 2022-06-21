// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package geo

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLocationWKB_UnmarshalText(t *testing.T) {
	var loc LocationWKB
	for i := range wkbPoints {
		p := &wkbPoints[i]
		err := loc.UnmarshalText([]byte(p.s))
		if p.err {
			require.Error(t, err)
			continue
		}
		require.NoError(t, err)
		require.True(t, loc.EqualsCoord(p.y, p.x))
	}
}

func TestLocationWKB_Scan(t *testing.T) {
	var loc LocationWKB
	require.NoError(t, loc.Scan(nil))
	require.NoError(t, loc.Scan(""))
	for i := range wkbPoints {
		p := &wkbPoints[i]
		err := loc.Scan(p.s)
		if p.err && len(p.s) > 0 {
			require.Error(t, err)
			continue
		}
		require.NoError(t, err)
		require.True(t, loc.EqualsCoord(p.y, p.x))
	}
	for i := range wktPoints {
		p := &wktPoints[i]
		err := loc.Scan([]byte(p.s))
		if p.err && len(p.s) > 0 {
			require.Error(t, err)
			continue
		}
		require.NoError(t, err)
		require.True(t, loc.EqualsCoord(p.y, p.x))
	}
}

func TestLocationWKB_EncodeWKB(t *testing.T) {
	var loc LocationWKB
	for i := range wkbPoints {
		p := &wkbPoints[i]
		if p.err || len(p.s) != defPointLenWKB {
			continue
		}
		require.NoError(t, loc.DecodeString(p.s))
		require.Equal(t, p.s, loc.EncodeWKB())
	}
}