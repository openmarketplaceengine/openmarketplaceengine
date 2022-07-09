// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package dao

import "strings"

func CreateTable(name string, cols ...string) Executable {
	var sb strings.Builder
	sb.Grow(128)
	sb.WriteString("create table if not exists ")
	sb.WriteString(name)
	sb.WriteString(" (\n")
	for i := 0; i < len(cols); i++ {
		if i > 0 {
			sb.WriteString(",\n")
		}
		sb.WriteString("  ")
		sb.WriteString(cols[i])
	}
	sb.WriteString("\n)")
	return SQLExec(sb.String())
}
