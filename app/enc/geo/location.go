// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package geo

import (
	"fmt"
	"unsafe"
)

type Location struct {
	Latitude  float64
	Longitude float64
}

type LocationWKB struct {
	Latitude  float64
	Longitude float64
}

//-----------------------------------------------------------------------------

func (loc *Location) Reset() {
	*loc = Location{}
}

//-----------------------------------------------------------------------------

func (loc *LocationWKB) Reset() {
	*loc = LocationWKB{}
}

//-----------------------------------------------------------------------------

func (loc Location) String() string {
	return fmt.Sprintf("[%f, %f]", loc.Latitude, loc.Longitude)
}

//-----------------------------------------------------------------------------

func (loc Location) EqualsCoord(lat, lon float64) bool {
	return loc.Latitude == lat && loc.Longitude == lon
}

//-----------------------------------------------------------------------------

func (loc LocationWKB) EqualsCoord(lat, lon float64) bool {
	return loc.Latitude == lat && loc.Longitude == lon
}

//-----------------------------------------------------------------------------

func (loc *LocationWKB) UnmarshalText(text []byte) (err error) {
	if len(text) == 0 {
		loc.Reset()
		return
	}
	src := *(*string)(unsafe.Pointer(&text))
	loc.Longitude, loc.Latitude, err = DecodePoint(src)
	return
}
