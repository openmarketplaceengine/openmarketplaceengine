// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package geo

import (
	"encoding/binary"
	"fmt"
	"math"
	"strings"
	"unsafe"

	"github.com/openmarketplaceengine/openmarketplaceengine/app/enc/hex"
)

const (
	minPointLenWKB = 42
	defPointLenWKB = 50
	minPointLenWKT = 10 // POINT(1 1)
)

const (
	wktPointPrefix = "POINT"
	wktSridPrefix  = "SRID="
)

//-----------------------------------------------------------------------------

func DecodePoint(s string) (x, y float64, err error) {
	n := len(s)
	if n == defPointLenWKB || n == minPointLenWKB {
		if checkEndian(s) {
			return DecodePointWKB(s)
		}
	}
	if n >= minPointLenWKT {
		s = strings.ToUpper(s)
		if hasPrefix(s, wktPointPrefix) {
			return decodePointWKT(s)
		}
		if hasPrefix(s, wktSridPrefix) {
			v := s[len(wktSridPrefix):]
			if v = skipLeft(v, ';'); hasPrefix(v, wktPointPrefix) {
				return decodePointWKT(s)
			}
		}
	}
	return 0, 0, fmt.Errorf("invalid geo point format: %q", s)
}

//-----------------------------------------------------------------------------
// Well-Known Binary Format
//-----------------------------------------------------------------------------

func DecodePointWKB(s string) (x, y float64, err error) {
	const fename = "DecodePointWKB"
	if len(s) < minPointLenWKB {
		return 0, 0, funcError{fename, ErrSrcLen}
	}
	littleEndian, ok := parseEndian(s)
	if !ok {
		return 0, 0, funcError{fename, ErrEndian}
	}
	off := 2
	typ, terr := decodeType(s[off:], littleEndian, wkbPoint)
	if terr != nil {
		return 0, 0, funcError{fename, terr}
	}
	off += 8
	if (typ & wkbSRID) != 0 {
		off += 8 // skip SRID
		if (off + 32) < len(s) {
			return 0, 0, funcError{fename, ErrSrcLen}
		}
	}
	x, err = hex.DecodeFloat64(s[off:], littleEndian)
	if err != nil {
		return 0, 0, funcError{fename, err}
	}
	off += 16
	y, err = hex.DecodeFloat64(s[off:], littleEndian)
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
	hex.Encode(b, b[defPointLenWKB/2:], hex.UpperCase)
	return *(*string)(unsafe.Pointer(&b))
}

//-----------------------------------------------------------------------------
// Well-Known Text Format
//-----------------------------------------------------------------------------

type wktPointError string

func (e wktPointError) Error() string {
	return fmt.Sprintf("invalid POINT WKT format: %q", string(e))
}

//-----------------------------------------------------------------------------

func DecodePointWKT(s string) (x, y float64, err error) {
	const fename = "DecodePointWKT"
	if len(s) < minPointLenWKT {
		return 0, 0, funcError{fename, ErrSrcLen}
	}
	v := strings.ToUpper(s)
	if !hasPrefix(v, wktPointPrefix) {
		if !hasPrefix(v, wktSridPrefix) {
			return 0, 0, wktPointError(s)
		}
		v = v[len(wktSridPrefix):]
		if v = skipLeft(v, ';'); len(v) < minPointLenWKT || !hasPrefix(v, wktPointPrefix) {
			return 0, 0, wktPointError(s)
		}
	}
	return decodePointWKT(v)
}

func decodePointWKT(s string) (x, y float64, err error) {
	v := s[len(wktPointPrefix):]
	if v = skipLeft(v, '('); len(v) > 0 {
		sx := readUntil(v, ' ', ' ')
		var ok bool
		if x, ok = parseFloat(sx); ok {
			v = trimLeft(v[len(sx):])
			sy := readUntil(v, ')', ' ')
			if y, ok = parseFloat(sy); ok {
				return
			}
		}
	}
	return 0, 0, wktPointError(s)
}

//-----------------------------------------------------------------------------

func EncodePointWKT(x, y float64) string {
	return fmt.Sprintf("POINT(%f %f)", x, y)
}
