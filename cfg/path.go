// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package cfg

import (
	"fmt"
	"os"
	"path/filepath"
)

//-----------------------------------------------------------------------------

func configSearch(fileName string, files []string) ([]string, error) {
	const cfgEnv = EnvPrefix + "_CONF"
	type pathFunc = func() (string, error)
	if cfg, ok := os.LookupEnv(cfgEnv); ok && len(cfg) > 0 {
		files = append(files, cfg)
	}
	run := []pathFunc{workPath, execPath, dataPath}
	for i := range run {
		dir, err := run[i]()
		if err != nil {
			return nil, err
		}
		files = append(files, filepath.Join(dir, fileName))
	}
	return files, nil
}

//-----------------------------------------------------------------------------

func workPath() (dir string, err error) {
	dir, err = os.Getwd()
	if err != nil {
		err = pathFail("os.Getwd()", err)
	}
	return
}

//-----------------------------------------------------------------------------

func execPath() (dir string, err error) {
	dir, err = os.Executable()
	if err != nil {
		err = pathFail("os.Executable()", err)
		return
	}
	var inf os.FileInfo
	inf, err = os.Lstat(dir)
	if err != nil {
		err = pathFail("os.Lstat", err)
		return
	}
	if (inf.Mode() & os.ModeSymlink) != 0 {
		dir, err = filepath.EvalSymlinks(dir)
		if err != nil {
			err = pathFail("EvalSymlinks()", err)
			return
		}
	}
	dir = filepath.Dir(dir)
	return
}

//-----------------------------------------------------------------------------

func dataPath() (dir string, err error) {
	dir, err = os.UserConfigDir()
	if err != nil {
		err = pathFail("os.UserConfigDir()", err)
		return
	}
	dir = filepath.Join(dir, AppName)
	return
}

//-----------------------------------------------------------------------------

func pathFail(name string, err error) error {
	return fmt.Errorf("%s failed: %w", name, err)
}
