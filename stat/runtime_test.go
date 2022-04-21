// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package stat

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSplitFunc(t *testing.T) {
	nary := make([]string, 0, 8)
	for i := range rstatNames {
		nary = nary[:0]
		name := rstatNames[i][0]
		if len(name) > 0 && name[0] == '/' {
			name = name[1:]
		}
		splitFunc(name, '/', func(part string, path string) bool {
			nary = append(nary, part)
			return true
		})
		sary := strings.Split(name, "/")
		require.Equal(t, sary, nary)
	}
}
