package v1

import (
	"context"
	"io"

	"github.com/go-redis/redis/v8"
	"github.com/openmarketplaceengine/openmarketplaceengine/pkg/location"
	"github.com/openmarketplaceengine/openmarketplaceengine/pkg/location/storage"
)

type Controller struct {
	UnimplementedLocationServiceServer
	storage *storage.Storage
	areaKey string
}

func New(client *redis.Client, areaKey string) *Controller {
	return &Controller{
		storage: storage.New(client),
		areaKey: areaKey,
	}
}

func (c *Controller) UpdateLocation(ctx context.Context, request *UpdateLocationRequest) (*UpdateLocationResponse, error) {
	err := c.storage.Update(ctx, c.areaKey, location.Location{
		WorkerID:  request.WorkerId,
		Longitude: request.Longitude,
		Latitude:  request.Latitude,
	})
	if err != nil {
		return nil, err
	}
	return &UpdateLocationResponse{
		WorkerId: request.WorkerId,
	}, nil
}
func (c *Controller) QueryLocation(ctx context.Context, request *QueryLocationRequest) (*QueryLocationResponse, error) {
	l := c.storage.QueryLocation(ctx, c.areaKey, request.WorkerId)
	if l != nil {
		return &QueryLocationResponse{
			WorkerId:  l.WorkerID,
			Longitude: l.Longitude,
			Latitude:  l.Latitude,
		}, nil
	}
	return nil, nil
}

func (c *Controller) UpdateLocationStreaming(stream LocationService_UpdateLocationStreamingServer) error {
	for {
		_, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&UpdateLocationResponse{})
		}
		if err != nil {
			return err
		}
	}
}

func (c *Controller) QueryLocationStreaming(request *QueryLocationRequest, stream LocationService_QueryLocationStreamingServer) error {
	if err := stream.Send(&QueryLocationResponse{
		WorkerId:  request.WorkerId,
		Longitude: 0,
		Latitude:  0,
		Distance:  0,
		LastSeen:  nil,
	}); err != nil {
		return err
	}
	return nil
}
