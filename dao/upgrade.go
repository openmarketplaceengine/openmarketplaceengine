// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package dao

import (
	"fmt"
	"io/fs"
	"time"

	"github.com/openmarketplaceengine/openmarketplaceengine/app/dir"
)

type Upgrade struct {
	Version int
	Disable bool
	Details string
	Confirm string
	fpath   *dir.FsysPath //nolint
	errtext string
	success bool
}

type upgradeManager struct {
	upfs []dir.FS
	list dir.FsysList
}

const upgradeTable = "upgrade"

//-----------------------------------------------------------------------------

func RegisterUpgrade(fs fs.FS) {
	Pgdb.upgr.registerUpgrade(fs)
}

//-----------------------------------------------------------------------------

func upgradeTableCreate(ctx Context) error {
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

//-----------------------------------------------------------------------------

func upgradeSelect(ctx Context, version int) (found bool, err error) {
	var status int
	var errmsg string
	sql := From(upgradeTable)
	sql.Select("status").To(&status)
	sql.Select("COALESCE(errmsg, '')").To(&errmsg)
	sql.Where("id = ?", version)
	found, err = sql.QueryOne(ctx)
	if err != nil || !found || status != 0 {
		return
	}
	return false, fmt.Errorf("Upgrade from version %d failed: %s", version, errmsg)
}

//-----------------------------------------------------------------------------

func upgradeDelete(ctx Context, version int) error {
	sql := Delete(upgradeTable)
	sql.Where("id = ?", version)
	return ExecTX(ctx, sql)
}

//-----------------------------------------------------------------------------

func (u *Upgrade) Insert(ctx Context) error {
	return ExecTX(ctx, u.insert())
}

//-----------------------------------------------------------------------------

func (u *Upgrade) insert() Executable {
	sql := Insert(upgradeTable)
	sql.Set("id", u.Version)
	sql.Set("info", u.Details)
	sql.Set("stamp", time.Now())
	if u.success {
		sql.Set("status", 1)
	} else {
		sql.Set("status", 0)
		sql.Set("errmsg", u.errtext)
	}
	return sql
}

//-----------------------------------------------------------------------------
// Upgrade Manager
//-----------------------------------------------------------------------------

func (u *upgradeManager) registerUpgrade(fs fs.FS) {
	if fs != nil {
		u.upfs = append(u.upfs, fs)
	}
}

//-----------------------------------------------------------------------------

func (u *upgradeManager) clear() {
	u.upfs = nil
}

//-----------------------------------------------------------------------------

func (u *upgradeManager) readFsys() error {
	if len(u.upfs) == 0 {
		return nil
	}
	u.list.Alloc(8)
	var err error
	for i := 0; i < len(u.upfs) && err == nil; i++ {
		err = u.list.ListFext(u.upfs[i], ".yaml")
	}
	return err
}
