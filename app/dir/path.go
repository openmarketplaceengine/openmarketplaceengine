// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package dir

import (
	"os"
	"path/filepath"
)

const (
	DirSep = string(os.PathSeparator)
)

//-----------------------------------------------------------------------------

func IsAbs(path string) bool {
	return filepath.IsAbs(path)
}

//-----------------------------------------------------------------------------

func Abs(path string) (string, error) {
	return filepath.Abs(path)
}

//-----------------------------------------------------------------------------

func Join(elem ...string) string {
	return filepath.Join(elem...)
}

//-----------------------------------------------------------------------------

func MkdirAll(path string, perm FileMode) error {
	return os.MkdirAll(path, perm)
}

//-----------------------------------------------------------------------------

func EvalSymlinks(path string) (string, error) {
	return filepath.EvalSymlinks(path)
}

//-----------------------------------------------------------------------------

func FastJoin(path string, name string) string {
	plen := len(path)
	if plen == 0 {
		return name
	}
	nlen := len(name)
	if nlen == 0 {
		return path
	}
	psep := IsPathSeparator(path[plen-1])
	nsep := IsPathSeparator(name[0])
	if psep && nsep {
		if nlen == 1 {
			return path
		}
		return path + name[1:]
	}
	if psep || nsep {
		return path + name
	}
	return path + DirSep + name
}

//-----------------------------------------------------------------------------

func IsPathSeparator(c byte) bool {
	return os.IsPathSeparator(c)
}
