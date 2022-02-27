package pickup

import (
	"github.com/openmarketplaceengine/openmarketplaceengine/core/model/step"

	"github.com/cocoonspace/fsm"
)

const (
	ReadyState fsm.State = iota
	CompletedState
	CanceledState
)

const (
	CompleteEvent fsm.Event = iota
	CancelEvent
)

var actionToEvent = map[step.Action]fsm.Event{
	CompleteAction: CompleteEvent,
	CancelAction:   CancelEvent,
}

var statusToState = map[step.State]fsm.State{
	Ready:     ReadyState,
	Completed: CompletedState,
	Canceled:  CanceledState,
}

var stateToStatus = map[fsm.State]step.State{
	ReadyState:     Ready,
	CompletedState: Completed,
	CanceledState:  Canceled,
}

var statusToAvailableActions = map[step.State][]step.Action{
	Ready:     {CompleteAction, CancelAction},
	Completed: {},
	Canceled:  {},
}

func newFsm(current step.State) *fsm.FSM {
	f := fsm.New(statusToState[current])
	f.Transition(fsm.On(CompleteEvent), fsm.Src(ReadyState), fsm.Dst(CompletedState))
	f.Transition(fsm.On(CancelEvent), fsm.Src(ReadyState), fsm.Dst(CanceledState))

	return f
}
