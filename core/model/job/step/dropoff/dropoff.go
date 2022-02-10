package dropoff

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

type DropOff struct {
	f *fsm.FSM
}

func (p *DropOff) Complete() {
	p.f.Event(CompletedEvent)
}

func NewDropOff(onExit func()) *DropOff {
	return &DropOff{
		f: newDropOffFsm(onExit),
	}
}

func newDropOffFsm(onExit func()) *fsm.FSM {
	f := fsm.New(ReadyState)
	f.Transition(fsm.On(CompletedEvent), fsm.Src(ReadyState), fsm.Dst(CompletedState))
	f.Exit(func(state fsm.State) {
		onExit()
	})
	return f
}
