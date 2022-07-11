// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package geo

import (
	"strconv"
)

//-----------------------------------------------------------------------------

func trimLeft(s string) string {
	for i := 0; i < len(s); i++ {
		if s[i] != ' ' {
			return s[i:]
		}
	}
	return s
}

//-----------------------------------------------------------------------------

func skipLeft(s string, c byte) string {
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			s = s[i+1:]
			for len(s) > 0 && s[0] == ' ' {
				s = s[1:]
			}
			return s
		}
	}
	return ""
}

//-----------------------------------------------------------------------------

func readUntil(s string, c1, c2 byte) string {
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == c1 || c == c2 {
			return s[:i]
		}
	}
	return ""
}

//-----------------------------------------------------------------------------

func parseFloat(s string) (float64, bool) {
	f, err := strconv.ParseFloat(s, 64)
	return f, err == nil
}

//-----------------------------------------------------------------------------

func hasPrefix(s string, prefix string) bool {
	return len(s) >= len(prefix) && s[0:len(prefix)] == prefix
}