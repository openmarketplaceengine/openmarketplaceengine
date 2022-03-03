package protos

import (
	"context"
	v1 "github.com/openmarketplaceengine/openmarketplaceengine/core/model/location/protos/v1"
	"io"
)

type Server struct {
	v1.UnimplementedLocationServiceServer
}

func (s Server) UpdateLocation(ctx context.Context, request *v1.UpdateLocationRequest) (*v1.UpdateLocationResponse, error) {
	return &v1.UpdateLocationResponse{WorkerId: request.WorkerId}, nil
}
func (s Server) QueryLocation(ctx context.Context, request *v1.QueryLocationRequest) (*v1.QueryLocationResponse, error) {
	return &v1.QueryLocationResponse{
		WorkerId:  request.WorkerId,
		Longitude: 0,
		Latitude:  0,
		Distance:  0,
		LastSeen:  nil,
	}, nil
}

func (s Server) UpdateLocationStreaming(stream v1.LocationService_UpdateLocationStreamingServer) error {
	for {
		_, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&v1.UpdateLocationResponse{})
		}
		if err != nil {
			return err
		}
	}
}

func (s Server) QueryLocationStreaming(request *v1.QueryLocationRequest, stream v1.LocationService_QueryLocationStreamingServer) error {
	if err := stream.Send(&v1.QueryLocationResponse{
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
