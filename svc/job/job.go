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

func (s *Service) GetAvailableJobs(ctx context.Context, areaKey string, workerID string, radiusMeters int32) (*EstimatedJobs, error) {
	lastLocation := s.tracker.QueryLastLocation(ctx, areaKey, workerID)

	if lastLocation == nil {
		return nil, errors.New("location of worker is not known")
	}

	fromLat := lastLocation.Latitude
	fromLon := lastLocation.Longitude

	jobs, err := job.QueryByPickupDistance(ctx, fromLon, fromLat, "AVAILABLE", radiusMeters, googleMatrixAPILimit)

	if err != nil {
		return nil, fmt.Errorf("query by pickup distance error: %w", err)
	}

	return PickupDistanceEstimatedJobs(lastLocation, jobs)
}
