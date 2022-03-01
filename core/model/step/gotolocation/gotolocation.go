package gotolocation

import (
	"fmt"

	"github.com/cocoonspace/fsm"
	"github.com/openmarketplaceengine/openmarketplaceengine/core/model/step"
)

const (
	Moving   step.State = "Moving"
	Near     step.State = "Near"
	Arrived  step.State = "Arrived"
	Canceled step.State = "Canceled"
)

const (
	NearBy step.Event = "NearBy"
	Arrive step.Event = "Arrive"
	Cancel step.Event = "Cancel"
)

var fsm2state = map[fsm.State]step.State{
	movingState:   Moving,
	nearState:     Near,
	arrivedState:  Arrived,
	canceledState: Canceled,
}

var events = map[fsm.State][]step.Event{
	movingState:   {NearBy, Cancel},
	nearState:     {Arrive, Cancel},
	arrivedState:  {},
	canceledState: {},
}

var event2fsm = map[step.Event]fsm.Event{
	NearBy: nearEvent,
	Arrive: arriveEvent,
	Cancel: cancelEvent,
}

var state2fsm = map[step.State]fsm.State{
	Moving:   movingState,
	Near:     nearState,
	Arrived:  arrivedState,
	Canceled: canceledState,
}

type GoToLocation struct {
	fsm *fsm.FSM
}

func (gtl *GoToLocation) CurrentState() step.State {
	state := gtl.fsm.Current()
	return fsm2state[state]
}

func (gtl *GoToLocation) AvailableEvents() []step.Event {
	state := gtl.fsm.Current()
	return events[state]
}

func (gtl *GoToLocation) Handle(event step.Event) error {
	ok := gtl.fsm.Event(event2fsm[event])
	if !ok {
		state := fsm2state[gtl.fsm.Current()]
		return fmt.Errorf("illegal transition from state=%v by event=%v", state, event)
	}
	return nil
}

func New(state step.State) *GoToLocation {
	gtl := &GoToLocation{
		fsm: newFsm(state2fsm[state]),
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
