// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package geo

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHexDecode(t *testing.T) {
	for i := range wkbPoints {
		src := wkbPoints[i].src
		d1 := make([]byte, len(src) / 2)
		d2 := make([]byte, len(d1))
		_, err := hex.Decode(d1, []byte(src))
		require.NoError(t, err)
		_, err = hexdec(d2, src)
		require.NoError(t, err)
		require.Equal(t, d1, d2)
	}
}

//-----------------------------------------------------------------------------
// Benchmark
//-----------------------------------------------------------------------------

var hexDst []byte

func BenchmarkHexDec(b *testing.B) {
	const hexSrc = "0101000020E6100000FE0C6FD6E07A52C02A00C633685C4440"
	hexDst = make([]byte, len(hexSrc) / 2)
	b.Run("STD", func(b *testing.B) {
		benchStdHexDec(b, hexDst, hexSrc)
	})
	b.Run("LUT", func(b *testing.B) {
		benchLutHexDec(b, hexDst, hexSrc)
	})
}

func benchStdHexDec(b *testing.B, dst []byte, src string) {
	inp := []byte(src)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := hex.Decode(dst, inp)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func benchLutHexDec(b *testing.B, dst []byte, src string) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := hexdec(dst, src)
		if err != nil {
			b.Fatal(err)
		}
	}
}