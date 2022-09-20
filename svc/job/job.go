package job

import (
	"context"
	"errors"
	"fmt"

	"github.com/openmarketplaceengine/openmarketplaceengine/dom/job"
	"github.com/openmarketplaceengine/openmarketplaceengine/svc/location"
)

// Limit of initial available jobs query to accommodate google usage constraints.
// https://developers.google.com/maps/documentation/distance-matrix/usage-and-billing#other-usage-limits
const googleMatrixAPILimit = int32(25)

type Service struct {
	tracker *location.Tracker
}

func NewService(tracker *location.Tracker) *Service {
	return &Service{
		tracker: tracker,
	}
}

func (s *Service) GetEstimatedJobs(ctx context.Context, areaKey string, workerID string, radiusMeters int32) ([]*Job, error) {
	workerLocation := s.tracker.QueryLastLocation(ctx, areaKey, workerID)

	if workerLocation == nil {
		return nil, errors.New("location of worker is not known")
	}

	fromLat := workerLocation.Latitude
	fromLon := workerLocation.Longitude

	jobs, err := job.QueryByPickupDistance(ctx, fromLon, fromLat, "AVAILABLE", radiusMeters, googleMatrixAPILimit)

	if err != nil {
		return nil, fmt.Errorf("query by pickup distance error: %w", err)
	}

	return estimateJobs(ctx, workerLocation, jobs)
}
