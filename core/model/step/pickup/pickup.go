package pickup

import (
	"fmt"

	"github.com/cocoonspace/fsm"
)

const (
	Ready fsm.State = iota
	Completed
	Canceled
)

const (
	Complete fsm.Event = iota
	Cancel
)

var events = map[fsm.State][]fsm.Event{
	readyState:     {Complete, Cancel},
	completedState: {},
	canceledState:  {},
}

type Pickup struct {
	fsm *fsm.FSM
}

func (gtl *Pickup) CurrentState() fsm.State {
	return gtl.fsm.Current()
}

func (gtl *Pickup) AvailableEvents() []fsm.Event {
	state := gtl.fsm.Current()
	return events[state]
}

func (gtl *Pickup) Handle(event fsm.Event) error {
	ok := gtl.fsm.Event(event)
	if !ok {
		state := gtl.fsm.Current()
		return fmt.Errorf("illegal transition from state=%v by event=%v", state, event)
	}
	return nil
}

func New(state fsm.State) (pickup *Pickup) {
	pickup = &Pickup{
		fsm: newFsm(state),
	}

	return
}

const (
	readyState fsm.State = iota
	completedState
	canceledState
)

const (
	completeEvent fsm.Event = iota
	cancelEvent
)

func newFsm(initial fsm.State) *fsm.FSM {
	f := fsm.New(initial)
	f.Transition(fsm.On(completeEvent), fsm.Src(readyState), fsm.Dst(completedState))
	f.Transition(fsm.On(cancelEvent), fsm.Src(readyState), fsm.Dst(canceledState))

	return f
}
