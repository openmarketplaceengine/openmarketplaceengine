// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package htp

import "testing"

var x []byte

func BenchmarkUnsafeBytes(b *testing.B) {
	var s = "this is a test string"
	b.Run("Copying", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			x = []byte(s)
			if len(x) != len(s) || x[0] != s[0] {
				b.Fatal()
			}
		}
	})
	b.Run("Unsafe", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			x = unsafeBytes(s)
			if len(x) != len(s) || x[0] != s[0] {
				b.Fatal()
			}
		}
	})
}
