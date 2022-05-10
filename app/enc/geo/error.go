// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package geo

import (
	"errors"
	"fmt"
)

var (
	ErrEndian = errors.New("invalid endian format")
	ErrSrcLen = errors.New("invalid input length")
	ErrHexDec = errors.New("invalid hex input")
)

type funcError struct {
	name string
	fail error
}

func (e funcError) Error() string {
	return fmt.Sprintf("%s failed: %v", e.name, e.fail)
}

func (e funcError) Unwrap() error {
	return e.fail
}
