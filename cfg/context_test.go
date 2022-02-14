// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package cfg

import (
	"flag"
	"sync"
	"testing"
	"time"
)

var ctxTestWait int

func init() {
	flag.IntVar(&ctxTestWait, "wait", 0, "Wait time in seconds before calling Context.Stop()")
}

//-----------------------------------------------------------------------------

func TestCtx_WaitDone(t *testing.T) {
	const waiters = 8
	var wg sync.WaitGroup
	wg.Add(waiters)
	debugContext(t)
	var c sigctx
	for i := 0; i < waiters; i++ {
		go testWaitDone(t, c.context(), &wg, i)
	}
	stopAfter(c.context(), &wg)
}

//-----------------------------------------------------------------------------

func TestCtx_WaitStop(t *testing.T) {
	const waiters = 6
	var wg sync.WaitGroup
	wg.Add(waiters)
	debugContext(t)
	var c sigctx
	for i := 0; i < waiters; i++ {
		go testWaitStop(t, c.context(), &wg, i)
	}
	stopAfter(c.context(), &wg)
}

//-----------------------------------------------------------------------------
// Test Helpers
//-----------------------------------------------------------------------------

func stopAfter(c SignalContext, g *sync.WaitGroup) {
	if testing.Short() || ctxTestWait <= 0 {
		c.Stop()
		g.Wait()
		c.WaitStop()
		return
	}
	tmr := time.AfterFunc(time.Duration(ctxTestWait)*time.Second, c.Stop)
	g.Wait()
	c.WaitStop()
	tmr.Stop()
}

//-----------------------------------------------------------------------------

func testWaitDone(t testing.TB, c SignalContext, g *sync.WaitGroup, i int) {
	t.Logf("wait: %d", i)
	c.WaitDone()
	g.Done()
	t.Logf("done: %d", i)
}

//-----------------------------------------------------------------------------

func testWaitStop(t testing.TB, c SignalContext, g *sync.WaitGroup, i int) {
	t.Logf("wait: %d", i)
	c.WaitStop()
	g.Done()
	t.Logf("stop: %d", i)
}

//-----------------------------------------------------------------------------

func debugContext(t *testing.T) {
	if testing.Verbose() {
		DebugContext = func(format string, args ...interface{}) {
			t.Helper()
			t.Logf(format, args...)
		}
	}
}

//-----------------------------------------------------------------------------
// Benchmarks
//-----------------------------------------------------------------------------

func BenchmarkContext_Done(b *testing.B) {
	b.Run("Error", BenchmarkContext_DoneErr)
	b.Run("Chan", BenchmarkContext_DoneChan)
	b.Run("Stop", BenchmarkContext_DoneStop)
}

func BenchmarkContext_DoneErr(b *testing.B) {
	benchDone(b, checkDoneErr)
}

func BenchmarkContext_DoneChan(b *testing.B) {
	benchDone(b, checkDoneChan)
}

func BenchmarkContext_DoneStop(b *testing.B) {
	benchDone(b, checkDoneStop)
}

//-----------------------------------------------------------------------------

func benchDone(b *testing.B, f func(ctx SignalContext) bool) {
	var x sigctx
	c := x.context()
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if f(c) {
				b.FailNow()
			}
		}
	})
}

//go:noinline
func checkDoneErr(c SignalContext) bool {
	return c.Err() != nil
}

//go:noinline
func checkDoneChan(c SignalContext) bool {
	select {
	case <-c.Done():
		return true
	default:
		return false
	}
}

//go:noinline
func checkDoneStop(c SignalContext) bool {
	return c.Stopped()
}
