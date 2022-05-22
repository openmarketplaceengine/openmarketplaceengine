package location

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/openmarketplaceengine/openmarketplaceengine/dao"
	"github.com/openmarketplaceengine/openmarketplaceengine/log"
	"github.com/openmarketplaceengine/openmarketplaceengine/pkg/detector"
)

type Handler func(ctx context.Context, l *Location) error
type CrossingHandler func(ctx context.Context, crossing *detector.Crossing) error

type Tracker struct {
	storage  *Storage
	detector *detector.Detector
}

func NewTracker(storage *Storage, detector *detector.Detector) (*Tracker, error) {
	return &Tracker{
		storage:  storage,
		detector: detector,
	}, nil
}

func (t *Tracker) TrackLocation(ctx context.Context, areaKey string, workerID string, lon float64, lat float64) (*detector.Crossing, error) {
	lastLocation := t.storage.LastLocation(ctx, areaKey, workerID)

	err := t.storage.Update(ctx, areaKey, &Location{
		WorkerID:  workerID,
		Longitude: lon,
		Latitude:  lat,
	}, time.Now())
	if err != nil {
		return nil, fmt.Errorf("update location error: %w", err)
	}

	tollgateCrossing, err := t.detectCrossing(ctx, lastLocation, workerID, lon, lat, PublishCrossing, PersistCrossing)
	if err != nil {
		return nil, fmt.Errorf("detect tollgate crossing error: %w", err)
	}
	return tollgateCrossing, nil
}

func (t *Tracker) QueryLastLocation(ctx context.Context, areaKey string, workerID string) *LastLocation {
	l := t.storage.LastLocation(ctx, areaKey, workerID)
	if l != nil {
		return l
	}
	return nil
}

func (t *Tracker) detectCrossing(ctx context.Context, lastLocation *LastLocation, workerID string, lon float64, lat float64, handlers ...CrossingHandler) (*detector.Crossing, error) {
	if lastLocation != nil {
		from := &detector.Location{
			Lon: lastLocation.Longitude,
			Lat: lastLocation.Latitude,
		}
		to := &detector.Location{
			Lon: lon,
			Lat: lat}

		movement := &detector.Movement{
			From: from,
			To:   to,
		}

		detected, err := t.detector.DetectCrossing(ctx, workerID, movement)
		if err != nil {
			return nil, fmt.Errorf("detect crossing error: %s", err)
		}
		if detected != nil {
			for _, handler := range handlers {
				err := handler(ctx, detected)
				if err != nil {
					return nil, fmt.Errorf("crossing handler error: %s", err)
				}
			}
			return detected, nil
		}
	}
	return nil, nil
}

func (t *Tracker) publishLocation(ctx context.Context, workerID string, lon float64, lat float64) {
	channel := locationChannel(workerID)

	bytes, err := json.Marshal(Location{
		WorkerID:  workerID,
		Longitude: lon,
		Latitude:  lat,
	})
	if err != nil {
		log.Errorf("location marshal error: %q", err)
		return
	}
	payload := string(bytes)
	pub := dao.Reds.PubSubClient
	err = pub.Publish(ctx, channel, payload).Err()

	if err != nil {
		log.Errorf("location publish error: %q", err)
		return
	}
}

func locationChannel(workerID string) string {
	return fmt.Sprintf("channel-location-%s", workerID)
}
