// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package dao

import (
	"errors"
	"fmt"
	"io/fs"
	"strings"
	"time"

	"github.com/openmarketplaceengine/openmarketplaceengine/app/dir"
)

type Upgrade struct {
	exec    ListExec
	stamped time.Time
	Version int
	Details string
	Drop    Drop
	Tables  []*UpgradeTable
	Disable bool
	Logonly bool
	Mustrun bool
}

type UpgradeTable struct {
	Name   string
	Create string
	Select string
	exec   ListExec
	count  int
}

type vermap = map[int]time.Time

type upgradeManager struct {
	upfs []dir.FS
	list dir.FsysList
	vmap vermap // version map
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

func (u *Upgrade) Stamp() time.Time {
	if u.stamped.IsZero() {
		return time.Now()
	}
	return u.stamped
}

//-----------------------------------------------------------------------------

func (u *Upgrade) insert() Executable {
	sql := Insert(upgradeTable)
	sql.Set("id", u.Version)
	sql.Set("info", u.Details)
	sql.Set("stamp", u.Stamp())
	return sql
}

//-----------------------------------------------------------------------------

func (u *Upgrade) Upsert(ctx Context) error {
	return ExecTX(ctx, u.upsert())
}

//-----------------------------------------------------------------------------

func (u *Upgrade) upsert() *AtomicUpsert {
	sql := UpsertAtomic(upgradeTable)
	sql.Key("id", u.Version)
	sql.Set("info", u.Details)
	sql.Set("stamp", u.Stamp())
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
	u.vmap = nil
	u.list.Clear()
}

//-----------------------------------------------------------------------------

func (u *upgradeManager) readFsys() (err error) {
	if len(u.upfs) == 0 {
		return
	}
	u.list.Alloc(8)
	for i := 0; i < len(u.upfs) && err == nil; i++ {
		err = u.list.ListFext(u.upfs[i], ".yaml")
	}
	return
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
	if u.vmap, err = u.loadVmap(ctx); err != nil {
		return
	}
	for i := 0; i < n && err == nil; i++ {
		err = u.readFile(ctx, u.list.Path(i))
	}
	return
}

//-----------------------------------------------------------------------------

func (u *upgradeManager) readFile(ctx Context, fpath *dir.FsysPath) (err error) {
	var upg Upgrade
	debugf("processing upgrade: %q", fpath.Name)
	if err = fpath.ReadYAML(&upg); err != nil {
		return
	}
	if upg.Disable {
		debugf("upgrade is disabled: %q", fpath.Name)
		return
	}
	if stamp, ok := u.vmap[upg.Version]; ok {
		infof("found schema version %d (upgraded on %s)", upg.Version, stamp.Format(time.RFC822))
		if !upg.Mustrun && !upg.Logonly {
			return
		}
	}
	return upg.upgrade(ctx)
}

//-----------------------------------------------------------------------------

func (u *Upgrade) upgrade(ctx Context) error {
	if u.Mustrun && !u.stamped.IsZero() {
		infof("force upgrade to version %d: %q", u.Version, u.Details)
	} else {
		infof("upgrading to version %d: %q", u.Version, u.Details)
	}
	if u.Logonly {
		debugf("log only upgrade mode")
	}
	if u.Drop.Len() > 0 {
		u.exec.Append(&u.Drop)
	}
	for i := range u.Tables {
		if tbl := u.Tables[i]; tbl != nil {
			if err := tbl.upgrade(ctx, u.Version); err != nil {
				return err
			}
			if tbl.exec.Len() > 0 {
				u.exec.Append(&tbl.exec)
			}
		}
	}
	u.exec.Append(u.upsert())
	if u.Logonly {
		return ExecFlog(ctx, infof, &u.exec)
	}
	return ExecTX(ctx, &u.exec)
}

//-----------------------------------------------------------------------------

func (t *UpgradeTable) validate() error {
	t.Name = strings.TrimSpace(t.Name)
	if len(t.Name) == 0 {
		return errors.New("upgrade table name is empty")
	}
	t.Select = strings.TrimSpace(t.Select)
	if len(t.Select) == 0 {
		return fmt.Errorf("empty 'select' value for table: %q", t.Name)
	}
	t.Create = strings.TrimSpace(t.Create)
	if len(t.Create) == 0 {
		return fmt.Errorf("empty 'create' value for table: %q", t.Name)
	}
	return nil
}

//-----------------------------------------------------------------------------

func (t *UpgradeTable) upgrade(ctx Context, ver int) error {
	if err := t.validate(); err != nil {
		return err
	}
	// count
	sql := From(t.Name).Count("*").To(&t.count)
	_, err := sql.QueryOne(SkipUndefErrors(ctx))
	if err != nil {
		if ErrUndefinedTable.Is(err) {
			debugf("skipping upgrade for not found table: %q", t.Name)
			return nil
		}
		return err
	}
	infof("upgrading table: %q: to version: %d", t.Name, ver)
	// verify columns
	sql = From(t.Name).Select(t.Select).Limit(1)
	_, err = sql.QueryOne(SkipUndefErrors(ctx))
	if err != nil {
		if ErrUndefinedColumn.Is(err) {
			// check if table already upgraded
			return nil
		}
		return fmt.Errorf("table %q select error: %w", t.Name, err)
	}
	return nil
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

//-----------------------------------------------------------------------------

func (u *upgradeManager) loadVmap(ctx Context) (vermap, error) {
	vmap := make(vermap)
	sql := From(upgradeTable).Select("id").Select("stamp")
	err := sql.QueryRows(ctx, func(rows *Rows) error {
		var ver int
		var stamp time.Time
		for rows.Next() {
			if serr := rows.Scan(&ver, &stamp); serr != nil {
				return serr
			}
			vmap[ver] = stamp
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return vmap, nil
}
