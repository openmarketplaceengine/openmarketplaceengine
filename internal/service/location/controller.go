package location

import (
	"context"
	"fmt"

	"github.com/openmarketplaceengine/openmarketplaceengine/dao"
	"github.com/openmarketplaceengine/openmarketplaceengine/internal/service/crossing"
	"github.com/openmarketplaceengine/openmarketplaceengine/internal/validate"
	"github.com/openmarketplaceengine/openmarketplaceengine/srv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/go-redis/redis/v8"
	locationV1beta1 "github.com/openmarketplaceengine/openmarketplaceengine/internal/omeapi/location/v1beta1"
	"github.com/openmarketplaceengine/openmarketplaceengine/internal/omeapi/type/v1beta1"
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
	areaKey := request.GetAreaKey()
	value := request.GetValue()
	workerID := value.GetWorkerId()
	location := value.GetLocation()
	updateTime := value.GetUpdateTime()
	var v validate.Validator
	v.ValidateString("area_key", areaKey, validate.IsNotNull)
	v.ValidateString("value_worker_id", workerID, validate.IsNotNull)
	v.ValidateFloat64("value_location_lon", location.GetLon(), validate.IsLongitude)
	v.ValidateFloat64("value_location_lat", location.GetLat(), validate.IsLatitude)
	v.ValidateTimestamp("value_update_time", updateTime)

	errorInfo := v.ErrorInfo()
	if errorInfo != nil {
		st, err := status.New(codes.InvalidArgument, "bad request").
			WithDetails(errorInfo)
		if err != nil {
			panic(fmt.Errorf("enrich grpc status with details error: %w", err))
		}
		return nil, st.Err()
	}

	x, err := c.tracker.TrackLocation(ctx, areaKey, workerID, location.GetLon(), location.GetLat())

	if err != nil {
		st := status.Newf(codes.Internal, "update location or detect tollgate error: %v", err)
		st, err = st.WithDetails(request)
		if err != nil {
			panic(fmt.Errorf("enrich grpc status with details error: %w", err))
		}
		return nil, st.Err()
	}
	return &locationV1beta1.UpdateLocationResponse{
		AreaKey:    areaKey,
		WorkerId:   workerID,
		Crossing:   transform(x),
		UpdateTime: updateTime,
	}, nil
}

func transform(c *crossing.TollgateCrossing) *v1beta1.Crossing {
	if c == nil {
		return nil
	}
	return &v1beta1.Crossing{
		Id:         c.ID,
		TollgateId: c.TollgateID,
		WorkerId:   c.WorkerID,
		Direction:  string(c.Crossing.Crossing.Direction),
		Alg:        string(c.Crossing.Crossing.Alg),
		Movement: &v1beta1.Movement{
			From: &v1beta1.Location{
				Lat: c.Crossing.Crossing.Movement.From.Lat,
				Lon: c.Crossing.Crossing.Movement.From.Lon,
			},
			To: &v1beta1.Location{
				Lat: c.Crossing.Crossing.Movement.To.Lat,
				Lon: c.Crossing.Crossing.Movement.To.Lon,
			},
		},
		CreateTime: timestamppb.New(c.Created.Time),
	}
}

func (c *Controller) GetLocation(ctx context.Context, request *locationV1beta1.GetLocationRequest) (*locationV1beta1.GetLocationResponse, error) {
	workerID := request.GetWorkerId()
	areaKey := request.GetAreaKey()
	var v validate.Validator
	v.ValidateString("worker_id", workerID, validate.IsNotNull)
	v.ValidateString("area_key", areaKey, validate.IsNotNull)

	errorInfo := v.ErrorInfo()
	if errorInfo != nil {
		st, err := status.New(codes.InvalidArgument, "bad request").
			WithDetails(errorInfo)
		if err != nil {
			panic(fmt.Errorf("enrich grpc status with details error: %w", err))
		}
		return nil, st.Err()
	}

	l := c.tracker.QueryLastLocation(ctx, areaKey, workerID)
	if l != nil {
		return &locationV1beta1.GetLocationResponse{
			AreaKey:  areaKey,
			WorkerId: l.WorkerID,
			Location: &v1beta1.Location{
				Lon: l.Longitude,
				Lat: l.Latitude,
			},
			LastSeenTime: timestamppb.New(l.LastSeenTime),
		}, nil
	}
	st := status.Newf(codes.NotFound, "location not found")
	st, err := st.WithDetails(request)
	if err != nil {
		panic(fmt.Errorf("enrich grpc status with details error: %w", err))
	}
	return nil, st.Err()
}
