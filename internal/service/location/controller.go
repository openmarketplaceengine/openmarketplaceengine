package location

import (
	"context"
	"fmt"

	"github.com/openmarketplaceengine/openmarketplaceengine/dao"
	"github.com/openmarketplaceengine/openmarketplaceengine/internal/service/tollgate/crossing"
	"github.com/openmarketplaceengine/openmarketplaceengine/internal/validate"
	"github.com/openmarketplaceengine/openmarketplaceengine/srv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/go-redis/redis/v8"
	locationV1beta1 "github.com/openmarketplaceengine/openmarketplaceengine/internal/omeapi/location/v1beta1"
	tollgateCrossingV1beta1 "github.com/openmarketplaceengine/openmarketplaceengine/internal/omeapi/tollgate_crossing/v1beta1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Controller struct {
	locationV1beta1.UnimplementedLocationServiceServer
	tracker *Tracker
}

func GrpcRegister() {
	srv.Grpc.Register(func(srv *grpc.Server) error {
		controller, err := newController(dao.Reds.StoreClient, dao.Reds.PubSubClient)
		if err != nil {
			return err
		}
		locationV1beta1.RegisterLocationServiceServer(srv, controller)
		return nil
	})
}

func newController(storeClient *redis.Client, pubSubClient *redis.Client) (*Controller, error) {
	tracker, err := NewTracker(storeClient, pubSubClient)
	if err != nil {
		return nil, err
	}
	return &Controller{
		tracker: tracker,
	}, nil
}

func (c *Controller) UpdateLocation(ctx context.Context, request *locationV1beta1.UpdateLocationRequest) (*locationV1beta1.UpdateLocationResponse, error) {
	var v validate.Validator
	v.ValidateString("worker_id", request.GetWorkerId(), validate.IsNotNull)
	v.ValidateTimestamp("timestamp", request.GetTimestamp())
	v.ValidateFloat64("lon", request.Lon, validate.IsLongitude)
	v.ValidateFloat64("lat", request.Lat, validate.IsLatitude)

	errorInfo := v.ErrorInfo()
	if errorInfo != nil {
		st, err := status.New(codes.InvalidArgument, "bad request").
			WithDetails(errorInfo)
		if err != nil {
			panic(fmt.Errorf("enrich grpc status with details error: %w", err))
		}
		return nil, st.Err()
	}

	x, err := c.tracker.TrackLocation(ctx, areaKey, request.WorkerId, request.Lon, request.Lat)

	if err != nil {
		st := status.Newf(codes.Internal, "update location or detect tollgate error: %v", err)
		st, err = st.WithDetails(request)
		if err != nil {
			panic(fmt.Errorf("enrich grpc status with details error: %w", err))
		}
		return nil, st.Err()
	}
	return &locationV1beta1.UpdateLocationResponse{
		WorkerId:         request.WorkerId,
		TollgateCrossing: transformCrossing(x),
		Timestamp:        request.Timestamp,
	}, nil
}

func transformCrossing(c *crossing.TollgateCrossing) *tollgateCrossingV1beta1.TollgateCrossing {
	if c == nil {
		return nil
	}
	return &tollgateCrossingV1beta1.TollgateCrossing{
		Id:         c.ID,
		TollgateId: c.TollgateID,
		WorkerId:   c.WorkerID,
		Direction:  string(c.Crossing.Crossing.Direction),
		Alg:        string(c.Crossing.Crossing.Alg),
		Movement: &tollgateCrossingV1beta1.Movement{
			FromLon: c.Crossing.Crossing.Movement.From.Lon,
			FromLat: c.Crossing.Crossing.Movement.From.Lat,
			ToLon:   c.Crossing.Crossing.Movement.To.Lon,
			ToLat:   c.Crossing.Crossing.Movement.To.Lat,
		},
		Created: &timestamppb.Timestamp{
			Seconds: c.Created.Unix(),
			Nanos:   0,
		},
	}
}

func (c *Controller) QueryLocation(ctx context.Context, request *locationV1beta1.QueryLocationRequest) (*locationV1beta1.QueryLocationResponse, error) {
	l := c.tracker.QueryLastLocation(ctx, areaKey, request.WorkerId)
	if l != nil {
		return &locationV1beta1.QueryLocationResponse{
			WorkerId: l.WorkerID,
			Lon:      l.Longitude,
			Lat:      l.Latitude,
			LastSeenTime: &timestamppb.Timestamp{
				Seconds: l.LastSeenTime.Unix(),
				Nanos:   0,
			},
		}, nil
	}
	st := status.Newf(codes.NotFound, "location not found")
	st, err := st.WithDetails(request)
	if err != nil {
		panic(fmt.Errorf("enrich grpc status with details error: %w", err))
	}
	return nil, st.Err()
}
