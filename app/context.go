// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package app

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
)

// SignalContext interface represents the internal context.Context implementation that
// hooks to SIGINT, SIGTERM, and SIGHUP operating system signals for graceful closing.
//
// SignalContext's methods can be called concurrently multiple times from any goroutine.
type SignalContext interface {
	context.Context

	// Stop stops the context by closing the context.Done() channel if necessary.
	Stop()

	// WaitDone pauses execution of the calling goroutine until the
	// Context.Done() channel closed.
	//
	// Returns immediately if the context has already been stopped.
	WaitDone()

	// WaitStop pauses execution of the calling goroutine until the SignalContext implementation
	// finishes all background tasks and transitions to invalid state.
	//
	// WaitStop will wait for all AsyncRun tasks to finish execution.
	WaitStop()

	// AsyncRun executes f asynchronously in a separate goroutine.
	//
	// WaitStop will not return until all AsyncRun tasks finish.
	//
	// Tip: Use for concurrent shutdown of app services after the context stopped.
	AsyncRun(f func())

	// Running reports if the context is in the running state.
	Running() bool

	// Stopped reports if the context stopped.
	Stopped() bool

	// Invalid indicates that the context finished execution and no longer valid for use.
	Invalid() bool
}

const (
	xrunning = 1
	xstopped = 2
	xinvalid = 3
)

// SignalContext internal implementation.
type sigctx struct {
	flag int32
	wait sync.WaitGroup // stop wait group
	context.Context
	stop context.CancelFunc
	sigs chan os.Signal
	lock sync.Mutex
}

var _ctx sigctx

// DebugContext function pointer.
//
// Set to a valid format function to enable SignalContext
// debug information output. Default is nil.
//
// See: https://pkg.go.dev/fmt for format specs.
var DebugContext func(format string, args ...interface{})

const ctxDebugPrefix = "SignalContext"

// Context returns SignalContext singleton.
//
// It lazily creates the implementer when called for the first time.
func Context() SignalContext {
	return _ctx.context()
}

// Created reports if the SignalContext singleton initialized.
func Created() bool {
	return _ctx.Created()
}

// Running reports if the SignalContext is in the running state.
func Running() bool {
	return _ctx.Running()
}

// Stopped reports that the SignalContext is stopped.
func Stopped() bool {
	return _ctx.Stopped()
}

// Done checks if the context.Context is done.
func Done(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}

//-----------------------------------------------------------------------------

func (c *sigctx) context() SignalContext {
	if atomic.LoadInt32(&c.flag) == 0 {
		c.create()
	}
	return c
}

//-----------------------------------------------------------------------------

func (c *sigctx) Created() bool {
	f := atomic.LoadInt32(&c.flag)
	return f > 0 && f < xinvalid
}

//-----------------------------------------------------------------------------

func (c *sigctx) Running() bool {
	return atomic.LoadInt32(&c.flag) == xrunning
}

//-----------------------------------------------------------------------------

func (c *sigctx) Stopped() bool {
	return atomic.LoadInt32(&c.flag) > xrunning
}

//-----------------------------------------------------------------------------

func (c *sigctx) Invalid() bool {
	return atomic.LoadInt32(&c.flag) == xinvalid
}

//-----------------------------------------------------------------------------

func (c *sigctx) Stop() {
	if c.tryStop() {
		c.debugf("%s: Stop()\n", ctxDebugPrefix)
		c.stop()
	}
}

//-----------------------------------------------------------------------------

func (c *sigctx) WaitDone() {
	if c.Running() {
		<-c.Done()
	}
}

//-----------------------------------------------------------------------------

func (c *sigctx) WaitStop() {
	c.wait.Wait()
}

//-----------------------------------------------------------------------------

func (c *sigctx) AsyncRun(f func()) {
	if f == nil {
		panic("RunAsync func argument is nil")
	}
	c.wait.Add(1)
	go func() {
		defer c.wait.Done() // user defer to safeguard from f() panic
		f()
	}()
}

//-----------------------------------------------------------------------------

func (c *sigctx) create() {
	c.lock.Lock()
	if c.flag > 0 {
		c.lock.Unlock()
		return
	}
	defer c.debugf("%s: created\n", ctxDebugPrefix)
	defer c.lock.Unlock()
	c.Context, c.stop = context.WithCancel(context.Background())
	c.sigs = make(chan os.Signal, 1)
	c.wait.Add(1)
	signal.Notify(c.sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	go c.waitStop()
	atomic.StoreInt32(&c.flag, xrunning)
}

//-----------------------------------------------------------------------------

func (c *sigctx) tryStop() bool {
	return atomic.CompareAndSwapInt32(&c.flag, xrunning, xstopped)
}

//-----------------------------------------------------------------------------

func (c *sigctx) debugf(format string, args ...interface{}) {
	if DebugContext != nil {
		DebugContext(format, args...)
	}
}

//-----------------------------------------------------------------------------

func (c *sigctx) waitStop() {
	select {
	case sig := <-c.sigs:
		c.debugf("%s: %s\n", ctxDebugPrefix, sig)
	case <-c.Done():
		c.debugf("%s: Done()\n", ctxDebugPrefix)
	}
	c.Stop()
	atomic.StoreInt32(&c.flag, xinvalid)
	c.wait.Done()
	// Ignore further signals
	for {
		sig := <-c.sigs
		c.debugf("%s: IGNORE: %s\n", ctxDebugPrefix, sig)
	}
}
