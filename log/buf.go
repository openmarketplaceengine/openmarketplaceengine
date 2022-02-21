// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package log

import "unsafe"

type AppendWriter struct {
	buf []byte
}

func (w *AppendWriter) Alloc(size int) {
	w.buf = make([]byte, 0, size)
}

func (w *AppendWriter) Write(p []byte) (n int, err error) {
	if n = len(p); n > 0 {
		w.buf = append(w.buf, p...)
	}
	return
}

func (w *AppendWriter) WriteByte(b byte) error {
	w.buf = append(w.buf, b)
	return nil
}

func (w *AppendWriter) WriteString(s string) (n int, err error) {
	if n = len(s); n > 0 {
		w.buf = append(w.buf, s...)
	}
	return
}

func (w *AppendWriter) String() string {
	if len(w.buf) == 0 {
		return ""
	}
	return *(*string)(unsafe.Pointer(&w.buf))
}
