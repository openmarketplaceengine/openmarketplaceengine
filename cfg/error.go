// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package cfg

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"unsafe"
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

//-----------------------------------------------------------------------------

type EmptyError string

func (e EmptyError) Error() string {
	return fmt.Sprintf("%s is empty", string(e))
}

//-----------------------------------------------------------------------------

type ConstError string

func (e ConstError) Error() string {
	return string(e)
}

//-----------------------------------------------------------------------------

type PanicError struct {
	Cause error
	Stack string
}

//-----------------------------------------------------------------------------

func NewPanicError(rec interface{}) error {
	stack := TracePanic()
	if err, ok := rec.(error); ok {
		return &PanicError{err, stack}
	}
	return &PanicError{
		fmt.Errorf("panic: %v", rec),
		stack,
	}
}

//-----------------------------------------------------------------------------

func (e *PanicError) Error() string {
	if len(e.Stack) > 0 {
		return fmt.Sprintf("%s\n%s", e.Cause, e.Stack)
	}
	return fmt.Sprintf("%s", e.Cause)
}

//-----------------------------------------------------------------------------

func (e *PanicError) Unwrap() error {
	return e.Cause
}

// CatchPanic is to be used with defer to capture panic stack trace.
func CatchPanic(dst *error) {
	if rec := recover(); rec != nil && dst != nil {
		*dst = NewPanicError(rec)
	}
}

//-----------------------------------------------------------------------------

func TracePanic() string {
	const fpanic = "panic.go"
	var b []byte
	c := make([]uintptr, 64)   // callers
	n := runtime.Callers(2, c) // number of callers
	if n == 0 {
		return ""
	}
	frames := runtime.CallersFrames(c[:n]) // frames
	b = make([]byte, 0, 1024)              // buffer
	var frame runtime.Frame
	skip := true
	more := true
	for more {
		frame, more = frames.Next()
		if skip {
			if strings.HasSuffix(frame.File, fpanic) {
				skip = false
			}
			continue
		}
		if len(frame.Function) > 0 {
			if len(b) > 0 {
				b = append(b, '\n')
			}
			b = append(b, frame.Function...)
			if len(frame.File) > 0 {
				b = append(b, '\n', ' ', ' ', ' ', ' ')
				b = append(b, frame.File...)
				b = append(b, ':')
				b = strconv.AppendInt(b, int64(frame.Line), 10)
			}
		}
	}
	return *(*string)(unsafe.Pointer(&b))
}
