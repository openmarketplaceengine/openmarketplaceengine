// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package hex

import (
	"encoding/binary"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

//-----------------------------------------------------------------------------

func TestDecodeUint32(t *testing.T) {
	b := make([]byte, 4)
	for i := uint32(0); i < 1000; i++ {
		binary.LittleEndian.PutUint32(b, i)
		s := EncodeToString(b)
		v, err := DecodeUint32(s, true)
		require.NoError(t, err)
		require.Equal(t, i, v)
	}
}

func TestDecodeUint64(t *testing.T) {
	b := make([]byte, 8)
	for i := uint64(0); i < 1000; i++ {
		binary.LittleEndian.PutUint64(b, i)
		s := EncodeToString(b)
		v, err := DecodeUint64(s, true)
		require.NoError(t, err)
		require.Equal(t, i, v)
	}
}

func TestDecodeFloat64(t *testing.T) {
	b := make([]byte, 8)
	for i := float64(0); i < 1000; i++ {
		u := math.Float64bits(i)
		binary.LittleEndian.PutUint64(b, u)
		s := EncodeToString(b)
		f, err := DecodeFloat64(s, true)
		require.NoError(t, err)
		require.Equal(t, i, f)
	}
}
