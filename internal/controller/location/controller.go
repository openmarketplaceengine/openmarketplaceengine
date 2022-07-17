package location

import (
	"context"
	"fmt"
	"time"

	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
	"github.com/openmarketplaceengine/openmarketplaceengine/dao"
	"github.com/openmarketplaceengine/openmarketplaceengine/dom/tollgate"
	"github.com/openmarketplaceengine/openmarketplaceengine/pkg/detector"
	"github.com/openmarketplaceengine/openmarketplaceengine/srv"
	"github.com/openmarketplaceengine/openmarketplaceengine/svc/location"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	locationV1beta1 "github.com/openmarketplaceengine/openmarketplaceengine/internal/idl/api/location/v1beta1"
	"github.com/openmarketplaceengine/openmarketplaceengine/internal/idl/api/type/v1beta1"
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

func (c *controller) UpdateLocation(ctx context.Context, req *locationV1beta1.UpdateLocationRequest) (*locationV1beta1.UpdateLocationResponse, error) {
	err := req.ValidateAll()

	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	areaKey := req.GetAreaKey()
	workerID := req.GetWorkerId()
	loc := req.GetLocation()

	x, err := c.tracker.TrackLocation(ctx, areaKey, workerID, loc.GetLongitude(), loc.GetLatitude())

	if err != nil {
		st := status.Newf(codes.Internal, "update location or detect tollgate error: %v", err)
		st, err = st.WithDetails(req)
		if err != nil {
			panic(fmt.Errorf("enrich grpc status with details error: %w", err))
		}
		return nil, st.Err()
	}
	return &locationV1beta1.UpdateLocationResponse{
		AreaKey:  areaKey,
		WorkerId: workerID,
		Crossing: transform(x),
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

func (c *controller) GetLocation(ctx context.Context, req *locationV1beta1.GetLocationRequest) (*locationV1beta1.GetLocationResponse, error) {
	err := req.ValidateAll()

	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	workerID := req.GetWorkerId()
	areaKey := req.GetAreaKey()

	l := c.tracker.GetLocation(ctx, areaKey, workerID)
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
	return nil, status.Errorf(codes.NotFound, "location not found")
}
