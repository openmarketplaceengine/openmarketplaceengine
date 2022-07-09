// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package dao

import "fmt"

type Drop struct {
	Index   []string
	Table   []string
	View    []string
	Cascade bool
}

//-----------------------------------------------------------------------------

func (d *Drop) Execute(ctx Context, exe Executor) (err error) {
	drop := d.makeDrop()
	for i := 0; i < len(drop) && err == nil; i++ {
		_, err = exe.ExecContext(ctx, drop[i])
	}
	return
}

func (d *Drop) Result() Result {
	return nil
}

//-----------------------------------------------------------------------------

func (d *Drop) AppendIndex(index ...string) *Drop {
	d.Index = append(d.Index, index...)
	return d
}

func (d *Drop) AppendTable(table ...string) *Drop {
	d.Table = append(d.Table, table...)
	return d
}

func (d *Drop) AppendView(view ...string) *Drop {
	d.View = append(d.View, view...)
	return d
}

//-----------------------------------------------------------------------------

func (d *Drop) Len() int {
	return len(d.Index) + len(d.Table) + len(d.View)
}

func (d *Drop) makeDrop() []string {
	if n := d.Len(); n > 0 {
		drop := make([]string, 0, n)
		drop = d.appendDrop("view", d.View, drop)
		drop = d.appendDrop("index", d.Index, drop)
		drop = d.appendDrop("table", d.Table, drop)
		return drop
	}
	return nil
}

func (d *Drop) appendDrop(kind string, name []string, drop []string) []string {
	for i := 0; i < len(name); i++ {
		drop = append(drop, MakeDrop(kind, name[i], d.Cascade))
	}
	return drop
}

//-----------------------------------------------------------------------------

func MakeDrop(kind string, name string, cascade bool) string {
	if cascade {
		return fmt.Sprintf("drop %s if exists %q cascade", kind, name)
	}
	return fmt.Sprintf("drop %s if exists %q", kind, name)
}

//-----------------------------------------------------------------------------

func DropTable(name string, cascade bool) Executable {
	return SQLExec(MakeDrop("table", name, cascade))
}

//-----------------------------------------------------------------------------

func DropIndex(name string, cascade bool) Executable {
	return SQLExec(MakeDrop("index", name, cascade))
}

//-----------------------------------------------------------------------------

func DropView(name string, cascade bool) Executable {
	return SQLExec(MakeDrop("view", name, cascade))
}
