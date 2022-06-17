// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package geo

import (
	"encoding/binary"
	"encoding/hex"
	"math"
	"unsafe"
)

const (
	WKBPointLen = 50
)

func DecodePointWKB(src string) (x, y float64, err error) {
	const fename = "DecodePointWKB"
	if len(src) != WKBPointLen {
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
	x = math.Float64frombits(ufunc(val[9:]))
	y = math.Float64frombits(ufunc(val[17:]))
	return //nolint
}

//-----------------------------------------------------------------------------

func EncodePointWKB(x, y float64) string {
	var b = make([]byte, WKBPointLen)
	off := WKBPointLen / 2
	b[off] = 1
	off++
	binary.LittleEndian.PutUint32(b[off:], uint32(wkbPoint|wkbSRID))
	off += 4
	binary.LittleEndian.PutUint32(b[off:], DefSRID)
	off += 4
	binary.LittleEndian.PutUint64(b[off:], math.Float64bits(x))
	off += 8
	binary.LittleEndian.PutUint64(b[off:], math.Float64bits(y))
	hexenc(b, b[WKBPointLen/2:])
	return *(*string)(unsafe.Pointer(&b))
}
