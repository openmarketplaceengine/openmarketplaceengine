// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package dir

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFilePath_BasePath(t *testing.T) {
	tests := [...][2]string{
		{"/", "/"},
		{"/hello", "/"},
		{"/hello/world.file", "/hello/"},
	}
	var f FilePath
	for i := range tests {
		p := &tests[i]
		f.Set(p[0])
		require.Equal(t, p[1], f.BasePath())
	}
}
