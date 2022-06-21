// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package dao

type Upgrade struct {
	Version int
	Details string
	Confirm string
}

const upgradeTable = "upgrade"

//-----------------------------------------------------------------------------

func createUpgradeTable(ctx Context) error {
	const upgradeSchema = SQLExec(
		`create table if not exists upgrade(
	    id     integer not null primary key,
	    info   text,
	    stamp  timestamptz,
	    status integer,
	    errmsg text
		)`)
	debugf("creating %q table", upgradeTable)
	return ExecDB(ctx, upgradeSchema)
}
