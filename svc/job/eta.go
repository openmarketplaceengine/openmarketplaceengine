package job

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/openmarketplaceengine/geoservices"
	"github.com/openmarketplaceengine/geoservices/distance"
	"github.com/openmarketplaceengine/geoservices/distance/google"
	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
	"github.com/openmarketplaceengine/openmarketplaceengine/dom/job"
	"github.com/openmarketplaceengine/openmarketplaceengine/svc/location"
	"googlemaps.github.io/maps"
)

type EstimatedJobs struct {
	Jobs          []*EstimatedJob
	OriginAddress string
}

type EstimatedJob struct {
	Job      *job.AvailableJob
	Address  string
	Distance int
	Duration time.Duration
}

func PickupDistanceEstimatedJobs(from *location.LastLocation, to []*job.AvailableJob) (*EstimatedJobs, error) {
	client, err := maps.NewClient(maps.WithAPIKey(cfg.Server.GoogleAPIKey))
	if err != nil {
		return nil, fmt.Errorf("maps client create error: %w", err)
	}
	m, err := google.Matrix(context.Background(), client, distance.MatrixPointsInput{
		Origins:      []geoservices.LatLng{{Lat: from.Latitude, Lng: from.Longitude}},
		Destinations: transform(to),
	})

	if err != nil {
		return nil, fmt.Errorf("google matrix API error: %w", err)
	}
	if len(m.Rows) > 1 {
		return nil, fmt.Errorf("google matrix API expected to return exactly 1 row because of 1 origin, but got %v", len(m.Rows))
	}
	row := m.Rows[0]

	var jobs = make([]*EstimatedJob, len(to))
	for i, e := range row.Elements {
		jobs[i] = &EstimatedJob{
			Job:      to[i],
			Address:  m.DestinationAddresses[i],
			Distance: e.Distance,
			Duration: e.Duration,
		}
	}

	sort.SliceStable(jobs, func(i, j int) bool {
		return jobs[i].Duration < jobs[j].Duration
	})
	return &EstimatedJobs{
		Jobs:          jobs,
		OriginAddress: m.OriginAddresses[0],
	}, nil
}

func transform(to []*job.AvailableJob) (ll []geoservices.LatLng) {
	for _, j := range to {
		ll = append(ll, geoservices.LatLng{
			Lat: j.PickupLat,
			Lng: j.PickupLon,
		})
	}
	return
}
