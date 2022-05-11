// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package geo

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLocationWKB_UnmarshalText(t *testing.T) {
	var loc LocationWKB
	for i := range wkbPoints {
		p := &wkbPoints[i]
		err := loc.UnmarshalText([]byte(p.src))
		require.NoError(t, err)
		require.True(t, loc.EqualsCoord(p.lat, p.lon))
	}
}
