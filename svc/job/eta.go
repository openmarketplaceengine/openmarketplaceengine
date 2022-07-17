package job

import (
	"context"
	"fmt"
	"time"

	"github.com/openmarketplaceengine/geoservices"
	"github.com/openmarketplaceengine/geoservices/distance"
	"github.com/openmarketplaceengine/geoservices/distance/google"
	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
	"github.com/openmarketplaceengine/openmarketplaceengine/dao"
	"github.com/openmarketplaceengine/openmarketplaceengine/dom/job"
	"github.com/openmarketplaceengine/openmarketplaceengine/pkg/validate"
	"github.com/openmarketplaceengine/openmarketplaceengine/svc/location"
	"googlemaps.github.io/maps"
)

type Estimate struct {
	DistanceMeters int
	Duration       time.Duration
}

type EstimatedJob struct {
	ID              dao.SUID
	WorkerToPickup  Estimate
	PickupToDropOff Estimate
	WorkerLocation  Location
	PickupLocation  Location
	DropOffLocation Location
}

type Location struct {
	Address string
	Lat     float64
	Lng     float64
}

func EstimateJobs(ctx context.Context, workerLocation *location.WorkerLocation, jobs []*job.AvailableJob) ([]*EstimatedJob, error) {
	if len(jobs) == 0 {
		return []*EstimatedJob{}, nil
	}

	client, err := maps.NewClient(maps.WithAPIKey(cfg.Server.GoogleAPIKey))
	if err != nil {
		return nil, fmt.Errorf("maps client create error: %w", err)
	}

	err = validateLatLng(jobs)
	if err != nil {
		return nil, err
	}

	origins, destinations := transform(jobs)
	originsPlus := append(origins, geoservices.LatLng{Lat: workerLocation.Latitude, Lng: workerLocation.Longitude})
	m, err := google.Matrix(ctx, client, distance.MatrixPointsInput{
		Origins:      originsPlus,
		Destinations: destinations,
	})

	if err != nil {
		return nil, fmt.Errorf("google matrix API error: %w", err)
	}

	totalOrigins := len(jobs) + 1
	if len(m.Rows) != totalOrigins {
		return nil, fmt.Errorf("google matrix API expected to return exactly %v rows because of +1 worker origin, but got %v", totalOrigins, len(m.Rows))
	}

	workerIdx := len(m.Rows) - 1
	workerRow := m.Rows[workerIdx]

	var eJobs = make([]*EstimatedJob, len(jobs))
	for i, row := range m.Rows {
		if i == workerIdx {
			break
		}
		for j, e := range row.Elements {
			eJobs[j] = &EstimatedJob{
				ID: jobs[j].ID,
				WorkerLocation: Location{
					Address: m.OriginAddresses[workerIdx],
					Lat:     workerLocation.Latitude,
					Lng:     workerLocation.Longitude,
				},
				PickupLocation: Location{
					Address: m.OriginAddresses[j],
					Lat:     origins[j].Lat,
					Lng:     origins[j].Lng,
				},
				DropOffLocation: Location{
					Address: m.DestinationAddresses[j],
					Lat:     destinations[j].Lat,
					Lng:     destinations[j].Lng,
				},
				WorkerToPickup: Estimate{
					DistanceMeters: workerRow.Elements[j].Distance,
					Duration:       workerRow.Elements[j].Duration,
				},
				PickupToDropOff: Estimate{
					DistanceMeters: e.Distance,
					Duration:       e.Duration,
				},
			}
		}
	}

	return eJobs, nil
}

func transform(jobs []*job.AvailableJob) (origins []geoservices.LatLng, destinations []geoservices.LatLng) {
	for _, j := range jobs {
		origins = append(origins, geoservices.LatLng{
			Lat: j.PickupLat,
			Lng: j.PickupLon,
		})
		destinations = append(destinations, geoservices.LatLng{
			Lat: j.DropoffLat,
			Lng: j.DropoffLon,
		})
	}
	return
}

func validateLatLng(jobs []*job.AvailableJob) error {
	v := validate.Validator{}
	for _, j := range jobs {
		v.ValidateFloat64(fmt.Sprintf("job(%s).pickup_lat", j.ID), j.PickupLat).Latitude()
		v.ValidateFloat64(fmt.Sprintf("job(%s).pickup_lon", j.ID), j.PickupLon).Longitude()
		v.ValidateFloat64(fmt.Sprintf("job(%s).dropoff_lat", j.ID), j.DropoffLat).Latitude()
		v.ValidateFloat64(fmt.Sprintf("job(%s).dropoff_lon", j.ID), j.DropoffLon).Longitude()
	}
	return v.Error()
}
