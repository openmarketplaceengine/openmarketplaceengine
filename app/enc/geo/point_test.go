// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package geo

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type pointTest struct {
	s   string
	x   float64
	y   float64
	err bool
}

//-----------------------------------------------------------------------------

var wkbPoints = []pointTest{
	{"0101000020E6100000FE0C6FD6E07A52C02A00C633685C4440", -73.919973, 40.72193, false},
	{"0101000020E6100000A1754309697B52C073E26190AA544440", -73.9282859, 40.6614552, false},
	{"0101000020E6100000000000000000F03F000000000000F03F", 1.0, 1.0, false},
	{"0101000020E6100000000000000000F0BF000000000000F03F", -1.0, 1.0, false},
	{"0101000000000000000000F03F000000000000F0BF", 1.0, -1.0, false},
}

//-----------------------------------------------------------------------------

var wktPoints = []pointTest{
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
	testDecodePoint(t, wkbPoints, DecodePointWKB)
}

//-----------------------------------------------------------------------------

func TestEncodePointWKB(t *testing.T) {
	for i := range wkbPoints {
		p := &wkbPoints[i]
		if !p.err && len(p.s) == defPointLenWKB {
			s := EncodePointWKB(p.x, p.y)
			require.Equal(t, p.s, s)
		}
	}
}

//-----------------------------------------------------------------------------
// WKT
//-----------------------------------------------------------------------------

func TestDecodePointWKT(t *testing.T) {
	testDecodePoint(t, wktPoints, DecodePointWKT)
}

//-----------------------------------------------------------------------------
// Any
//-----------------------------------------------------------------------------

func TestDecodePoint(t *testing.T) {
	testDecodePoint(t, wkbPoints, DecodePoint)
	testDecodePoint(t, wktPoints, DecodePoint)
}

//-----------------------------------------------------------------------------

func testDecodePoint(t *testing.T, ary []pointTest, dec func(string) (float64, float64, error)) {
	for i := range ary {
		p := &ary[i]
		x, y, err := dec(p.s)
		if p.err {
			require.Error(t, err)
			continue
		}
		require.NoError(t, err)
		require.Equal(t, p.x, x)
		require.Equal(t, p.y, y)
	}
}
