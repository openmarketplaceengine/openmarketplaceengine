// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package dir

import (
	"embed"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

//go:embed testdata/intpfx
var testIntPfx embed.FS

func TestFsysList_String(t *testing.T) {
	var list FsysList
	require.Equal(t, "null", list.String())
	list.Alloc(8)
	require.Equal(t, "[]", list.String())
	list.Append(nil, "/path/©", "©", 0)
	require.Equal(t, `["/path/©"]`, list.String())
	list.Append(nil, "/path/2", "2", 0)
	require.Equal(t, "[\n  \"/path/©\",\n  \"/path/2\"\n]", list.String())
}

//-----------------------------------------------------------------------------

func TestFsysPath_IntPrefix(t *testing.T) {
	var file FsysPath
	err := WalkDir(testIntPfx, ".", func(path string, d Entry) error {
		if d.Type().IsRegular() && strings.HasSuffix(path, ".yaml") {
			file.Reset(testIntPfx, path, d)
			pfxInt, _ := file.IntPrefix()
			var ymlInt int
			if err := file.ReadYAML(&ymlInt); err != nil {
				return err
			}
			if pfxInt != ymlInt {
				return fmt.Errorf("Int prefix does not match content: %d <> %d", pfxInt, ymlInt)
			}
		}
		return nil
	})
	require.NoError(t, err)
}
