package gotolocation

import (
	"context"
	"fmt"
	"github.com/openmarketplaceengine/openmarketplaceengine/core/model/job/step/gotolocation/state"
	"time"
)

type Service interface {
	Moving(ctx context.Context, driverId string, latitude, longitude float64) error
	Near(ctx context.Context, driverId string, latitude, longitude float64) error
	Arrived(ctx context.Context, driverId string, latitude, longitude float64) error
	Canceled(ctx context.Context, driverId string, latitude, longitude float64) error
}

type service struct {
	storage *state.Storage
}

func NewService() Service {
	return &service{
		storage: state.NewStorage(3 * time.Hour),
	}
}

func getKey(driverId string) string {
	return fmt.Sprintf("gotolocation-%s", driverId)
}

func (s *service) Moving(ctx context.Context, driverId string, latitude, longitude float64) error {

	st := state.State{
		DriverID:             driverId,
		DestinationLatitude:  latitude,
		DestinationLongitude: longitude,
		CreatedAt:            time.Now(),
		LastState:            "moving",
	}

	key := getKey(driverId)
	err := s.storage.Store(ctx, key, st)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) Near(ctx context.Context, driverId string, latitude, longitude float64) error {
	key := getKey(driverId)
	st, err := s.storage.Retrieve(ctx, key)
	if err != nil {
		return err
	}

	st.LastModifiedAt = time.Now()
	st.LastModifiedAtLatitude = latitude
	st.LastModifiedAtLongitude = longitude
	st.LastState = "near"

	err = s.storage.Store(ctx, key, st)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) Arrived(ctx context.Context, driverId string, latitude, longitude float64) error {
	key := getKey(driverId)
	st, err := s.storage.Retrieve(ctx, key)
	if err != nil {
		return err
	}

	st.LastModifiedAt = time.Now()
	st.LastModifiedAtLatitude = latitude
	st.LastModifiedAtLongitude = longitude
	st.LastState = "arrived"

	err = s.storage.Store(ctx, key, st)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) Canceled(ctx context.Context, driverId string, latitude, longitude float64) error {
	key := getKey(driverId)
	st, err := s.storage.Retrieve(ctx, key)
	if err != nil {
		return err
	}

	st.LastModifiedAt = time.Now()
	st.LastModifiedAtLatitude = latitude
	st.LastModifiedAtLongitude = longitude
	st.LastState = "cancelled"

	err = s.storage.Store(ctx, key, st)
	if err != nil {
		return err
	}
	return nil
}
