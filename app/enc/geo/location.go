// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package geo

import (
	"database/sql/driver"
	"fmt"
	"reflect"
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

//-----------------------------------------------------------------------------

func (loc *LocationWKB) DecodeString(s string) (err error) {
	if len(s) == 0 {
		loc.Reset()
		return
	}
	loc.Longitude, loc.Latitude, err = DecodePoint(s)
	return
}

//-----------------------------------------------------------------------------
// SQL Interface
//-----------------------------------------------------------------------------

func (loc *LocationWKB) Value() (driver.Value, error) {
	return EncodePointWKB(loc.Longitude, loc.Latitude), nil
}

func (loc *LocationWKB) Scan(src interface{}) error {
	switch v := src.(type) {
	case nil:
		loc.Reset()
		return nil
	case string:
		return loc.DecodeString(v)
	case []byte:
		return loc.UnmarshalText(v)
	default:
		return fmt.Errorf("invalid SQL scan point type: %q", reflect.TypeOf(src))
	}
}