// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package cfg

import (
	"fmt"
	"strings"
)

type MultiError struct {
	err []error
}

//-----------------------------------------------------------------------------

func (e *MultiError) Add(err error) error {
	if err != nil {
		e.err = append(e.err, err)
	}
	return err
}

//-----------------------------------------------------------------------------

func (e *MultiError) Formatf(format string, args ...interface{}) error {
	return e.Add(fmt.Errorf(format, args...))
}

//-----------------------------------------------------------------------------

func (e *MultiError) WithLast(f func(err error) error) error {
	if n := len(e.err) - 1; n >= 0 {
		e.err[n] = f(e.err[n])
		return e.err[n]
	}
	return nil
}

//-----------------------------------------------------------------------------

func (e *MultiError) Last() error {
	if n := len(e.err); n > 0 {
		return e.err[n-1]
	}
	return nil
}

//-----------------------------------------------------------------------------

func (e *MultiError) Error() string {
	n := len(e.err)
	if n == 0 {
		return "<empty error>"
	}
	if n == 1 {
		return e.err[0].Error()
	}
	var b strings.Builder
	b.Grow(n * 32)
	for i := 0; i < n; i++ {
		if i > 0 {
			b.WriteByte('\n')
		}
		b.WriteString(e.err[i].Error())
	}
	return b.String()
}

//-----------------------------------------------------------------------------

func (e *MultiError) Return() error {
	if len(e.err) > 0 {
		return e
	}
	return nil
}

//-----------------------------------------------------------------------------

type CantStop string

func (e CantStop) Error() string {
	return fmt.Sprintf("cannot stop %s without successful boot", string(e))
}
