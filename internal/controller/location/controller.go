package location

import (
	"context"
	"fmt"
	"time"

	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
	"github.com/openmarketplaceengine/openmarketplaceengine/dao"
	"github.com/openmarketplaceengine/openmarketplaceengine/dom/tollgate"
	"github.com/openmarketplaceengine/openmarketplaceengine/internal/validate"
	"github.com/openmarketplaceengine/openmarketplaceengine/pkg/detector"
	"github.com/openmarketplaceengine/openmarketplaceengine/srv"
	"github.com/openmarketplaceengine/openmarketplaceengine/svc/location"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	locationV1beta1 "github.com/openmarketplaceengine/openmarketplaceengine/internal/api/location/v1beta1"
	"github.com/openmarketplaceengine/openmarketplaceengine/internal/api/type/v1beta1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type controller struct {
	locationV1beta1.UnimplementedLocationServiceServer
	tracker *location.Tracker
}

func newController() (*controller, error) {
	tollgates, err := tollgate.QueryAll(cfg.Context(), 100)
	if err != nil {
		return nil, err
	}
	storeClient := dao.Reds.StoreClient
	s := location.NewStorage(storeClient)
	d := detector.NewDetector(transformTollgates(tollgates), location.NewBBoxStorage(storeClient))
	tracker := location.NewTracker(s, d)
	if err != nil {
		return nil, err
	}
	return &controller{
		tracker: tracker,
	}, nil
}

func init() {
	srv.Grpc.Register(func(s *grpc.Server) error {
		srv.Grpc.Infof("registering: %s", locationV1beta1.LocationService_ServiceDesc.ServiceName)
		controller, err := newController()
		if err != nil {
			return err
		}
		locationV1beta1.RegisterLocationServiceServer(s, controller)
		return nil
	})
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

func (c *controller) UpdateLocation(ctx context.Context, request *locationV1beta1.UpdateLocationRequest) (*locationV1beta1.UpdateLocationResponse, error) {
	areaKey := request.GetAreaKey()
	value := request.GetValue()
	workerID := value.GetWorkerId()
	loc := value.GetLocation()
	updateTime := value.GetUpdateTime()
	var v validate.Validator
	v.ValidateString("area_key", areaKey).NotEmpty()
	v.ValidateString("value_worker_id", workerID).NotEmpty()
	v.ValidateFloat64("value_location_longitude", loc.GetLongitude(), validate.IsLongitude)
	v.ValidateFloat64("value_location_latitude", loc.GetLatitude(), validate.IsLatitude)
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

	x, err := c.tracker.TrackLocation(ctx, areaKey, workerID, loc.GetLongitude(), loc.GetLatitude())

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

func transform(c *detector.Crossing) *v1beta1.Crossing {
	if c == nil {
		return nil
	}
	return &v1beta1.Crossing{
		Id:         c.ID,
		TollgateId: c.TollgateID,
		WorkerId:   c.WorkerID,
		Direction:  string(c.Direction),
		Alg:        string(c.Alg),
		Movement: &v1beta1.Movement{
			From: &v1beta1.Location{
				Latitude:  c.Movement.From.Latitude,
				Longitude: c.Movement.From.Longitude,
			},
			To: &v1beta1.Location{
				Latitude:  c.Movement.To.Latitude,
				Longitude: c.Movement.To.Longitude,
			},
		},
		CreateTime: timestamppb.New(time.Now()),
	}
}

func (c *controller) GetLocation(ctx context.Context, request *locationV1beta1.GetLocationRequest) (*locationV1beta1.GetLocationResponse, error) {
	workerID := request.GetWorkerId()
	areaKey := request.GetAreaKey()
	var v validate.Validator
	v.ValidateString("worker_id", workerID).NotEmpty()
	v.ValidateString("area_key", areaKey).NotEmpty()

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
				Longitude: l.Longitude,
				Latitude:  l.Latitude,
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
