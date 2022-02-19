// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package cfg

import (
	"strings"
	"unsafe"
)

// Address represents a connection URL
// with safe password printing/serialization.
type Address string

// MarshalText function protects exposing address password as a plain text
// in JSON/YAML/Text serialization.
func (a Address) MarshalText() ([]byte, error) {
	if len(a) == 0 {
		return zeroPass, nil
	}
	if off, end := a.findPass(); end > 0 {
		return a.hidePass(off, end), nil
	}
	return []byte(a), nil
}

// String hides address password from printing/formatting.
func (a Address) String() string {
	if len(a) == 0 {
		return ""
	}
	if off, end := a.findPass(); end > 0 {
		b := a.hidePass(off, end)
		return *(*string)(unsafe.Pointer(&b))
	}
	return string(a)
}

//-----------------------------------------------------------------------------

func (a Address) hidePass(off, end int) []byte {
	const hide = "********"
	b := make([]byte, 0, len(a)+len(hide)-(end-off))
	b = append(b, a[:off]...)
	b = append(b, hide...)
	b = append(b, a[end:]...)
	return b
}

//-----------------------------------------------------------------------------

func (a Address) findPass() (off, end int) {
	s := string(a)
	if end = strings.IndexByte(s, '@'); end > 0 {
		s = s[:end+1]
		if off = strings.LastIndexByte(s, ':'); off != -1 {
			if off++; s[off] != '/' {
				return
			}
			off = 0
		}
		end = 0
	}
	return
}
