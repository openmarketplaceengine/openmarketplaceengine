package job

import (
	"context"
	"fmt"
	"github.com/openmarketplaceengine/openmarketplaceengine/dom/job"
	"github.com/openmarketplaceengine/openmarketplaceengine/svc/location"
)

type Service struct {
	tracker *location.Tracker
}

func NewService(tracker *location.Tracker) *Service {
	return &Service{
		tracker: tracker,
	}
}

func (s *Service) GetAvailableJobs(ctx context.Context, areaKey string, workerID string, rangeLimit float64, limit int) ([]*job.AvailableJob, error) {
	lastLocation := s.tracker.QueryLastLocation(ctx, areaKey, workerID)

	if lastLocation == nil {
		return nil, fmt.Errorf("location of worker is not known")
	}

	fromLat := lastLocation.Latitude
	fromLon := lastLocation.Longitude

	jobs, err := job.QueryByPickupDistance(ctx, fromLat, fromLon, "AVAILABLE", rangeLimit, job.Mile, limit)

	if err != nil {
		return nil, fmt.Errorf("query by pickup distance error: %w", err)
	}

	return jobs, nil
}
