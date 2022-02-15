package pickup

import (
	"github.com/cocoonspace/fsm"
)

const (
	ReadyEvent fsm.Event = iota
	CompletedEvent
)

const (
	ReadyState fsm.State = iota
	CompletedState
)

type Pickup struct {
	f *fsm.FSM
}

func (p *Pickup) Complete() {
	p.f.Event(CompletedEvent)
}

func NewPickup(onExit func()) *Pickup {
	return &Pickup{
		f: newPickUpFsm(onExit),
	}
}

func newPickUpFsm(onExit func()) *fsm.FSM {
	f := fsm.New(ReadyState)
	f.Transition(fsm.On(CompletedEvent), fsm.Src(ReadyState), fsm.Dst(CompletedState))
	f.Exit(func(state fsm.State) {
		onExit()
	})
	return f
}
