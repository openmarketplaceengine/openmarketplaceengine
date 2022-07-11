// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

// Some code in this file is from encoding/hex package.
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hex

import "unsafe"

type EncodeOption uint32

const (
	UpperCase EncodeOption = 1 << iota
)

//goland:noinspection SpellCheckingInspection
const (
	hexlower = "0123456789abcdef"
	hexupper = "0123456789ABCDEF"
)

// EncodedLen returns the length of an encoding of n source bytes.
// Specifically, it returns n * 2.
func EncodedLen(n int) int { return n * 2 }

// Encode encodes src into EncodedLen(len(src))
// bytes of dst. As a convenience, it returns the number
// of bytes written to dst, but this value is always EncodedLen(len(src)).
// Encode implements hexadecimal encoding.
func Encode(dst, src []byte, opt ...EncodeOption) int {
	hexchr := hexlower
	if hasEncOpt(opt, UpperCase) {
		hexchr = hexupper
	}
	n := len(src)
	for i, j := 0, 0; i < n; i++ {
		c := src[i]
		dst[j] = hexchr[c>>4]
		dst[j+1] = hexchr[c&0x0f]
		j += 2
	}
	return n * 2
}

// EncodeToString returns the hexadecimal encoding of src.
func EncodeToString(src []byte, opt ...EncodeOption) string {
	dst := make([]byte, len(src)*2)
	Encode(dst, src, opt...)
	return *(*string)(unsafe.Pointer(&dst))
}

//-----------------------------------------------------------------------------

func hasEncOpt(all []EncodeOption, opt EncodeOption) bool {
	for i := 0; i < len(all); i++ {
		if (all[i] & opt) != 0 {
			return true
		}
	}
	return false
}
