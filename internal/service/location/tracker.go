package location

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
	"github.com/openmarketplaceengine/openmarketplaceengine/internal/service/location/storage"
	"github.com/openmarketplaceengine/openmarketplaceengine/internal/service/tollgate"
	"github.com/openmarketplaceengine/openmarketplaceengine/internal/service/tollgate/crossing"
	"github.com/openmarketplaceengine/openmarketplaceengine/internal/service/tollgate/detector"
	"github.com/openmarketplaceengine/openmarketplaceengine/log"
)

const areaKey = "global"

type Tracker struct {
	storage      *storage.Storage
	pubSubClient *redis.Client
	detector     *detector.Detector
}

func NewTracker(storeClient *redis.Client, pubSubClient *redis.Client) (*Tracker, error) {
	tollgates, err := tollgate.QueryAll(cfg.Context(), 100)
	if err != nil {
		return nil, err
	}

	d := detector.NewDetector(transformTollgates(tollgates), storeClient)

	return &Tracker{
		storage:      storage.New(storeClient),
		pubSubClient: pubSubClient,
		detector:     d,
	}, nil
}

func transformTollgates(tollgates []*tollgate.Tollgate) (result []*detector.Tollgate) {
	for _, t := range tollgates {
		var line *detector.Line
		var bBoxes []*detector.BBox
		var bBoxesRequired int32

		if t.GateLine != nil {
			line = t.GateLine.Line
		}

		if t.BBoxes != nil {
			bBoxes = t.BBoxes.BBoxes
			bBoxesRequired = t.BBoxes.Required
		}

		result = append(result, &detector.Tollgate{
			ID:             t.ID,
			Line:           line,
			BBoxes:         bBoxes,
			BBoxesRequired: bBoxesRequired,
		})
	}
	return
}

func (s *Tracker) TrackLocation(ctx context.Context, areaKey string, workerID string, lon float64, lat float64) (*crossing.TollgateCrossing, error) {
	lastLocation := s.storage.LastLocation(ctx, areaKey, workerID)

	err := s.storage.Update(ctx, areaKey, &storage.Location{
		WorkerID:  workerID,
		Longitude: lon,
		Latitude:  lat,
	}, time.Now())
	if err != nil {
		return nil, fmt.Errorf("update location error: %w", err)
	}

	s.publishLocation(ctx, workerID, lon, lat)

	tollgateCrossing, err := s.detectTollgateCrossing(ctx, lastLocation, workerID, lon, lat)
	if err != nil {
		return nil, fmt.Errorf("detect tollgate crossing error: %w", err)
	}
	return tollgateCrossing, nil
}

func (s *Tracker) QueryLastLocation(ctx context.Context, areaKey string, workerID string) *storage.LastLocation {
	l := s.storage.LastLocation(ctx, areaKey, workerID)
	if l != nil {
		return l
	}
	return nil
}

func (s *Tracker) detectTollgateCrossing(ctx context.Context, lastLocation *storage.LastLocation, workerID string, lon float64, lat float64) (*crossing.TollgateCrossing, error) {
	if lastLocation != nil {
		from := &detector.Location{
			Lon: lastLocation.Longitude,
			Lat: lastLocation.Latitude,
		}
		to := &detector.Location{
			Lon: lon,
			Lat: lat}
		movement := &detector.Movement{
			SubjectID: workerID,
			From:      from,
			To:        to,
		}

		detected, err := s.detector.DetectCrossing(ctx, movement)
		if err != nil {
			return nil, fmt.Errorf("detect crossing error: %s", err)
		}
		if detected != nil {
			tollgateCrossing := crossing.NewTollgateCrossing(detected.TollgateID, movement.SubjectID, detected)
			err := tollgateCrossing.Insert(ctx)
			if err != nil {
				return nil, fmt.Errorf("crossing insert error: %s", err)
			}
			s.publishTollgateCrossing(ctx, detected)
			return tollgateCrossing, nil
		}
	}
	return nil, nil
}

func (s *Tracker) publishLocation(ctx context.Context, workerID string, lon float64, lat float64) {
	channel := locationChannel(workerID)

	bytes, err := json.Marshal(storage.Location{
		WorkerID:  workerID,
		Longitude: lon,
		Latitude:  lat,
	})
	if err != nil {
		log.Errorf("location marshal error: %q", err)
		return
	}
	payload := string(bytes)
	err = s.pubSubClient.Publish(ctx, channel, payload).Err()

	if err != nil {
		log.Errorf("location publish error: %q", err)
		return
	}
}

func locationChannel(workerID string) string {
	return fmt.Sprintf("channel-location-%s", workerID)
}

func (s *Tracker) publishTollgateCrossing(ctx context.Context, crossing *detector.Crossing) {
	channel := crossingChannel(crossing.TollgateID)

	bytes, err := json.Marshal(crossing)
	if err != nil {
		log.Errorf("crossing marshal error: %q", err)
		return
	}
	payload := string(bytes)
	err = s.pubSubClient.Publish(ctx, channel, payload).Err()

	if err != nil {
		log.Errorf("crossing publish error: %q", err)
		return
	}
}

func crossingChannel(tollgateID string) string {
	return fmt.Sprintf("channel:crossing:%s", tollgateID)
}
