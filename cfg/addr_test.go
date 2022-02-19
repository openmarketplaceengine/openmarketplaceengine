// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package cfg

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddress_String(t *testing.T) {
	data := addrTestData()
	for i := range data {
		kv := data[i]
		require.Equal(t, kv[1], Address(kv[0]).String())
	}
}

//-----------------------------------------------------------------------------

func TestAddress_MarshalText(t *testing.T) {
	data := addrTestData()
	for i := range data {
		kv := data[i]
		tx, _ := Address(kv[0]).MarshalText()
		require.Equal(t, kv[1], string(tx))
	}
}

//-----------------------------------------------------------------------------

func addrTestData() [][2]string {
	//goland:noinspection GrazieInspection
	return [][2]string{
		{"", ""},
		{"user@host", "user@host"},
		{"user:@host", "user:********@host"},
		{"user:pass@host", "user:********@host"},
		{":@host", ":********@host"},
		{":pass@host", ":********@host"},
		{"http://user@host", "http://user@host"},
		{"http://user:pass@host", "http://user:********@host"},
		{"http://host", "http://host"},
	}
}
