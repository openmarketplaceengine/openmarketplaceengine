// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package dao

import (
	"fmt"
	"strconv"
	"unsafe"
)

type StrBuf struct {
	buf []byte
}

func (b *StrBuf) Alloc(n int) *StrBuf {
	b.buf = make([]byte, 0, n)
	return b
}

func (b *StrBuf) Len() int {
	return len(b.buf)
}

func (b *StrBuf) Cap() int {
	return cap(b.buf)
}

func (b *StrBuf) PutStr(s string) *StrBuf {
	b.buf = append(b.buf, s...)
	return b
}

func (b *StrBuf) PutInt(n int) *StrBuf {
	b.buf = strconv.AppendInt(b.buf, int64(n), 10)
	return b
}

//-----------------------------------------------------------------------------

func (b *StrBuf) Write(p []byte) (n int, err error) {
	if n = len(p); n > 0 {
		b.buf = append(b.buf, p...)
	}
	return
}

//-----------------------------------------------------------------------------

func (b *StrBuf) Printf(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(b, format, args...)
}

//-----------------------------------------------------------------------------

func (b *StrBuf) String() string {
	return *(*string)(unsafe.Pointer(&b.buf))
}
