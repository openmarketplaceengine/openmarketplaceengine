// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package cfg

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPassword_MarshalText(t *testing.T) {
	data := map[Password]string{
		"":       "\"\"\n",
		"secret": "'****************'\n",
	}
	for k, v := range data {
		y := string(yamlBytes(t, k))
		require.Equal(t, v, y)
	}
}

func TestPassword_String(t *testing.T) {
	data := map[Password]string{
		"":       "",
		"secret": mockPass,
	}
	for k, v := range data {
		s := fmt.Sprintf("%s", k) //nolint
		require.Equal(t, v, s)
	}
}
