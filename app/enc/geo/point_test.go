// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package geo

import (
	"testing"

	"github.com/stretchr/testify/require"
)

//-----------------------------------------------------------------------------

var wkbPoints = []struct {
	src string
	lat float64
	lon float64
}{
	{"0101000020E6100000FE0C6FD6E07A52C02A00C633685C4440", 40.72193, -73.919973},
	{"0101000020E6100000A1754309697B52C073E26190AA544440", 40.6614552, -73.9282859},
	{"0101000020E6100000000000000000F03F000000000000F03F", 1.0, 1.0},
	{"0101000020E6100000000000000000F0BF000000000000F03F", 1.0, -1.0},
	{"0101000000000000000000F03F000000000000F0BF", -1.0, 1.0},
}

//-----------------------------------------------------------------------------

var wktPoints = []struct {
	s   string
	x   float64
	y   float64
	err bool
}{
	{"", 0, 0, true},
	{"POINT", 0, 0, true},
	{"SRID=4326;", 0, 0, true},
	{"POINT(1.1 2.2)", 1.1, 2.2, false},
	{"POINT  (  1.1   2.2  )", 1.1, 2.2, false},
	{"SRID=4326;POINT(1 -1)", 1, -1, false},
	{"SRID=4326;  POINT(1 -1)", 1, -1, false},
	{"POINT Z (  1     2   3)", 1, 2, false},
}

//-----------------------------------------------------------------------------
// WKB
//-----------------------------------------------------------------------------

func TestDecodePointWKB(t *testing.T) {
	for i := range wkbPoints {
		v := &wkbPoints[i]
		lon, lat, err := DecodePointWKB(v.src)
		require.NoError(t, err)
		require.Equal(t, v.lat, lat)
		require.Equal(t, v.lon, lon)
	}
}

//-----------------------------------------------------------------------------

func TestEncodePointWKB(t *testing.T) {
	for i := range wkbPoints {
		v := &wkbPoints[i]
		if len(v.src) == defPointLenWKB {
			s := EncodePointWKB(v.lon, v.lat)
			require.Equal(t, v.src, s)
		}
	}
}

//-----------------------------------------------------------------------------
// WKT
//-----------------------------------------------------------------------------

func TestDecodePointWKT(t *testing.T) {
	for i := range wktPoints {
		p := &wktPoints[i]
		x, y, err := DecodePointWKT(p.s)
		if p.err {
			require.Error(t, err)
			continue
		}
		require.NoError(t, err)
		require.Equal(t, p.x, x)
		require.Equal(t, p.y, y)
	}
}
