// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package dir

import "os"

type (
	FileInfo = os.FileInfo
	FileMode = os.FileMode
)

func Stat(name string) (FileInfo, error) {
	return os.Stat(name)
}

func Lstat(name string) (FileInfo, error) {
	return os.Lstat(name)
}

func IsDir(name string) bool {
	inf, err := Stat(name)
	return err == nil && inf.Mode()&os.ModeDir != 0
}

func IsFile(name string) bool {
	inf, err := Stat(name)
	return err == nil && inf.Mode()&os.ModeType == 0
}

func IsLink(name string) bool {
	inf, err := Lstat(name)
	return err == nil && inf.Mode()&os.ModeSymlink != 0
}

func FileSize(name string) int64 {
	inf, err := Stat(name)
	if err == nil {
		if inf.Mode().IsRegular() {
			return inf.Size()
		}
		return -1
	}
	return -2
}
