// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package dao

import (
	"database/sql/driver"
	"fmt"
	"time"
	"unsafe"
)

type Time struct {
	time.Time
}

var zeroTime time.Time

var timezFormats = []string{
	"2006-01-02 15:04:05.999999999Z07:00:00",
	"2006-01-02 15:04:05.999999999Z07:00",
	"2006-01-02 15:04:05.999999999Z07",
}

var timeuFormats = []string{
	"2006-01-02 15:04:05.999999999",
	"2006-01-02",
}

func (t *Time) Now() {
	t.Time = time.Now()
}

func (t *Time) Reset() {
	t.Time = zeroTime
}

func (t Time) IsNull() bool {
	return t.IsZero() && t.Location() == time.UTC
}

func (t *Time) Scan(src interface{}) error {
	if src == nil {
		t.Time = zeroTime
		return nil
	}
	switch val := src.(type) {
	case time.Time:
		t.Time = val
	case string:
		return t.parse(val)
	case []byte:
		s := *(*string)(unsafe.Pointer(&val))
		return t.parse(s)
	}
	return nil
}

func (t *Time) parse(s string) error {
	for _, layout := range timezFormats {
		if val, err := time.Parse(layout, s); err == nil {
			t.Time = val
			return nil
		}
	}
	for _, layout := range timeuFormats {
		if val, err := time.ParseInLocation(layout, s, time.UTC); err == nil {
			t.Time = val
			return nil
		}
	}
	return fmt.Errorf("invalid Time value: %q", s)
}

func (t Time) Value() (driver.Value, error) {
	if t.IsNull() {
		return nil, nil
	}
	return t.Time, nil
}
