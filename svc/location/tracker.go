package location

import (
	"context"
	"fmt"

	"github.com/openmarketplaceengine/openmarketplaceengine/pkg/detector"
)

type Handler func(ctx context.Context, areaKey string, l *Location) error

type Tracker struct {
	storage          *Storage
	locationHandlers []Handler
	crossingHandlers []detector.Handler
	detector         *detector.Detector
}

func NewTracker(s *Storage, d *detector.Detector) *Tracker {
	return &Tracker{
		storage:          s,
		locationHandlers: []Handler{s.RedisUpdateHandler(), persistLocation, publishLocation},
		crossingHandlers: []detector.Handler{persistCrossing, publishCrossing},
		detector:         d,
	}
}

func (t *Tracker) TrackLocation(ctx context.Context, areaKey string, workerID string, longitude float64, latitude float64) (*detector.Crossing, error) {
	workerLocation := t.storage.WorkerLocation(ctx, areaKey, workerID)

	l := &Location{
		WorkerID:  workerID,
		Longitude: longitude,
		Latitude:  latitude,
	}

	for _, handler := range t.locationHandlers {
		err := handler(ctx, areaKey, l)
		if err != nil {
			return nil, fmt.Errorf("handler error: %s", err)
		}
	}

	if workerLocation != nil {
		from := &detector.Location{
			Longitude: workerLocation.Longitude,
			Latitude:  workerLocation.Latitude,
		}
		to := &detector.Location{
			Longitude: longitude,
			Latitude:  latitude}

		movement := &detector.Movement{
			From: from,
			To:   to,
		}

		detected, err := t.detector.DetectCrossing(ctx, workerID, movement, t.crossingHandlers...)
		if err != nil {
			return nil, fmt.Errorf("detect crossing error: %s", err)
		}
		return detected, nil
	}

	return nil, nil
}

func (t *Tracker) GetLocation(ctx context.Context, areaKey string, workerID string) *WorkerLocation {
	l := t.storage.WorkerLocation(ctx, areaKey, workerID)
	if l != nil {
		return l
	}
	return nil
}
