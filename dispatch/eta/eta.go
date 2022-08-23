package eta

import (
	"context"
	"fmt"

	"github.com/openmarketplaceengine/openmarketplaceengine/dispatch/google"
	"github.com/openmarketplaceengine/openmarketplaceengine/dispatch/job"
	"github.com/openmarketplaceengine/openmarketplaceengine/pkg/validate"
	"googlemaps.github.io/maps"
)

func EstimateJobs(ctx context.Context, googleAPIKey string, latLon job.LatLon, jobs []*job.Job) ([]*job.EstimatedJob, error) {
	if len(jobs) == 0 {
		return []*job.EstimatedJob{}, nil
	}

	client, err := maps.NewClient(maps.WithAPIKey(googleAPIKey))
	if err != nil {
		return nil, fmt.Errorf("maps client create error: %w", err)
	}

	err = validateLatLon(jobs)
	if err != nil {
		return nil, err
	}

	origins, destinations := jobsToLocations(jobs)
	originsPlus := append(origins, latLon)
	m, err := google.Matrix(ctx, client, google.MatrixPointsInput{
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

	var eJobs = make([]*job.EstimatedJob, len(jobs))
	for i, row := range m.Rows {
		if i == workerIdx {
			break
		}
		for j, e := range row.Elements {
			eJobs[j] = &job.EstimatedJob{
				ID: jobs[j].ID,
				WorkerLocation: job.Location{
					Address: m.OriginAddresses[workerIdx],
					Lat:     latLon.Lat,
					Lon:     latLon.Lon,
				},
				Pickup: job.Location{
					Address: m.OriginAddresses[j],
					Lat:     origins[j].Lat,
					Lon:     origins[j].Lon,
				},
				DropOff: job.Location{
					Address: m.DestinationAddresses[j],
					Lat:     destinations[j].Lat,
					Lon:     destinations[j].Lon,
				},
				ToPickup: job.Estimate{
					DistanceMeters: workerRow.Elements[j].Distance,
					Duration:       workerRow.Elements[j].Duration,
				},
				PickupToDropOff: job.Estimate{
					DistanceMeters: e.Distance,
					Duration:       e.Duration,
				},
			}
		}
	}

	return eJobs, nil
}

func jobsToLocations(jobs []*job.Job) (origins []job.LatLon, destinations []job.LatLon) {
	for _, j := range jobs {
		origins = append(origins, job.LatLon{
			Lat: j.Pickup.Lat,
			Lon: j.Pickup.Lon,
		})
		destinations = append(destinations, job.LatLon{
			Lat: j.DropOff.Lat,
			Lon: j.DropOff.Lon,
		})
	}
	return
}

func validateLatLon(jobs []*job.Job) error {
	v := validate.Validator{}
	for _, j := range jobs {
		v.ValidateFloat64(fmt.Sprintf("job(%s).pickup_lat", j.ID), j.Pickup.Lat).Latitude()
		v.ValidateFloat64(fmt.Sprintf("job(%s).pickup_lon", j.ID), j.Pickup.Lon).Longitude()
		v.ValidateFloat64(fmt.Sprintf("job(%s).dropoff_lat", j.ID), j.DropOff.Lat).Latitude()
		v.ValidateFloat64(fmt.Sprintf("job(%s).dropoff_lon", j.ID), j.DropOff.Lon).Longitude()
	}
	return v.Error()
}
