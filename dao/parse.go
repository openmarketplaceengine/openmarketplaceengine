// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package dao

import (
	"fmt"
	"strings"
)

func splitCreateTableColumns(s string) ([]string, error) {
	s = strings.ToLower(strings.TrimSpace(s))
	if strings.HasPrefix(s, "create table") {
		x := strings.IndexByte(s, '(')
		if x == -1 {
			return nil, fmt.Errorf("opening bracket not found")
		}
		s = s[x+1:]
	}
	n := len(s)
	if n > 0 && s[0] == '(' {
		s = s[1:]
		if n--; n == 0 {
			return nil, fmt.Errorf("dangling opening bracket")
		}
	}
	if n == 0 {
		return nil, nil
	}
	cols := make([]string, 0, 8)
	nest := 0
	// start := -1
	// space := -1
loop:
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch c {
		case ',':
			if nest > 0 {
				continue loop
			}
		case '(':
			nest++
		case ')':
			if nest == 0 {
				break loop
			}
			nest--
		}
	}
	if nest > 0 {
		return nil, fmt.Errorf("brackets do not match")
	}
	return cols, nil
}
