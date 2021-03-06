// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package dao

import (
	"fmt"
	"sync/atomic"
	"unsafe"

	"github.com/rs/xid"
)

type (
	// SUID represents opaque UUID string coming from external sources.
	// Neither format, nor length known in advance.
	SUID = string
	// XUID represents locally generated, sortable 12-bytes UUID.
	XUID = string
)

func NewXid() XUID {
	buf, _ := xid.New().MarshalText()
	return *(*string)(unsafe.Pointer(&buf))
}

//-----------------------------------------------------------------------------

var mockUUID uint32

// MockUUID returns continuous pseudo UUID's for testing.
func MockUUID() string {
	return fmt.Sprintf("%08x", atomic.AddUint32(&mockUUID, 1))
}
