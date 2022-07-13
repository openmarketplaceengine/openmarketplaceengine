// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package dir

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"path/filepath"
	"strconv"
	"strings"
	"unsafe"

	"gopkg.in/yaml.v2"
)

type (
	FS    = fs.FS
	Entry = fs.DirEntry
)

type FsysPath struct {
	Fsys fs.FS
	Path string
	Name string
	Mode FileMode
}

type FsysList struct {
	list []*FsysPath
}

//goland:noinspection GoUnusedGlobalVariable
var SkipDir = fs.SkipDir

//-----------------------------------------------------------------------------

func WalkDir(fsys FS, root string, fn func(path string, d Entry) error) error {
	err := fs.WalkDir(fsys, root, func(path string, d Entry, err error) error {
		if err == nil {
			err = fn(path, d)
		}
		return err
	})
	return err
}

//-----------------------------------------------------------------------------
// FsysPath
//-----------------------------------------------------------------------------

func (f *FsysPath) Reset(fsys FS, path string, entry Entry) *FsysPath {
	f.Fsys = fsys
	f.Path = path
	if entry != nil {
		f.Name = entry.Name()
		f.Mode = entry.Type()
	} else {
		f.Name = filepath.Base(path)
		f.Mode = 0
	}
	return f
}

//-----------------------------------------------------------------------------

func (f *FsysPath) ReadFile() ([]byte, error) {
	return fs.ReadFile(f.Fsys, f.Path)
}

//-----------------------------------------------------------------------------

func (f *FsysPath) ReadText() (string, error) {
	buf, err := fs.ReadFile(f.Fsys, f.Path)
	if err != nil {
		return "", err
	}
	return *(*string)(unsafe.Pointer(&buf)), nil
}

//-----------------------------------------------------------------------------

func (f *FsysPath) ReadYAML(dst interface{}) error {
	return f.Decode(dst, yaml.Unmarshal)
}

//-----------------------------------------------------------------------------

func (f *FsysPath) ReadJSON(dst interface{}) error {
	return f.Decode(dst, json.Unmarshal)
}

//-----------------------------------------------------------------------------

func (f *FsysPath) Decode(dst interface{}, dec func(buf []byte, dst interface{}) error) error {
	buf, err := f.ReadFile()
	if err == nil {
		err = dec(buf, dst)
	}
	return err
}

//-----------------------------------------------------------------------------

func (f *FsysPath) IntPrefix() (int, bool) {
	i := 0
	s := f.Name
	n := len(s)
	if n > 0 && s[0] == '-' {
		i++
	}
loop:
	for i < n {
		switch s[i] {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			i++
			continue loop
		default:
			break loop
		}
	}
	if i == 0 || (i == 1 && s[0] == '-') {
		return 0, false
	}
	v, err := strconv.Atoi(s[:i])
	return v, err == nil
}

//-----------------------------------------------------------------------------

func (f *FsysPath) String() string {
	return f.Path
}

//-----------------------------------------------------------------------------
// FsysList
//-----------------------------------------------------------------------------

func (f *FsysList) Alloc(size int) {
	f.list = make([]*FsysPath, 0, size)
}

//-----------------------------------------------------------------------------

func (f *FsysList) Len() int {
	return len(f.list)
}

//-----------------------------------------------------------------------------

func (f *FsysList) Path(i int) *FsysPath {
	return f.list[i]
}

//-----------------------------------------------------------------------------

func (f *FsysList) Append(fsys FS, path string, name string, mode FileMode) {
	f.list = append(f.list, &FsysPath{fsys, path, name, mode})
}

//-----------------------------------------------------------------------------

func (f *FsysList) ListFsys(fsys fs.FS, filter func(path string, d Entry) (bool, error)) error {
	var ok bool
	err := fs.WalkDir(fsys, ".", func(path string, d Entry, err error) error {
		if err == nil {
			if ok, err = filter(path, d); ok {
				f.Append(fsys, path, d.Name(), d.Type())
			}
		}
		return err
	})
	return err
}

//-----------------------------------------------------------------------------

func (f *FsysList) ListFext(fsys fs.FS, fext string) error {
	err := fs.WalkDir(fsys, ".", func(path string, d Entry, err error) error {
		mod := d.Type()
		if err == nil && mod.IsRegular() && strings.HasSuffix(path, fext) {
			f.Append(fsys, path, d.Name(), mod)
		}
		return err
	})
	return err
}

//-----------------------------------------------------------------------------

func (f *FsysList) String() string {
	if f.list == nil {
		return "null"
	}
	switch n := len(f.list); n {
	case 0:
		return "[]"
	case 1:
		return fmt.Sprintf("[%q]", f.list[0].Path)
	default:
		b := make([]byte, 0, 32)
		b = append(b, '[', '\n', ' ', ' ')
		for i := 0; i < n; i++ {
			if i > 0 {
				b = append(b, ',', '\n', ' ', ' ')
			}
			b = strconv.AppendQuote(b, f.list[i].Path)
		}
		b = append(b, '\n', ']')
		return *(*string)(unsafe.Pointer(&b))
	}
}
