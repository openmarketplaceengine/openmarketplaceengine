package job

import (
	"context"
	"errors"
	"fmt"

	"github.com/openmarketplaceengine/openmarketplaceengine/dom/job"
	"github.com/openmarketplaceengine/openmarketplaceengine/svc/location"
)

const googleMatrixAPILimit = int32(25)

type Service struct {
	tracker *location.Tracker
}

func NewService(tracker *location.Tracker) *Service {
	return &Service{
		tracker: tracker,
	}
}

func (s *Service) EstimatedJobs(ctx context.Context, areaKey string, workerID string, radiusMeters int32) ([]*EstimatedJob, error) {
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

	return estimatedJobs(ctx, workerLocation, jobs)
}
