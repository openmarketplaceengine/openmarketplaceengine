// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package geo

import (
	"encoding/binary"
	"fmt"
	"math"
)

//goland:noinspection SpellCheckingInspection
const hextable = "0123456789ABCDEF"

func hexenc(dst, src []byte) int {
	j := 0
	for _, v := range src {
		dst[j] = hextable[v>>4]
		dst[j+1] = hextable[v&0x0f]
		j += 2
	}
	return len(src) * 2
}

//-----------------------------------------------------------------------------
// Decoding
//-----------------------------------------------------------------------------

func hexdec(dst []byte, src string) (int, error) { //nolint:unparam
	n := len(src)
	if (n & 1) != 0 {
		return 0, fmt.Errorf("hex input length %d is not even", n)
	}
	i, j := 0, 1
	for j < n {
		a := hexlut[src[j-1]]
		if a == 0 {
			return i, InvalidByteError(src[j-1])
		}
		b := hexlut[src[j]]
		if b == 0 {
			return i, InvalidByteError(src[j])
		}
		a--
		b--
		dst[i] = (a << 4) | b
		i++
		j += 2
	}
	return i, nil
}

//-----------------------------------------------------------------------------

func hexdecU32(s string, littleEndian bool) (uint32, error) {
	const minLen = 8
	if len(s) < minLen {
		return 0, fmt.Errorf("hex to uint32 input length %d is too small", len(s))
	}
	var buf [minLen / 2]byte
	dst := buf[:]
	_, err := hexdec(dst, s[:minLen])
	if err != nil {
		return 0, err
	}
	if littleEndian {
		return binary.LittleEndian.Uint32(dst), nil
	}
	return binary.BigEndian.Uint32(dst), nil
}

//-----------------------------------------------------------------------------

func hexdecU64(s string, littleEndian bool) (uint64, error) {
	const minLen = 16
	if len(s) < minLen {
		return 0, fmt.Errorf("hex to uint64 input length %d is too small", len(s))
	}
	var buf [minLen / 2]byte
	dst := buf[:]
	_, err := hexdec(dst, s[:minLen])
	if err != nil {
		return 0, err
	}
	if littleEndian {
		return binary.LittleEndian.Uint64(dst), nil
	}
	return binary.BigEndian.Uint64(dst), nil
}

//-----------------------------------------------------------------------------

func hexdecF64(s string, littleEndian bool) (float64, error) {
	u, err := hexdecU64(s, littleEndian)
	if err != nil {
		return 0, err
	}
	return math.Float64frombits(u), nil
}

//-----------------------------------------------------------------------------

// InvalidByteError values describe errors resulting from an invalid byte in a hex string.
type InvalidByteError byte

func (e InvalidByteError) Error() string {
	return fmt.Sprintf("invalid hex byte: %#U", rune(e))
}

//-----------------------------------------------------------------------------
// HEX Lookup Table
//-----------------------------------------------------------------------------

var hexlut = [256]byte{
	'0': 1, '1': 2, '2': 3, '3': 4, '4': 5,
	'5': 6, '6': 7, '7': 8, '8': 9, '9': 10,
	'a': 11, 'b': 12, 'c': 13, 'd': 14, 'e': 15, 'f': 16,
	'A': 11, 'B': 12, 'C': 13, 'D': 14, 'E': 15, 'F': 16,
}
