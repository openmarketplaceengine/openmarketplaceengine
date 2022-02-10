package gotolocation

import (
	"github.com/cocoonspace/fsm"
)

type Location struct {
	Latitude  float64
	Longitude float64
}

const (
	TravelingEvent fsm.Event = iota
	NearDestinationEvent
	AtDestinationEvent
	CompletedEvent
	CanceledEvent
)

const (
	TravelingState fsm.State = iota
	NearDestinationState
	AtDestinationState
	CompletedState
	CanceledState
)

type GoToLocation struct {
	f *fsm.FSM
}

func (p *GoToLocation) Traveling() {
	p.f.Event(TravelingEvent)
}

func (p *GoToLocation) NearDestination() {
	p.f.Event(NearDestinationEvent)
}

func (p *GoToLocation) AtDestination() {
	p.f.Event(AtDestinationEvent)
}

func (p *GoToLocation) Complete() {
	p.f.Event(CompletedEvent)
}

func (p *GoToLocation) Cancel() {
	p.f.Event(CanceledEvent)
}

func NewGoToLocation(onExit func()) *GoToLocation {
	return &GoToLocation{
		f: newGoToLocationFsm(onExit),
	}
}

func newGoToLocationFsm(onExit func()) *fsm.FSM {
	f := fsm.New(TravelingState)
	f.Transition(fsm.On(NearDestinationEvent), fsm.Src(TravelingState), fsm.Dst(NearDestinationState))
	f.Transition(fsm.On(AtDestinationEvent), fsm.Src(NearDestinationState), fsm.Dst(AtDestinationState))
	f.Transition(fsm.On(CompletedEvent), fsm.Src(AtDestinationState), fsm.Dst(CompletedState))
	f.Transition(fsm.On(CanceledEvent), fsm.Src(AtDestinationState), fsm.Dst(CanceledState))
	f.Exit(func(state fsm.State) {
		onExit()
	})
	return f
}
