package pickup

import (
	"fmt"

	"github.com/cocoonspace/fsm"
	"github.com/openmarketplaceengine/openmarketplaceengine/core/model/step"
)

const (
	Ready     step.State = "Ready"
	Completed step.State = "Completed"
	Canceled  step.State = "Canceled"
)

const (
	Complete step.Event = "Complete"
	Cancel   step.Event = "Cancel"
)

var event2fsm = map[step.Event]fsm.Event{
	Complete: completeEvent,
	Cancel:   cancelEvent,
}

var state2fsm = map[step.State]fsm.State{
	Ready:     readyState,
	Completed: completedState,
	Canceled:  canceledState,
}

var fsm2state = map[fsm.State]step.State{
	readyState:     Ready,
	completedState: Completed,
	canceledState:  Canceled,
}

var events = map[fsm.State][]step.Event{
	readyState:     {Complete, Cancel},
	completedState: {},
	canceledState:  {},
}

type Pickup struct {
	fsm *fsm.FSM
}

func (gtl *Pickup) CurrentState() step.State {
	state := gtl.fsm.Current()
	return fsm2state[state]
}

func (gtl *Pickup) AvailableEvents() []step.Event {
	state := gtl.fsm.Current()
	return events[state]
}

func (gtl *Pickup) Handle(event step.Event) error {
	ok := gtl.fsm.Event(event2fsm[event])
	if !ok {
		state := fsm2state[gtl.fsm.Current()]
		return fmt.Errorf("illegal transition from state=%v by event=%v", state, event)
	}
	return nil
}

func New(state step.State) (pickup *Pickup) {
	pickup = &Pickup{
		fsm: newFsm(state2fsm[state]),
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
