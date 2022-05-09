// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package app

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	domain  = "openmarketplaceengine"
	appName = "test"
	cfgFile = "test.yaml"
)

func TestAppDir_Init(t *testing.T) {
	var dirs Dir
	err := dirs.Init(domain, appName, cfgFile)
	require.NoError(t, err)
	require.FileExists(t, dirs.ExecFile.FullPath)
	require.DirExists(t, dirs.ConfFile.BasePath())
	require.DirExists(t, dirs.CacheDir)
}

func dumpJson(t testing.TB, v interface{}) { //nolint
	buf, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%s\n", buf)
}
