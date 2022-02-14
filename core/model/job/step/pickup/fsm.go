package pickup

import (
	"fmt"

	"github.com/cocoonspace/fsm"
)

func checkTransition(current State, event Event) error {
	f := fsm.New(fsm.State(current))
	f.Transition(fsm.On(fsm.Event(ReadyEvent)), fsm.Src(fsm.State(New)), fsm.Dst(fsm.State(Ready)))
	f.Transition(fsm.On(fsm.Event(CompletedEvent)), fsm.Src(fsm.State(Ready)), fsm.Dst(fsm.State(Completed)))
	f.Transition(fsm.On(fsm.Event(CancelledEvent)), fsm.Src(fsm.State(Ready)), fsm.Dst(fsm.State(Cancelled)))

	ok := f.Event(fsm.Event(event))
	if !ok {
		return fmt.Errorf("illegal transition from state=%v by event=%v", current, event)
	}
	return nil
}
