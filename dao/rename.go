// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package dao

func RenameTable(oldName string, newName string) Executable {
	return SQLExecf("alter table if exists %s rename to %s", oldName, newName)
}