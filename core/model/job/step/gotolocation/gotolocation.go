package gotolocation

import (
	"context"
	"fmt"
	"time"
)

type State int

const (
	New State = iota
	Moving
	Near
	Arrived
	Cancelled
)

type GoToLocation struct {
	DriverID                string
	DestinationLatitude     float64
	DestinationLongitude    float64
	LastModifiedAt          time.Time
	LastModifiedAtLatitude  float64
	LastModifiedAtLongitude float64
	State                   State
}

func key(driverId string) string {
	return fmt.Sprintf("gotolocation-%s", driverId)
}

func NewGoToLocation(ctx context.Context, driverID string, latitude, longitude float64) (gtl *GoToLocation, err error) {
	gtl = &GoToLocation{
		DriverID:             driverID,
		DestinationLatitude:  latitude,
		DestinationLongitude: longitude,
		LastModifiedAt:       time.Now(),
		State:                New,
	}

	err = storage.Store(ctx, key(gtl.DriverID), *gtl)
	if err != nil {
		return
	}
	return
}

func (gtl *GoToLocation) updateState(ctx context.Context, latitude float64, longitude float64, state State) error {
	gtl.LastModifiedAtLatitude = latitude
	gtl.LastModifiedAtLongitude = longitude
	gtl.LastModifiedAt = time.Now()
	gtl.State = state

	err := storage.Store(ctx, key(gtl.DriverID), *gtl)
	if err != nil {
		return err
	}

	return nil
}

func (gtl *GoToLocation) Moving(ctx context.Context, latitude, longitude float64) error {
	return gtl.updateState(ctx, latitude, longitude, Moving)
}

func (gtl *GoToLocation) Near(ctx context.Context, latitude, longitude float64) error {
	return gtl.updateState(ctx, latitude, longitude, Near)
}

func (gtl *GoToLocation) Arrived(ctx context.Context, latitude, longitude float64) error {
	return gtl.updateState(ctx, latitude, longitude, Arrived)
}

func (gtl *GoToLocation) Canceled(ctx context.Context, latitude, longitude float64) error {
	return gtl.updateState(ctx, latitude, longitude, Cancelled)
}
