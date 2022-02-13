package fsm

import (
	"github.com/cocoonspace/fsm"
)

const (
	Moving fsm.Event = iota
	Near
	Arrived
	Canceled
)

const (
	MovingState fsm.State = iota
	NearState
	ArrivedState
	CanceledState
)

type GoToLocationFSM struct {
	f *fsm.FSM
}

func NewGoToLocationFSM(enterStateHooks map[fsm.State]func()) *GoToLocationFSM {

	f := fsm.New(MovingState)
	f.Transition(fsm.On(Near), fsm.Src(MovingState), fsm.Dst(NearState))
	f.Transition(fsm.On(Arrived), fsm.Src(NearState), fsm.Dst(ArrivedState))
	f.Transition(fsm.On(Canceled), fsm.Src(MovingState), fsm.Dst(CanceledState))
	f.Transition(fsm.On(Canceled), fsm.Src(NearState), fsm.Dst(CanceledState))
	f.Transition(fsm.On(Canceled), fsm.Src(ArrivedState), fsm.Dst(CanceledState))

	for state, fun := range enterStateHooks {
		f.EnterState(state, fun)
	}

	return &GoToLocationFSM{
		f: f,
	}
}
