// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package dao

import (
	"unsafe"

	"github.com/rs/xid"
)

type (
	// XUID represents locally generated, sortable 12-bytes UUID.
	XUID = string
)

func NewXid() XUID {
	buf, _ := xid.New().MarshalText()
	return *(*string)(unsafe.Pointer(&buf))
}
