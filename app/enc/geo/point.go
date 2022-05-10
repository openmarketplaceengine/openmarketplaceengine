// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package geo

import (
	"encoding/binary"
	"encoding/hex"
	"math"
)

func DecodePointWKB(src string) (lat, lon float64, err error) {
	const fename = "DecodePointWKB"
	const srclen = 50
	if len(src) != srclen {
		return 0, 0, funcError{fename, ErrSrcLen}
	}
	var val []byte
	val, err = hex.DecodeString(src)
	if err != nil {
		return 0, 0, funcError{fename, ErrHexDec}
	}
	var ufunc func([]byte) uint64
	switch val[0] {
	case 0:
		ufunc = binary.BigEndian.Uint64
	case 1:
		ufunc = binary.LittleEndian.Uint64
	default:
		return 0, 0, funcError{fename, ErrEndian}
	}
	lon = math.Float64frombits(ufunc(val[9:]))
	lat = math.Float64frombits(ufunc(val[17:]))
	return //nolint
}
