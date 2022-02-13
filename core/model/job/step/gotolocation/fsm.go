package gotolocation

import (
	"fmt"

	"github.com/cocoonspace/fsm"
)

func checkTransition(current State, event Event) error {
	f := fsm.New(fsm.State(current))
	f.Transition(fsm.On(fsm.Event(MoveEvent)), fsm.Src(fsm.State(New)), fsm.Dst(fsm.State(Moving)))
	f.Transition(fsm.On(fsm.Event(NearEvent)), fsm.Src(fsm.State(Moving)), fsm.Dst(fsm.State(Near)))
	f.Transition(fsm.On(fsm.Event(ArrivedEvent)), fsm.Src(fsm.State(Near)), fsm.Dst(fsm.State(Arrived)))
	f.Transition(fsm.On(fsm.Event(CancelledEvent)), fsm.Src(fsm.State(New)), fsm.Dst(fsm.State(Cancelled)))
	f.Transition(fsm.On(fsm.Event(CancelledEvent)), fsm.Src(fsm.State(Moving)), fsm.Dst(fsm.State(Cancelled)))
	f.Transition(fsm.On(fsm.Event(CancelledEvent)), fsm.Src(fsm.State(Near)), fsm.Dst(fsm.State(Cancelled)))
	f.Transition(fsm.On(fsm.Event(CancelledEvent)), fsm.Src(fsm.State(Arrived)), fsm.Dst(fsm.State(Cancelled)))

	ok := f.Event(fsm.Event(event))
	if !ok {
		return fmt.Errorf("illegal transition from state=%v by event=%v", current, event)
	}
	return nil
}
