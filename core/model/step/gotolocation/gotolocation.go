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
	Moving:   {NearBy, Cancel},
	Near:     {Arrive, Cancel},
	Arrived:  {},
	Canceled: {},
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

func newFsm(initial fsm.State) *fsm.FSM {
	f := fsm.New(initial)
	f.Transition(fsm.On(NearBy), fsm.Src(Moving), fsm.Dst(Near))
	f.Transition(fsm.On(Arrive), fsm.Src(Near), fsm.Dst(Arrived))
	f.Transition(fsm.On(Cancel), fsm.Src(Moving), fsm.Dst(Canceled))
	f.Transition(fsm.On(Cancel), fsm.Src(Near), fsm.Dst(Canceled))
	return f
}
