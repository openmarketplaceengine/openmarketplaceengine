package location

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/openmarketplaceengine/openmarketplaceengine/log"
	v1 "github.com/openmarketplaceengine/openmarketplaceengine/pkg/api/location/proto/v1"
	"github.com/openmarketplaceengine/openmarketplaceengine/pkg/api/location/storage"
	"github.com/openmarketplaceengine/openmarketplaceengine/redis/publisher"
	"github.com/openmarketplaceengine/openmarketplaceengine/redis/subscriber"
)

type Controller struct {
	v1.UnimplementedLocationServiceServer
	store   *storage.Storage
	sub     subscriber.Subscriber
	pub     publisher.Publisher
	areaKey string
}

func New(storeClient *redis.Client, subClient *redis.Client, pubClient *redis.Client, areaKey string) *Controller {
	return &Controller{
		store:   storage.New(storeClient),
		sub:     subscriber.New(subClient),
		pub:     publisher.New(pubClient),
		areaKey: areaKey,
	}
}

func (c *Controller) UpdateLocation(ctx context.Context, request *v1.UpdateLocationRequest) (*v1.UpdateLocationResponse, error) {
	err := c.store.Update(ctx, c.areaKey, storage.Location{
		WorkerID:  request.WorkerId,
		Longitude: request.Longitude,
		Latitude:  request.Latitude,
	})
	if err != nil {
		return nil, err
	}

	channel := fmt.Sprintf("location-%s", request.WorkerId)

	bytes, err := json.Marshal(storage.Location{
		WorkerID:  request.WorkerId,
		Longitude: request.Longitude,
		Latitude:  request.Latitude,
	})
	if err != nil {
		return nil, err
	}
	payload := string(bytes)
	err = c.pub.Publish(ctx, channel, payload)

	if err != nil {
		return nil, err
	}

	return &v1.UpdateLocationResponse{
		WorkerId: request.WorkerId,
	}, nil
}
func (c *Controller) QueryLocation(ctx context.Context, request *v1.QueryLocationRequest) (*v1.QueryLocationResponse, error) {
	l := c.store.QueryLocation(ctx, c.areaKey, request.WorkerId)
	if l != nil {
		return &v1.QueryLocationResponse{
			WorkerId:  l.WorkerID,
			Longitude: l.Longitude,
			Latitude:  l.Latitude,
		}, nil
	}
	return nil, nil
}

func (c *Controller) QueryLocationStreaming(request *v1.QueryLocationStreamingRequest, stream v1.LocationService_QueryLocationStreamingServer) error {

	channel := fmt.Sprintf("location-%s", request.WorkerId)
	rcv := make(chan string)
	c.sub.Subscribe(context.Background(), channel, rcv)

	go func() {

		for m := range rcv {

			var loc storage.Location
			err := json.Unmarshal([]byte(m), &loc)
			if err != nil {
				log.Errorf("skipping due to unmarshal location err=%s", err)
				continue
			}

			if err := stream.Send(&v1.QueryLocationStreamingResponse{
				WorkerId:  loc.WorkerID,
				Longitude: loc.Longitude,
				Latitude:  loc.Latitude,
			}); err != nil {
				log.Errorf("sending location to stream err=%s", err)
			}
		}
	}()

	return nil
}
