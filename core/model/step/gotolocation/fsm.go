package gotolocation

import (
	"github.com/openmarketplaceengine/openmarketplaceengine/core/model/step"

	"github.com/cocoonspace/fsm"
)

const (
	NearEvent fsm.Event = iota
	ArriveEvent
	CancelEvent
)

var actionToEvent = map[step.Action]fsm.Event{
	NearAction:   NearEvent,
	ArriveAction: ArriveEvent,
	CancelAction: CancelEvent,
}

const (
	MovingState fsm.State = iota
	NearState
	ArrivedState
	CanceledState
)

var statusToState = map[step.State]fsm.State{
	Moving:   MovingState,
	Near:     NearState,
	Arrived:  ArrivedState,
	Canceled: CanceledState,
}

var stateToStatus = map[fsm.State]step.State{
	MovingState:   Moving,
	NearState:     Near,
	ArrivedState:  Arrived,
	CanceledState: Canceled,
}

var statusToAvailableActions = map[step.State][]step.Action{
	Moving:   {NearAction, CancelAction},
	Near:     {ArriveAction, CancelAction},
	Arrived:  {},
	Canceled: {},
}

func newFsm(current step.State) *fsm.FSM {
	f := fsm.New(statusToState[current])
	f.Transition(fsm.On(NearEvent), fsm.Src(MovingState), fsm.Dst(NearState))
	f.Transition(fsm.On(ArriveEvent), fsm.Src(NearState), fsm.Dst(ArrivedState))
	f.Transition(fsm.On(CancelEvent), fsm.Src(MovingState), fsm.Dst(CanceledState))
	f.Transition(fsm.On(CancelEvent), fsm.Src(NearState), fsm.Dst(CanceledState))
	return f
}
