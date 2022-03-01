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
	Ready:     {Complete, Cancel},
	Completed: {},
	Canceled:  {},
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

func newFsm(initial fsm.State) *fsm.FSM {
	f := fsm.New(initial)
	f.Transition(fsm.On(Complete), fsm.Src(Ready), fsm.Dst(Completed))
	f.Transition(fsm.On(Cancel), fsm.Src(Ready), fsm.Dst(Canceled))

	return f
}
