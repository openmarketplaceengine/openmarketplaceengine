// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package dir

import (
	"os"
	"path/filepath"
)

type FilePath struct {
	FullPath string
	FileName string
}

func (p *FilePath) Set(fullPath string) {
	p.FullPath = fullPath
	p.FileName = filepath.Base(fullPath)
}

func (p *FilePath) Assign(path string, name string) {
	p.FullPath = Join(path, name)
	p.FileName = name
}

func (p *FilePath) Empty() bool {
	return len(p.FullPath) == 0
}

//-----------------------------------------------------------------------------

func (p *FilePath) Valid() bool {
	return FileSize(p.FullPath) > 0
}

//-----------------------------------------------------------------------------

func (p *FilePath) Read() ([]byte, error) {
	return os.ReadFile(p.FullPath)
}

//-----------------------------------------------------------------------------

func (p *FilePath) Write(data []byte, perm FileMode) error {
	if perm == 0 {
		perm = UserFilePerm
	}
	return os.WriteFile(p.FullPath, data, perm)
}

//-----------------------------------------------------------------------------

func (p *FilePath) Decode(dst interface{}, dec func(buf []byte, dst interface{}) error) error {
	buf, err := p.Read()
	if err == nil {
		err = dec(buf, dst)
	}
	return err
}

//-----------------------------------------------------------------------------

func (p *FilePath) Encode(src interface{}, perm FileMode, enc func(src interface{}) ([]byte, error)) error {
	buf, err := enc(src)
	if err == nil {
		err = p.Write(buf, perm)
	}
	return err
}
