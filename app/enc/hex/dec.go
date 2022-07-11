// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

// Some code in this file is from encoding/hex package.
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hex

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"unsafe"
)

// ErrLength reports an attempt to decode an odd-length input.
var ErrLength = errors.New("odd length hex string")

// ErrShort indicates that the input string is too short for decoding.
var ErrShort = errors.New("hex string is too short")

// InvalidByteError values describe errors resulting from an invalid byte in a hex string.
type InvalidByteError byte

func (e InvalidByteError) Error() string {
	return fmt.Sprintf("invalid hex byte: %#U", rune(e))
}

// DecodedLen returns the length of a decoding of x source bytes.
// Specifically, it returns x / 2.
func DecodedLen(x int) int { return x / 2 }

// Decode decodes src into DecodedLen(len(src)) bytes,
// returning the actual number of bytes written to dst.
//
// Decode expects that src contains only hexadecimal
// characters and that src has even length.
// If the input is malformed, Decode returns the number
// of bytes decoded before the error.
func Decode(dst []byte, src []byte) (int, error) {
	str := *(*string)(unsafe.Pointer(&src))
	return DecodeFromString(dst, str)
}

// DecodeFromString decodes src into DecodedLen(len(src)) bytes,
// returning the actual number of bytes written to dst.
//
// DecodeFromString expects that src contains only hexadecimal
// characters and that src has even length.
// If the input is malformed, DecodeFromString returns the number
// of bytes decoded before the error.
func DecodeFromString(dst []byte, src string) (int, error) {
	n := len(src)
	if (n & 1) != 0 {
		return 0, ErrLength
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

// DecodeString returns the bytes represented by the hexadecimal string s.
//
// DecodeString expects that src contains only hexadecimal
// characters and that src has even length.
// If the input is malformed, DecodeString returns
// the bytes decoded before the error.
func DecodeString(s string) ([]byte, error) {
	dst := make([]byte, len(s)/2)
	n, err := DecodeFromString(dst, s)
	return dst[:n], err
}

//-----------------------------------------------------------------------------
// Numbers
//-----------------------------------------------------------------------------

func DecodeUint32(s string, littleEndian bool) (uint32, error) {
	const minLen = 8
	if len(s) < minLen {
		return 0, ErrShort
	}
	var buf [minLen / 2]byte
	dst := buf[:]
	_, err := DecodeFromString(dst, s[:minLen])
	if err != nil {
		return 0, err
	}
	if littleEndian {
		return binary.LittleEndian.Uint32(dst), nil
	}
	return binary.BigEndian.Uint32(dst), nil
}

//-----------------------------------------------------------------------------

func DecodeUint64(s string, littleEndian bool) (uint64, error) {
	const minLen = 16
	if len(s) < minLen {
		return 0, ErrShort
	}
	var buf [minLen / 2]byte
	dst := buf[:]
	_, err := DecodeFromString(dst, s[:minLen])
	if err != nil {
		return 0, err
	}
	if littleEndian {
		return binary.LittleEndian.Uint64(dst), nil
	}
	return binary.BigEndian.Uint64(dst), nil
}

//-----------------------------------------------------------------------------

func DecodeFloat64(s string, littleEndian bool) (float64, error) {
	u, err := DecodeUint64(s, littleEndian)
	if err != nil {
		return 0, err
	}
	return math.Float64frombits(u), nil
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
