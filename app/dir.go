// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package app

import (
	"os"

	"github.com/driverscooperative/geosrv/app/dir"
)

type Dir struct {
	Domain   string
	ExecFile dir.FilePath
	ConfFile dir.FilePath
	WorkDir  string
	CacheDir string
}

//-----------------------------------------------------------------------------

func (d *Dir) Init(domain string, appUID string, cfgFile string) error {
	if len(appUID) == 0 {
		panic("appUID argument is empty")
	}
	d.Domain = domain
	appUID = dir.FastJoin(domain, appUID)
	path, err := os.Executable()
	if err != nil {
		return FuncError{"os.Executable", err}
	}
	if dir.IsLink(path) {
		path, err = dir.EvalSymlinks(path)
		if err != nil {
			return FuncError{"dir.EvalSymlinks", err}
		}
	}
	d.ExecFile.Set(path)
	if len(cfgFile) > 0 {
		path, err = pathCall(os.UserConfigDir, appUID, "os.UserConfigDir")
		if err != nil {
			return err
		}
		d.ConfFile.Assign(path, cfgFile)
	}
	d.CacheDir, err = pathCall(os.UserCacheDir, appUID, "os.UserCacheDir")
	if err != nil {
		return err
	}
	d.WorkDir, err = os.Getwd()
	return err
}

//-----------------------------------------------------------------------------

func pathCall(call func() (string, error), join string, callName string) (string, error) {
	path, err := call()
	if err != nil {
		return "", FuncError{callName, err}
	}
	path = dir.FastJoin(path, join)
	err = dir.MkdirAll(path, dir.UserPathPerm)
	if err != nil {
		return "", FuncError{"dir.MkdirAll", err}
	}
	return path, nil
}
