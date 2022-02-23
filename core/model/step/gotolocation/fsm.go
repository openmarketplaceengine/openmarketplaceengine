package gotolocation

import (
	"github.com/openmarketplaceengine/openmarketplaceengine/core/model/step"

	"github.com/cocoonspace/fsm"
)

const (
	MoveEvent fsm.Event = iota
	NearEvent
	ArriveEvent
	CancelEvent
)

var actionToEvent = map[step.Action]fsm.Event{
	Move:   MoveEvent,
	Near:   NearEvent,
	Arrive: ArriveEvent,
	Cancel: CancelEvent,
}

const (
	NewState fsm.State = iota
	MovingState
	NearByState
	ArrivedState
	CanceledState
)

var statusToState = map[step.Status]fsm.State{
	New:      NewState,
	Moving:   MovingState,
	NearBy:   NearByState,
	Arrived:  ArrivedState,
	Canceled: CanceledState,
}

var stateToStatus = map[fsm.State]step.Status{
	NewState:      New,
	MovingState:   Moving,
	NearByState:   NearBy,
	ArrivedState:  Arrived,
	CanceledState: Canceled,
}

var statusToAvailableActions = map[step.Status][]step.Action{
	New:      {Move, Cancel},
	Moving:   {Near, Cancel},
	NearBy:   {Arrive, Cancel},
	Arrived:  {},
	Canceled: {},
}

func newFsm(current step.Status) *fsm.FSM {
	f := fsm.New(statusToState[current])
	f.Transition(fsm.On(MoveEvent), fsm.Src(NewState), fsm.Dst(MovingState))
	f.Transition(fsm.On(NearEvent), fsm.Src(MovingState), fsm.Dst(NearByState))
	f.Transition(fsm.On(ArriveEvent), fsm.Src(NearByState), fsm.Dst(ArrivedState))
	f.Transition(fsm.On(CancelEvent), fsm.Src(NewState), fsm.Dst(CanceledState))
	f.Transition(fsm.On(CancelEvent), fsm.Src(MovingState), fsm.Dst(CanceledState))
	f.Transition(fsm.On(CancelEvent), fsm.Src(NearByState), fsm.Dst(CanceledState))
	f.Transition(fsm.On(CancelEvent), fsm.Src(ArrivedState), fsm.Dst(CanceledState))
	return f
}
