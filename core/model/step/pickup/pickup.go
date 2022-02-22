package pickup

import (
	"context"
	"fmt"
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
	switch s {
	case New:
		return "New"
	case Ready:
		return "Ready"
	case Completed:
		return "Completed"
	case Cancelled:
		return "Cancelled"
	default:
		return fmt.Sprintf("%d", s)
	}
}

const (
	ReadyEvent Event = iota
	CompletedEvent
	CancelledEvent
)

func (e Event) String() string {
	switch e {
	case ReadyEvent:
		return "ReadyEvent"
	case CompletedEvent:
		return "CompletedEvent"
	case CancelledEvent:
		return "CancelledEvent"
	default:
		return fmt.Sprintf("%d", e)
	}
}

type Pickup struct {
	DriverID        string  `json:",string"`
	PickupLatitude  float64 `json:",string"`
	PickupLongitude float64 `json:",string"`
	PassengerIds    string  `json:",string"`
	UpdatedAt       string  `json:",string"`
	State           State   `json:",string"`
}

func NewPickup(ctx context.Context, driverID string, latitude, longitude float64) (pickup *Pickup, err error) {
	pickup = &Pickup{
		DriverID:        driverID,
		PickupLatitude:  latitude,
		PickupLongitude: longitude,
		UpdatedAt:       time.Now().Format(time.RFC3339Nano),
		State:           New,
	}

	err = storage.Store(ctx, *pickup)
	if err != nil {
		return
	}
	return
}

func (p *Pickup) updateState(ctx context.Context, state State) error {
	p.UpdatedAt = time.Now().Format(time.RFC3339Nano)
	p.State = state

	err := storage.Store(ctx, *p)
	if err != nil {
		return err
	}

	return nil
}

func (p *Pickup) Ready(ctx context.Context) error {
	err := checkTransition(p.State, ReadyEvent)
	if err != nil {
		return err
	}
	return p.updateState(ctx, Ready)
}

func (p *Pickup) Complete(ctx context.Context) error {
	err := checkTransition(p.State, CompletedEvent)
	if err != nil {
		return err
	}
	return p.updateState(ctx, Completed)
}

func (p *Pickup) Cancel(ctx context.Context) error {
	err := checkTransition(p.State, CancelledEvent)
	if err != nil {
		return err
	}
	return p.updateState(ctx, Cancelled)
}
