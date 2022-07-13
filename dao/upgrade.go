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
	Disable bool
	Details string
	Confirm string
	fpath   *dir.FsysPath //nolint
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

func (u *upgradeManager) start(ctx Context) error {
	err := u.upgradeTableCreate(ctx)
	return err
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
