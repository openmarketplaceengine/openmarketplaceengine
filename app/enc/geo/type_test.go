// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package geo

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestType_String(t *testing.T) {
	var types = []struct {
		t Type
		s string
	}{
		{wkbPoint, "POINT"},
		{wkbPoint | wkbZ, "POINT Z"},
		{wkbPoint | wkbM, "POINT M"},
		{wkbPoint | wkbZM, "POINT ZM"},
		{wkbPoint | wkbZM | wkbSRID, "POINT ZM"},
		{wkbMultiPolygon | wkbZM | wkbSRID, "MULTIPOLYGON ZM"},
	}
	for i := range types {
		ts := &types[i]
		require.Equal(t, ts.s, ts.t.String())
	}
}
