// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package cfg

import (
	"fmt"
	"reflect"
	"sync"
)

// BootFunc is an aliases for boot function.
type BootFunc = func() error

// Bootable interface describes a generic named service
// with support for Boot/Stop lifecycle.
type Bootable interface {
	Boot() error
	Stop() error
}

type BootStep uint

// bootCall struct holds information about a bootable service.
type bootCall struct {
	// name identifies boot call for debug/error reporting
	name string
	// boot called when starting services
	boot BootFunc
	// stop called when stopping services
	stop BootFunc
	// step holds current boot phase
	step BootStep
}

const (
	stepNone = iota
	stepBoot
	stepStop
	stepFail
)

// BootList holds an array of Bootable services.
type BootList struct {
	lock   sync.Mutex
	call   []*bootCall
	step   BootStep
	Debugf PrintFunc
	Errorf PrintFunc
}

// Use appends a Bootable object to the list.
func (b *BootList) Use(name string, boot Bootable) *BootList {
	if boot == nil {
		panic("Bootable object argument is nil")
	}
	if len(name) == 0 {
		name = reflect.TypeOf(boot).String()
	}
	return b.Add(name, boot.Boot, boot.Stop)
}

// Add appends boot functions the list.
func (b *BootList) Add(name string, boot BootFunc, stop BootFunc) *BootList {
	if len(name) == 0 {
		panic("empty name argument")
	}
	if boot == nil && stop == nil {
		panic("both boot and stop functions are nil")
	}
	if b.call == nil {
		b.call = make([]*bootCall, 0, 8)
	}
	c := &bootCall{name: name, boot: boot, stop: stop}
	b.call = append(b.call, c)
	return b
}

// SetLog assigns Debug/Error print functions.
func (b *BootList) SetLog(debugf, errorf PrintFunc) {
	b.Debugf = debugf
	b.Errorf = errorf
}

// Boot starts boot sequence.
func (b *BootList) Boot() (err error) {
	b.lock.Lock()
	defer b.lock.Unlock()
	if b.step != stepNone {
		err = fmt.Errorf("cannot boot with unexpected state: %q", b.step)
		return
	}
	b.step = stepBoot
	for i := 0; i < len(b.call); i++ {
		c := b.call[i]
		if c.boot != nil {
			b.debugf("BOOT %s", c.name)
			if err := c.boot(); err != nil {
				err = &BootError{stepBoot, c.name, err}
				c.step = stepFail
				b.debugf("BOOT aborting...")
				_ = b.stop(i + 1)
				return err
			}
		}
		c.step = stepBoot
	}
	return
}

// Stop stops boot sequence.
func (b *BootList) Stop() (err error) {
	b.lock.Lock()
	defer b.lock.Unlock()
	if b.step != stepBoot {
		err = fmt.Errorf("cannot stop with unexpected state: %q", b.step)
		b.errorf("%s", err)
		return
	}
	b.step = stepStop
	err = b.stop(len(b.call))
	if err != nil {
		b.step = stepFail
	}
	return
}

//-----------------------------------------------------------------------------

func (b *BootList) stop(downFrom int) error {
	var mulerr MultiError
	for i := downFrom; i > 0; i-- {
		c := b.call[i-1]
		if c.step != stepBoot {
			continue
		}
		if c.stop != nil {
			b.debugf("STOP %s", c.name)
			if err := c.stop(); err != nil {
				err = mulerr.Add(&BootError{stepStop, c.name, err})
				b.errorf("%s", err)
				c.step = stepFail
				continue
			}
		}
		c.step = stepStop
	}
	return mulerr.Return()
}

//-----------------------------------------------------------------------------
// Logging
//-----------------------------------------------------------------------------

func (b *BootList) debugf(format string, args ...interface{}) {
	if b.Debugf != nil {
		b.Debugf(format, args...)
	}
}

func (b *BootList) errorf(format string, args ...interface{}) { //nolint
	if b.Errorf != nil {
		b.Errorf(format, args...)
	}
}

//-----------------------------------------------------------------------------
// boot step
//-----------------------------------------------------------------------------

func (s BootStep) String() string {
	switch s {
	case stepNone:
		return "NONE"
	case stepBoot:
		return "BOOT"
	case stepStop:
		return "STOP"
	case stepFail:
		return "FAIL"
	}
	return "<invalid step>"
}

//-----------------------------------------------------------------------------

type BootError struct {
	Step BootStep
	Name string
	Err  error
}

func (e *BootError) Error() string {
	return fmt.Sprintf("%s %s: %s", e.Step, e.Name, e.Err)
}

func (e *BootError) Unwrap() error {
	return e.Err
}
