package location

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
	locationV1beta1 "github.com/openmarketplaceengine/openmarketplaceengine/internal/omeapi/location/v1beta1"
	"github.com/openmarketplaceengine/openmarketplaceengine/internal/service/location/storage"
	"github.com/openmarketplaceengine/openmarketplaceengine/internal/service/tollgate"
	"github.com/openmarketplaceengine/openmarketplaceengine/log"
	"github.com/openmarketplaceengine/openmarketplaceengine/redis/publisher"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Controller struct {
	locationV1beta1.UnimplementedLocationServiceServer
	store    *storage.Storage
	pub      publisher.Publisher
	areaKey  string
	detector tollgate.Detector
}

func New(storeClient *redis.Client, pubClient *redis.Client, areaKey string, detector tollgate.Detector) *Controller {
	return &Controller{
		store:    storage.New(storeClient),
		pub:      publisher.New(pubClient),
		areaKey:  areaKey,
		detector: detector,
	}
}

func (c *Controller) UpdateLocation(ctx context.Context, request *locationV1beta1.UpdateLocationRequest) (*locationV1beta1.UpdateLocationResponse, error) {
	lastLocation := c.store.LastLocation(ctx, c.areaKey, request.WorkerId)

	err := c.store.Update(ctx, c.areaKey, &storage.Location{
		WorkerID:  request.WorkerId,
		Longitude: request.Longitude,
		Latitude:  request.Latitude,
	})
	if err != nil {
		return nil, err
	}

	c.publishLocation(ctx, request.WorkerId, request.Longitude, request.Latitude)

	var tollgateCrossing *locationV1beta1.TollgateCrossing

	if lastLocation != nil {
		from := &tollgate.Location{
			Lon: lastLocation.Longitude,
			Lat: lastLocation.Latitude,
		}
		to := &tollgate.Location{
			Lon: request.Longitude,
			Lat: request.Latitude}
		movement := &tollgate.Movement{
			SubjectID: request.WorkerId,
			From:      from,
			To:        to,
		}

		crossing, err := c.detector.DetectCrossing(ctx, movement)
		if err != nil {
			return nil, err
		}
		if crossing != nil {
			c.publishTollgateCrossing(ctx, crossing)
			tollgateCrossing = &locationV1beta1.TollgateCrossing{
				TollgateId: crossing.TollgateID,
				Longitude:  crossing.Location.Lon,
				Latitude:   crossing.Location.Lat,
				Direction:  string(crossing.Direction),
			}
		}
	}

	return &locationV1beta1.UpdateLocationResponse{
		WorkerId:         request.WorkerId,
		TollgateCrossing: tollgateCrossing,
		UpdateTime:       request.UpdateTime,
	}, nil
}

func (c *Controller) publishLocation(ctx context.Context, workerID string, longitude float64, latitude float64) {
	channel := locationChannel(workerID)

	bytes, err := json.Marshal(storage.Location{
		WorkerID:  workerID,
		Longitude: longitude,
		Latitude:  latitude,
	})
	if err != nil {
		log.Errorf("location marshal error: %q", err)
		return
	}
	payload := string(bytes)
	err = c.pub.Publish(ctx, channel, payload)

	if err != nil {
		log.Errorf("location publish error: %q", err)
		return
	}
}

func locationChannel(workerID string) string {
	return fmt.Sprintf("channel-location-%s", workerID)
}

func (c *Controller) publishTollgateCrossing(ctx context.Context, crossing *tollgate.Crossing) {
	channel := crossingChannel(crossing.TollgateID)

	bytes, err := json.Marshal(crossing)
	if err != nil {
		log.Errorf("crossing marshal error: %q", err)
		return
	}
	payload := string(bytes)
	err = c.pub.Publish(ctx, channel, payload)

	if err != nil {
		log.Errorf("crossing publish error: %q", err)
		return
	}
}

func crossingChannel(tollgateID string) string {
	return fmt.Sprintf("channel:crossing:%s", tollgateID)
}

func (c *Controller) QueryLocation(ctx context.Context, request *locationV1beta1.QueryLocationRequest) (*locationV1beta1.QueryLocationResponse, error) {
	l := c.store.LastLocation(ctx, c.areaKey, request.WorkerId)
	if l != nil {
		return &locationV1beta1.QueryLocationResponse{
			WorkerId:  l.WorkerID,
			Longitude: l.Longitude,
			Latitude:  l.Latitude,
			LastSeenTime: &timestamppb.Timestamp{
				Seconds: l.LastSeenTime.Unix(),
				Nanos:   0,
			},
		}, nil
	}
	return nil, fmt.Errorf("location not found for WorkerId=%s", request.WorkerId)
}
