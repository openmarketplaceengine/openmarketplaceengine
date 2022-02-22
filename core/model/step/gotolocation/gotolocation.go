package gotolocation

import (
	"context"
	"fmt"
	"time"
)

type State int
type Event int

const (
	New State = iota
	Moving
	Near
	Arrived
	Cancelled
)

func (s State) String() string {
	switch s {
	case New:
		return "New"
	case Moving:
		return "Moving"
	case Near:
		return "Near"
	case Arrived:
		return "Arrived"
	case Cancelled:
		return "Cancelled"
	default:
		return fmt.Sprintf("%d", s)
	}
}

const (
	MoveEvent Event = iota
	NearEvent
	ArrivedEvent
	CancelledEvent
)

func (e Event) String() string {
	switch e {
	case MoveEvent:
		return "MoveEvent"
	case NearEvent:
		return "NearEvent"
	case ArrivedEvent:
		return "ArrivedEvent"
	case CancelledEvent:
		return "CancelledEvent"
	default:
		return fmt.Sprintf("%d", e)
	}
}

type GoToLocation struct {
	DriverID             string  `json:",string"`
	DestinationLatitude  float64 `json:",string"`
	DestinationLongitude float64 `json:",string"`
	UpdatedAt            string  `json:",string"`
	UpdatedAtLatitude    float64 `json:",string"`
	UpdatedAtLongitude   float64 `json:",string"`
	State                State   `json:",string"`
}

func NewGoToLocation(ctx context.Context, driverID string, latitude, longitude float64) (gtl *GoToLocation, err error) {
	gtl = &GoToLocation{
		DriverID:             driverID,
		DestinationLatitude:  latitude,
		DestinationLongitude: longitude,
		UpdatedAt:            time.Now().Format(time.RFC3339Nano),
		State:                New,
	}

	err = storage.Store(ctx, *gtl)
	if err != nil {
		return
	}
	return
}

func (gtl *GoToLocation) updateState(ctx context.Context, latitude float64, longitude float64, state State) error {
	gtl.UpdatedAtLatitude = latitude
	gtl.UpdatedAtLongitude = longitude
	gtl.UpdatedAt = time.Now().Format(time.RFC3339Nano)
	gtl.State = state

	err := storage.Store(ctx, *gtl)
	if err != nil {
		return err
	}

	return nil
}

func (gtl *GoToLocation) Moving(ctx context.Context, latitude, longitude float64) error {
	err := checkTransition(gtl.State, MoveEvent)
	if err != nil {
		return err
	}
	return gtl.updateState(ctx, latitude, longitude, Moving)
}

func (gtl *GoToLocation) Near(ctx context.Context, latitude, longitude float64) error {
	err := checkTransition(gtl.State, NearEvent)
	if err != nil {
		return err
	}
	return gtl.updateState(ctx, latitude, longitude, Near)
}

func (gtl *GoToLocation) Arrived(ctx context.Context, latitude, longitude float64) error {
	err := checkTransition(gtl.State, ArrivedEvent)
	if err != nil {
		return err
	}

	return gtl.updateState(ctx, latitude, longitude, Arrived)
}

func (gtl *GoToLocation) Canceled(ctx context.Context, latitude, longitude float64) error {
	err := checkTransition(gtl.State, CancelledEvent)
	if err != nil {
		return err
	}
	return gtl.updateState(ctx, latitude, longitude, Cancelled)
}
