package dao

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"path/filepath"
	"unsafe"

	"gopkg.in/yaml.v2"
)

// FS represents a subset of fs.ReadFileFS interface.
type FS interface {
	ReadFile(name string) ([]byte, error)
}

type fsysPath struct {
	fsys fs.FS
	path string
	name string
}

//-----------------------------------------------------------------------------
// FsysExec
//-----------------------------------------------------------------------------

// FsysExec executes each file from the Index list.
type FsysExec struct {
	Fsys  FS
	Path  string
	Index Index
}

//-----------------------------------------------------------------------------

func NewFsysExec(fsys FS, path string, index Index) *FsysExec {
	return &FsysExec{
		Fsys:  fsys,
		Path:  path,
		Index: index,
	}
}

//-----------------------------------------------------------------------------

func (f *FsysExec) Names() ([]string, error) {
	return f.Index.Names(f.Fsys, f.Path)
}

//-----------------------------------------------------------------------------

func (f *FsysExec) Execute(ctx Context, exe Executor) error {
	names, err := f.Names()
	if err != nil {
		return err
	}
	var data []byte
	for i := range names {
		path := fullPath(f.Path, names[i])
		debugf("executing %q", path)
		data, err = f.Fsys.ReadFile(path)
		if err != nil {
			return err
		}
		if len(data) == 0 {
			continue
		}
		stmt := *(*string)(unsafe.Pointer(&data))
		_, err = exe.ExecContext(ctx, stmt)
		if err != nil {
			return fmt.Errorf("failed executing %q: %w", path, err)
		}
	}
	return nil
}

//-----------------------------------------------------------------------------

func (f *FsysExec) Result() Result {
	return nil
}

//-----------------------------------------------------------------------------
// Index File
//-----------------------------------------------------------------------------

type Index string

type indexReader = func(data []byte) ([]string, error)

//-----------------------------------------------------------------------------

func (i Index) Names(fsys FS, path string) ([]string, error) {
	const (
		yamlExt = ".yaml"
		jsonExt = ".json"
	)

	var run indexReader

	ext := filepath.Ext(string(i))

	switch ext {
	case yamlExt:
		run = i.readYAML
	case jsonExt:
		run = i.readJSON
	default:
		return nil, fmt.Errorf("invalid SQL index file extension: %q", ext)
	}

	buf, err := fsys.ReadFile(fullPath(path, string(i)))

	if err != nil {
		return nil, err
	}

	return run(buf)
}

//-----------------------------------------------------------------------------

func (i Index) readYAML(data []byte) ([]string, error) {
	names := make([]string, 0, 16)
	err := yaml.Unmarshal(data, &names)
	return names, err
}

//-----------------------------------------------------------------------------

func (i Index) readJSON(data []byte) ([]string, error) {
	names := make([]string, 0, 16)
	err := json.Unmarshal(data, &names)
	return names, err
}

//-----------------------------------------------------------------------------

func fullPath(path, name string) string {
	n := len(path)
	if n == 0 {
		return name
	}
	if len(name) == 0 {
		return path
	}
	if path[n-1] == '/' {
		return path + name
	}
	return path + "/" + name
}

//-----------------------------------------------------------------------------
// file path
//-----------------------------------------------------------------------------

func (f *fsysPath) String() string {
	return f.path
}
