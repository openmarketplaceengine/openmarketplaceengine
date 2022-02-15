package dropoff

import (
	"context"
	"time"
)

type State int
type Event int

const (
	New State = iota
	Ready
	Completed
	Cancelled
)

func (s State) String() string {
	return [...]string{"New", "Ready", "Completed", "Cancelled"}[s]
}

const (
	ReadyEvent Event = iota
	CompletedEvent
	CancelledEvent
)

func (e Event) String() string {
	return [...]string{"ReadyEvent", "CompletedEvent", "CancelledEvent"}[e]
}

type DropOff struct {
	DriverID         string  `json:",string"`
	DropOffLatitude  float64 `json:",string"`
	DropOffLongitude float64 `json:",string"`
	PassengerIds     string  `json:",string"`
	UpdatedAt        string  `json:",string"`
	State            State   `json:",string"`
}

func NewDropOff(ctx context.Context, driverID string, latitude, longitude float64) (dropOff *DropOff, err error) {
	dropOff = &DropOff{
		DriverID:         driverID,
		DropOffLatitude:  latitude,
		DropOffLongitude: longitude,
		UpdatedAt:        time.Now().Format(time.RFC3339Nano),
		State:            New,
	}

	err = storage.Store(ctx, *dropOff)
	if err != nil {
		return
	}
	return
}

func (p *DropOff) updateState(ctx context.Context, state State) error {
	p.UpdatedAt = time.Now().Format(time.RFC3339Nano)
	p.State = state

	err := storage.Store(ctx, *p)
	if err != nil {
		return err
	}

	return nil
}

func (p *DropOff) Ready(ctx context.Context) error {
	err := checkTransition(p.State, ReadyEvent)
	if err != nil {
		return err
	}
	return p.updateState(ctx, Ready)
}

func (p *DropOff) Complete(ctx context.Context) error {
	err := checkTransition(p.State, CompletedEvent)
	if err != nil {
		return err
	}
	return p.updateState(ctx, Completed)
}

func (p *DropOff) Cancel(ctx context.Context) error {
	err := checkTransition(p.State, CancelledEvent)
	if err != nil {
		return err
	}
	return p.updateState(ctx, Cancelled)
}
