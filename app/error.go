// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package app

import "fmt"

type FuncError struct {
	Func string
	Fail error
}

func (e FuncError) Error() string {
	return fmt.Sprintf("%s failed: %v", e.Func, e.Fail)
}

func (e FuncError) Unwrap() error {
	return e.Fail
}
