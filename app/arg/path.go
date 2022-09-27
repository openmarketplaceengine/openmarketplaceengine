// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package arg

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/driverscooperative/geosrv/app/dir"
)

type PathFlag uint

const (
	FileMustExist PathFlag = 1 << iota
	FileSkipInvalid
	PathPrintError
	DirCreate
)

var (
	errEmptyFilePath = errors.New("empty file path")
	errEmptyFileList = errors.New("empty file list")
)

//-----------------------------------------------------------------------------

func (f *FlagSet) Files(name string, help string, flags PathFlag, call func(ctx Context, files []string) error) *Flag {
	checkNameCall(name, call == nil)
	flag := f.Var(new(stringValue), name, help)
	flag.fset = func(ctx Context, set *FlagSet, val string) error {
		aval := flag.aval
		if aval == nil {
			aval = FileValidator(FileMustExist)
		}
		var args []string
		if len(val) > 0 {
			args = []string{val}
		} else {
			args = set.rest()
		}
		files, err := checkFiles(aval, flags, args)
		if err != nil {
			return err
		}
		return call(ctx, files)
	}
	return flag
}

//-----------------------------------------------------------------------------

func (f *FlagSet) Dir(name string, help string, flags PathFlag, call func(ctx Context, path string) error) *Flag {
	checkNameCall(name, call == nil)
	flag := f.Var(new(stringValue), name, help)
	flag.fset = func(ctx Context, set *FlagSet, val string) error {
		if len(val) == 0 {
			val, _ = set.next()
		}
		if len(val) == 0 {
			return fmt.Errorf("flag needs an argument: -%s", name)
		}
		// current working directory
		if len(val) == 1 && val[0] == '.' {
			cwd, err := os.Getwd()
			if err == nil {
				err = call(ctx, cwd)
			}
			return err
		}
		var err error
		val, err = dir.Abs(val)
		if err != nil {
			return err
		}
		var inf dir.FileInfo
		inf, err = dir.Stat(val)
		if err != nil {
			if os.IsNotExist(err) {
				if flags.Has(DirCreate) {
					err = dir.MkdirAll(val, 0700)
					if err == nil {
						return call(ctx, val)
					}
				}
			}
			return pathError("directory", val, err)
		}
		if !inf.IsDir() {
			return fmt.Errorf("not a directory: %q", val)
		}
		return call(ctx, val)
	}
	return flag
}

//-----------------------------------------------------------------------------

func checkFiles(valid Validator, flags PathFlag, args []string) ([]string, error) {
	alen := len(args)
	if alen == 0 {
		return nil, errEmptyFileList
	}
	eskip := flags&FileSkipInvalid != 0
	debug := flags&PathPrintError != 0
	files := make([]string, 0, alen)
	for i := 0; i < alen; i++ {
		arg, err := valid(args[i])
		if err != nil {
			if debug {
				println(err.Error())
			}
			if eskip {
				continue
			}
			return nil, err
		}
		files = append(files, arg)
	}
	if len(files) == 0 {
		return nil, errEmptyFileList
	}
	return files, nil
}

//-----------------------------------------------------------------------------

func FileValidator(flags PathFlag, fexts ...string) Validator {
	return func(arg string) (string, error) {
		if len(arg) == 0 {
			return "", errEmptyFilePath
		}
		if n := len(fexts); n > 0 {
			hasext := false
			for i := 0; i < n; i++ {
				if strings.HasSuffix(arg, fexts[i]) {
					hasext = true
					break
				}
			}
			if !hasext {
				return "", fmt.Errorf("invalid file extension: %q", arg)
			}
		}
		if flags&FileMustExist != 0 {
			inf, err := dir.Stat(arg)
			if err != nil {
				return "", pathError("file", arg, err)
			}
			if !inf.Mode().IsRegular() {
				return "", fmt.Errorf("not a file: %q", arg)
			}
			if inf.Size() == 0 {
				return "", fmt.Errorf("empty file: %q", arg)
			}
		}
		return arg, nil
	}
}

//-----------------------------------------------------------------------------

func pathError(prefix, arg string, err error) error {
	if os.IsNotExist(err) {
		return fmt.Errorf("%s does not exist: %q", prefix, arg)
	}
	if os.IsExist(err) {
		return fmt.Errorf("%s already exists: %q", prefix, arg)
	}
	if os.IsPermission(err) {
		return fmt.Errorf("permission denied: %q", arg)
	}
	if os.IsTimeout(err) {
		return fmt.Errorf("timeout occurred: %q", arg)
	}
	if pe, ok := err.(*os.PathError); ok {
		err = pe.Err
	}
	return fmt.Errorf("%s error %w: %q", prefix, err, arg)
}

//-----------------------------------------------------------------------------

func (f PathFlag) Has(flag PathFlag) bool {
	return (f & flag) != 0
}
