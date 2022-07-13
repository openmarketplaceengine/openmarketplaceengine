// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package dao

import (
	"io/fs"
	"time"

	"github.com/openmarketplaceengine/openmarketplaceengine/app/dir"
)

type Upgrade struct {
	Version int
	Details string
	Confirm string
	Disable bool
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

func (u *Upgrade) Insert(ctx Context) error {
	return ExecTX(ctx, u.insert())
}

//-----------------------------------------------------------------------------

func (u *Upgrade) insert() Executable {
	sql := Insert(upgradeTable)
	sql.Set("id", u.Version)
	sql.Set("info", u.Details)
	sql.Set("stamp", time.Now())
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

//-----------------------------------------------------------------------------

func (u *upgradeManager) start(ctx Context) (err error) {
	if err = u.upgradeTableCreate(ctx); err != nil {
		return
	}
	if err = u.readFsys(); err != nil {
		return
	}
	n := u.list.Len()
	if n == 0 {
		debugf("upgrade scripts not registered")
		return
	}
	for i := 0; i < n && err == nil; i++ {
		err = u.readFile(u.list.Path(i))
	}
	return
}

//-----------------------------------------------------------------------------

func (u *upgradeManager) readFile(fpath *dir.FsysPath) (err error) {
	var upg Upgrade
	debugf("%s: processing upgrade", fpath.Name)
	if err = fpath.ReadYAML(&upg); err != nil {
		return
	}
	if upg.Disable {
		debugf("%s: upgrade is disabled", fpath.Name)
		return
	}
	return
}

//-----------------------------------------------------------------------------

func (u *upgradeManager) upgradeTableCreate(ctx Context) error {
	var exe ListExec
	// TODO: remove in the next PR
	exe.Append(DropColumn(upgradeTable, "status"))
	exe.Append(DropColumn(upgradeTable, "errmsg"))
	exe.Append(CreateTable(upgradeTable,
		"id     integer not null primary key",
		"info   text not null",
		"stamp  timestamptz not null",
	))
	return ExecDB(ctx, &exe)
}

//-----------------------------------------------------------------------------

func (u *upgradeManager) upgradeSelect(ctx Context, version int) (found bool, stamp time.Time, err error) {
	sql := From(upgradeTable)
	sql.Select("stamp").To(&stamp)
	sql.Where("id = ?", version)
	found, err = sql.QueryOne(ctx)
	return
}

//-----------------------------------------------------------------------------

func (u *upgradeManager) upgradeDelete(ctx Context, version int) error {
	sql := Delete(upgradeTable)
	sql.Where("id = ?", version)
	return ExecTX(ctx, sql)
}
