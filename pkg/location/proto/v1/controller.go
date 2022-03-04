package v1

import (
	"context"
	"io"
)

type Controller struct {
	UnimplementedLocationServiceServer
}

func (s Controller) UpdateLocation(ctx context.Context, request *UpdateLocationRequest) (*UpdateLocationResponse, error) {
	return &UpdateLocationResponse{WorkerId: request.WorkerId}, nil
}
func (s Controller) QueryLocation(ctx context.Context, request *QueryLocationRequest) (*QueryLocationResponse, error) {
	return &QueryLocationResponse{
		WorkerId:  request.WorkerId,
		Longitude: 0,
		Latitude:  0,
		Distance:  0,
		LastSeen:  nil,
	}, nil
}

func (s Controller) UpdateLocationStreaming(stream LocationService_UpdateLocationStreamingServer) error {
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

func (s Controller) QueryLocationStreaming(request *QueryLocationRequest, stream LocationService_QueryLocationStreamingServer) error {
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
