// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package dir

import (
	"testing"

	"github.com/stretchr/testify/require"
)

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
