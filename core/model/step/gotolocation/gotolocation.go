package gotolocation

import (
	"fmt"

	"github.com/cocoonspace/fsm"
)

const (
	Moving fsm.State = iota
	Near
	Arrived
	Canceled
)

const (
	NearBy fsm.Event = iota
	Arrive
	Cancel
)

var events = map[fsm.State][]fsm.Event{
	movingState:   {NearBy, Cancel},
	nearState:     {Arrive, Cancel},
	arrivedState:  {},
	canceledState: {},
}

type GoToLocation struct {
	fsm *fsm.FSM
}

func (gtl *GoToLocation) CurrentState() fsm.State {
	return gtl.fsm.Current()
}

func (gtl *GoToLocation) AvailableEvents() []fsm.Event {
	state := gtl.fsm.Current()
	return events[state]
}

func (gtl *GoToLocation) Handle(event fsm.Event) error {
	ok := gtl.fsm.Event(event)
	if !ok {
		state := gtl.fsm.Current()
		return fmt.Errorf("illegal transition from state=%v by event=%v", state, event)
	}
	return nil
}

func New(state fsm.State) *GoToLocation {
	gtl := &GoToLocation{
		fsm: newFsm(state),
	}

	return gtl
}

const (
	nearEvent fsm.Event = iota
	arriveEvent
	cancelEvent
)

const (
	movingState fsm.State = iota
	nearState
	arrivedState
	canceledState
)

func newFsm(initial fsm.State) *fsm.FSM {
	f := fsm.New(initial)
	f.Transition(fsm.On(nearEvent), fsm.Src(movingState), fsm.Dst(nearState))
	f.Transition(fsm.On(arriveEvent), fsm.Src(nearState), fsm.Dst(arrivedState))
	f.Transition(fsm.On(cancelEvent), fsm.Src(movingState), fsm.Dst(canceledState))
	f.Transition(fsm.On(cancelEvent), fsm.Src(nearState), fsm.Dst(canceledState))
	return f
}
