// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package geo

import (
	"fmt"
	"unsafe"

	"github.com/openmarketplaceengine/openmarketplaceengine/app/enc/hex"
)

type Type uint32
type SRID uint32

const (
	DefSRID = 4326
)

const (
	wkbPoint              = 1
	wkbLineString         = 2
	wkbPolygon            = 3
	wkbMultiPoint         = 4
	wkbMultiLineString    = 5
	wkbMultiPolygon       = 6
	wkbGeometryCollection = 7
)

const (
	wkbZ    Type = 0x80000000
	wkbM    Type = 0x40000000
	wkbSRID Type = 0x20000000
	wkbZM        = wkbZ | wkbM
	wkbMask      = wkbZ | wkbM | wkbSRID
)

//-----------------------------------------------------------------------------

func (t Type) String() string {
	v := t & ^wkbMask
	var s string
	switch v {
	case wkbPoint:
		s = "POINT"
	case wkbLineString:
		s = "LINESTRING"
	case wkbPolygon:
		s = "POLYGON"
	case wkbMultiPoint:
		s = "MULTIPOINT"
	case wkbMultiLineString:
		s = "MULTILINESTRING"
	case wkbMultiPolygon:
		s = "MULTIPOLYGON"
	case wkbGeometryCollection:
		s = "GEOMETRYCOLLECTION"
	default:
		return "<invalid>"
	}
	if t&wkbZM != 0 {
		b := make([]byte, 0, len(s)+3)
		b = append(b, s...)
		b = append(b, ' ')
		if t&wkbZ != 0 {
			b = append(b, 'Z')
		}
		if t&wkbM != 0 {
			b = append(b, 'M')
		}
		return *(*string)(unsafe.Pointer(&b))
	}
	return s
}

//-----------------------------------------------------------------------------

func checkType(must Type, have Type) bool {
	must &= ^wkbMask
	have &= ^wkbMask
	return (must == have)
}

//-----------------------------------------------------------------------------

func decodeType(s string, littleEndian bool, require Type) (Type, error) {
	u32, err := hex.DecodeUint32(s, littleEndian)
	if err != nil {
		return 0, err
	}
	typ := Type(u32)
	if require > 0 && !checkType(require, typ) {
		return typ, fmt.Errorf("invalid geo type: %q", typ)
	}
	return typ, nil
}
