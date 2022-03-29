// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package cfg

import (
	"fmt"
	"sync/atomic"
)

type State int

const (
	StateUnused = iota
	StateBooting
	StateRunning
	StateClosing
	StateStopped
	StateFailure
)

type State64 struct {
	state int64
}

func (s *State64) Current() State {
	return State(atomic.LoadInt64(&s.state))
}

func (s *State64) Running() bool {
	return atomic.LoadInt64(&s.state) == StateRunning
}

func (s *State64) Invalid() bool {
	return atomic.LoadInt64(&s.state) != StateRunning
}

func (s *State64) TryBoot() bool {
	return s.TrySet(StateUnused, StateBooting)
}

func (s *State64) TryStop() bool {
	return s.TrySet(StateRunning, StateClosing)
}

func (s *State64) SetUnused() {
	atomic.StoreInt64(&s.state, StateUnused)
}

func (s *State64) SetRunning() {
	atomic.StoreInt64(&s.state, StateRunning)
}

func (s *State64) SetStopped() {
	atomic.StoreInt64(&s.state, StateStopped)
}

func (s *State64) SetFailure() {
	atomic.StoreInt64(&s.state, StateFailure)
}

func (s *State64) TrySet(old State, new State) bool { //nolint
	return atomic.CompareAndSwapInt64(&s.state, int64(old), int64(new))
}

func (s *State64) BootOrFail(err *error) {
	if err != nil && *err != nil {
		s.SetFailure()
		return
	}
	s.SetRunning()
}

func (s *State64) StopOrFail(f func() error) error {
	if err := f(); err != nil {
		s.SetFailure()
		return err
	}
	s.SetStopped()
	return nil
}

func (s *State64) ResetIfStopped() bool {
	return s.TrySet(StateStopped, StateUnused)
}

//-----------------------------------------------------------------------------

func (s *State64) StateError(prefix string) error {
	return &StateError{s.Current(), prefix}
}

//-----------------------------------------------------------------------------

func (s State) String() string {
	switch s {
	case StateUnused:
		return "unused"
	case StateBooting:
		return "booting"
	case StateRunning:
		return "running"
	case StateClosing:
		return "closing"
	case StateStopped:
		return "stopped"
	case StateFailure:
		return "failure"
	}
	return fmt.Sprintf("State[%d]", s)
}

//-----------------------------------------------------------------------------
// StateError
//-----------------------------------------------------------------------------

type StateError struct {
	State State
	Sname string
}

func (e *StateError) Error() string {
	if len(e.Sname) > 0 {
		return fmt.Sprintf("%s: invalid state: %s", e.Sname, e.State)
	}
	return fmt.Sprintf("invalid state: %s", e.State)
}
