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
}

//-----------------------------------------------------------------------------

func TestDecodePointWKB(t *testing.T) {
	for i := range wkbPoints {
		v := &wkbPoints[i]
		lat, lon, err := DecodePointWKB(v.src)
		require.NoError(t, err)
		require.Equal(t, v.lat, lat)
		require.Equal(t, v.lon, lon)
	}
}
