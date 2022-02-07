// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package cfg

// Password is a string type with protection from plain text serialization.
type Password string

const mockPass = "****************"

var (
	zeroPass = []byte("")
	safePass = []byte(mockPass)
)

// MarshalText function protects exposing Password as a plain text
// in JSON/YAML/Text serialization.
func (p Password) MarshalText() ([]byte, error) {
	if len(p) == 0 {
		return zeroPass, nil
	}
	return safePass, nil
}

// String mocks Password string from printing/formatting.
func (p Password) String() string {
	if len(p) > 0 {
		return mockPass
	}
	return ""
}
