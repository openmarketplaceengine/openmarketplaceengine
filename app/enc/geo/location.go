// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package geo

import (
	"fmt"
	"unsafe"
)

type Location struct {
	Lat float64
	Lon float64
}

type LocationWKB struct {
	Location
}

//-----------------------------------------------------------------------------

func (loc *Location) Reset() {
	*loc = Location{}
}

//-----------------------------------------------------------------------------

func (loc Location) String() string {
	return fmt.Sprintf("[%f, %f]", loc.Lat, loc.Lon)
}

//-----------------------------------------------------------------------------

func (loc Location) EqualsCoord(lat, lon float64) bool {
	return loc.Lat == lat && loc.Lon == lon
}

//-----------------------------------------------------------------------------

func (loc *LocationWKB) UnmarshalText(text []byte) (err error) {
	switch len(text) {
	case 0:
		loc.Reset()
	case WKBPointLen:
		src := *(*string)(unsafe.Pointer(&text))
		loc.Lat, loc.Lon, err = DecodePointWKB(src)
	default:
		return ErrSrcLen
	}
	return
}
