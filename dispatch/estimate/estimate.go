package estimate

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/openmarketplaceengine/openmarketplaceengine/dispatch/estimate/google"
	"github.com/openmarketplaceengine/openmarketplaceengine/pkg/validate"
	"googlemaps.github.io/maps"
	"time"
)

type Request struct {
	ID      string
	PickUp  LatLon
	DropOff LatLon
}

type LatLon struct {
	Lat float64
	Lon float64
}

type Location struct {
	Address string
	Lat     float64
	Lon     float64
}

type Eta struct {
	DistanceMeters int
	Duration       time.Duration
}

type Estimate struct {
	ID              string
	From            Location
	PickUp          Location
	DropOff         Location
	ToPickup        Eta
	PickupToDropOff Eta
}

func (e Estimate) MarshalBinary() ([]byte, error) {
	return json.Marshal(e)
}

func (e *Estimate) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &e)
}

func Estimates(ctx context.Context, googleAPIKey string, from LatLon, requests []*Request) ([]*Estimate, error) {
	if len(requests) == 0 {
		return []*Estimate{}, nil
	}

	client, err := maps.NewClient(maps.WithAPIKey(googleAPIKey))
	if err != nil {
		return nil, fmt.Errorf("maps client create error: %w", err)
	}

	err = validateLatLon(requests)
	if err != nil {
		return nil, err
	}

	origins, destinations := toOriginsDestinations(requests)
	originsPlus := append(origins, google.LatLon{Lat: from.Lat, Lon: from.Lon})
	m, err := google.Matrix(ctx, client, google.MatrixPointsInput{
		Origins:      originsPlus,
		Destinations: destinations,
	})

	if err != nil {
		return nil, fmt.Errorf("google matrix API error: %w", err)
	}

	totalOrigins := len(requests) + 1
	if len(m.Rows) != totalOrigins {
		return nil, fmt.Errorf("google matrix API expected to return exactly %v rows because of +1 worker origin, but got %v", totalOrigins, len(m.Rows))
	}

	workerIdx := len(m.Rows) - 1
	workerRow := m.Rows[workerIdx]

	var estimates = make([]*Estimate, len(requests))
	for i, row := range m.Rows {
		if i == workerIdx {
			break
		}
		for j, e := range row.Elements {
			estimates[j] = &Estimate{
				ID: requests[j].ID,
				From: Location{
					Address: m.OriginAddresses[workerIdx],
					Lat:     from.Lat,
					Lon:     from.Lon,
				},
				PickUp: Location{
					Address: m.OriginAddresses[j],
					Lat:     origins[j].Lat,
					Lon:     origins[j].Lon,
				},
				DropOff: Location{
					Address: m.DestinationAddresses[j],
					Lat:     destinations[j].Lat,
					Lon:     destinations[j].Lon,
				},
				ToPickup: Eta{
					DistanceMeters: workerRow.Elements[j].Distance,
					Duration:       workerRow.Elements[j].Duration,
				},
				PickupToDropOff: Eta{
					DistanceMeters: e.Distance,
					Duration:       e.Duration,
				},
			}
		}
	}

	return estimates, nil
}

func toOriginsDestinations(requests []*Request) (origins []google.LatLon, destinations []google.LatLon) {
	for _, d := range requests {
		origins = append(origins, google.LatLon{
			Lat: d.PickUp.Lat,
			Lon: d.PickUp.Lon,
		})
		destinations = append(destinations, google.LatLon{
			Lat: d.DropOff.Lat,
			Lon: d.DropOff.Lon,
		})
	}
	return
}

func validateLatLon(requests []*Request) error {
	v := validate.Validator{}
	for _, r := range requests {
		v.ValidateFloat64(fmt.Sprintf("request(%s).pickUp_lat", r.ID), r.PickUp.Lat).Latitude()
		v.ValidateFloat64(fmt.Sprintf("request(%s).pickUp_lon", r.ID), r.PickUp.Lon).Longitude()
		v.ValidateFloat64(fmt.Sprintf("request(%s).dropOff_lat", r.ID), r.DropOff.Lat).Latitude()
		v.ValidateFloat64(fmt.Sprintf("request(%s).dropOff_lon", r.ID), r.DropOff.Lon).Longitude()
	}
	return v.Error()
}
