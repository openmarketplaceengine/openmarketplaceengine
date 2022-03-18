package location

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
	locationV1beta1 "github.com/openmarketplaceengine/openmarketplaceengine/internal/omeapi/location/v1beta1"
	"github.com/openmarketplaceengine/openmarketplaceengine/internal/service/location/storage"
	"github.com/openmarketplaceengine/openmarketplaceengine/internal/service/tollgate"
	"github.com/openmarketplaceengine/openmarketplaceengine/redis/publisher"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Controller struct {
	locationV1beta1.UnimplementedLocationServiceServer
	store            *storage.Storage
	pub              publisher.Publisher
	areaKey          string
	tollgateDetector *tollgate.Detector
}

func New(storeClient *redis.Client, pubClient *redis.Client, areaKey string) *Controller {
	return &Controller{
		store:            storage.New(storeClient),
		pub:              publisher.New(pubClient),
		areaKey:          areaKey,
		tollgateDetector: tollgate.New(),
	}
}

func (c *Controller) UpdateLocation(ctx context.Context, request *locationV1beta1.UpdateLocationRequest) (*locationV1beta1.UpdateLocationResponse, error) {
	err := c.store.Update(ctx, c.areaKey, &storage.Location{
		WorkerID:  request.WorkerId,
		Longitude: request.Longitude,
		Latitude:  request.Latitude,
	})
	if err != nil {
		return nil, err
	}

	err = c.publishLocation(ctx, request.WorkerId, request.Longitude, request.Latitude)
	if err != nil {
		return nil, err
	}

	//TODO need to keep previous location
	from := tollgate.LocationXY{
		LongitudeX: request.Longitude,
		LatitudeY:  request.Latitude,
	}
	movement := &tollgate.Movement{SubjectID: request.WorkerId,
		From: from, To: tollgate.LocationXY{
			LongitudeX: request.Longitude,
			LatitudeY:  request.Latitude}}
	crossing := c.tollgateDetector.Detect(movement)
	if crossing != nil {
		err = c.publishTollgateCrossing(ctx, crossing)
		if err != nil {
			return nil, err
		}
	}

	return &locationV1beta1.UpdateLocationResponse{
		WorkerId: request.WorkerId,
	}, nil
}

func (c *Controller) publishLocation(ctx context.Context, workerID string, longitude float64, latitude float64) error {
	channel := fmt.Sprintf("location-%s", workerID)

	bytes, err := json.Marshal(storage.Location{
		WorkerID:  workerID,
		Longitude: longitude,
		Latitude:  latitude,
	})
	if err != nil {
		return err
	}
	payload := string(bytes)
	err = c.pub.Publish(ctx, channel, payload)

	if err != nil {
		return err
	}
	return nil
}

func (c *Controller) publishTollgateCrossing(ctx context.Context, crossing *tollgate.Crossing) error {
	channel := fmt.Sprintf("tollgate-x-%s", crossing.TollgateID)

	bytes, err := json.Marshal(crossing)
	if err != nil {
		return err
	}
	payload := string(bytes)
	err = c.pub.Publish(ctx, channel, payload)

	if err != nil {
		return err
	}
	return nil
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
