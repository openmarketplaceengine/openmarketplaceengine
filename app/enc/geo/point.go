// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package geo

import (
	"encoding/binary"
	"fmt"
	"math"
	"unsafe"
)

const (
	minPointLenWKB = 42
	defPointLenWKB = 50
)

func DecodePointWKB(s string) (x, y float64, err error) {
	const fename = "DecodePointWKB"
	if len(s) < minPointLenWKB {
		return 0, 0, funcError{fename, ErrSrcLen}
	}
	s1, s0 := s[1], s[0]
	littleEndian := (s1 == '1' && s0 == '0')
	if !littleEndian && s1 != '0' && s0 != '0' {
		return 0, 0, funcError{fename, ErrEndian}
	}
	off := 2
	var u32 uint32
	u32, err = hexdecU32(s[off:], littleEndian)
	if err != nil {
		return 0, 0, funcError{fename, err}
	}
	typ := Type(u32)
	off += 8
	if !checkType(wkbPoint, typ) {
		return 0, 0, funcError{fename, fmt.Errorf("invalid geo type: %q", typ)}
	}
	if (typ & wkbSRID) != 0 { // skip SRID
		off += 8
		if (off + 32) < len(s) {
			return 0, 0, funcError{fename, ErrSrcLen}
		}
	}
	x, err = hexdecF64(s[off:], littleEndian)
	if err != nil {
		return 0, 0, funcError{fename, err}
	}
	off += 16
	y, err = hexdecF64(s[off:], littleEndian)
	if err != nil {
		err = funcError{fename, err}
	}
	return //nolint
}

//-----------------------------------------------------------------------------

func EncodePointWKB(x, y float64) string {
	var b = make([]byte, defPointLenWKB)
	off := defPointLenWKB / 2
	b[off] = 1
	off++
	binary.LittleEndian.PutUint32(b[off:], uint32(wkbPoint|wkbSRID))
	off += 4
	binary.LittleEndian.PutUint32(b[off:], DefSRID)
	off += 4
	binary.LittleEndian.PutUint64(b[off:], math.Float64bits(x))
	off += 8
	binary.LittleEndian.PutUint64(b[off:], math.Float64bits(y))
	hexenc(b, b[defPointLenWKB/2:])
	return *(*string)(unsafe.Pointer(&b))
}
