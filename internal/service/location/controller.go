package location

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/openmarketplaceengine/openmarketplaceengine/dao"
	"github.com/openmarketplaceengine/openmarketplaceengine/internal/service/tollgate/crossing"
	"github.com/openmarketplaceengine/openmarketplaceengine/internal/validate"
	"github.com/openmarketplaceengine/openmarketplaceengine/srv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/go-redis/redis/v8"
	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
	locationV1beta1 "github.com/openmarketplaceengine/openmarketplaceengine/internal/omeapi/location/v1beta1"
	"github.com/openmarketplaceengine/openmarketplaceengine/internal/service/location/storage"
	"github.com/openmarketplaceengine/openmarketplaceengine/internal/service/tollgate"
	"github.com/openmarketplaceengine/openmarketplaceengine/internal/service/tollgate/detector"
	"github.com/openmarketplaceengine/openmarketplaceengine/log"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const areaKey = "global"

type Controller struct {
	locationV1beta1.UnimplementedLocationServiceServer
	store        *storage.Storage
	pubSubClient *redis.Client
	areaKey      string
	detector     *detector.Detector
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
	tollgates, err := tollgate.QueryAll(cfg.Context(), 100)
	if err != nil {
		return nil, err
	}

	d := detector.NewDetector(transformTollgates(tollgates), storeClient)

	return &Controller{
		store:        storage.New(storeClient),
		pubSubClient: pubSubClient,
		areaKey:      areaKey,
		detector:     d,
	}, nil
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

	lastLocation := c.store.LastLocation(ctx, c.areaKey, request.WorkerId)

	err := c.store.Update(ctx, c.areaKey, &storage.Location{
		WorkerID:  request.WorkerId,
		Longitude: request.Lon,
		Latitude:  request.Lat,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "update location error: %s", err)
	}

	c.publishLocation(ctx, request.WorkerId, request.Lon, request.Lat)

	tollgateCrossing, err := c.detectTollgateCrossing(ctx, lastLocation, request)
	if err != nil {
		st := status.Newf(codes.Internal, "detect tollgate crossing error: %v", err)
		st, err = st.WithDetails(request)
		if err != nil {
			panic(fmt.Errorf("enrich grpc status with details error: %w", err))
		}
		return nil, st.Err()
	}
	return &locationV1beta1.UpdateLocationResponse{
		WorkerId:         request.WorkerId,
		TollgateCrossing: tollgateCrossing,
		Timestamp:        request.Timestamp,
	}, nil
}

func (c *Controller) detectTollgateCrossing(ctx context.Context, lastLocation *storage.LastLocation, request *locationV1beta1.UpdateLocationRequest) (*locationV1beta1.TollgateCrossing, error) {
	if lastLocation != nil {
		from := &detector.Location{
			Lon: lastLocation.Longitude,
			Lat: lastLocation.Latitude,
		}
		to := &detector.Location{
			Lon: request.Lon,
			Lat: request.Lat}
		movement := &detector.Movement{
			SubjectID: request.WorkerId,
			From:      from,
			To:        to,
		}

		detected, err := c.detector.DetectCrossing(ctx, movement)
		if err != nil {
			return nil, fmt.Errorf("detect crossing error: %s", err)
		}
		if detected != nil {
			tollgateCrossing := crossing.NewTollgateCrossing(detected.TollgateID, movement.SubjectID, detected)
			err := tollgateCrossing.Insert(ctx)
			if err != nil {
				return nil, fmt.Errorf("crossing insert error: %s", err)
			}
			c.publishTollgateCrossing(ctx, detected)
			return transformCrossing(tollgateCrossing), nil
		}
	}
	return nil, nil
}

func transformCrossing(c *crossing.TollgateCrossing) *locationV1beta1.TollgateCrossing {
	return &locationV1beta1.TollgateCrossing{
		Id:         c.ID,
		TollgateId: c.TollgateID,
		WorkerId:   c.WorkerID,
		Direction:  string(c.Crossing.Crossing.Direction),
		Alg:        string(c.Crossing.Crossing.Alg),
		Movement: &locationV1beta1.Movement{
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
	err = c.pubSubClient.Publish(ctx, channel, payload).Err()

	if err != nil {
		log.Errorf("location publish error: %q", err)
		return
	}
}

func locationChannel(workerID string) string {
	return fmt.Sprintf("channel-location-%s", workerID)
}

func (c *Controller) publishTollgateCrossing(ctx context.Context, crossing *detector.Crossing) {
	channel := crossingChannel(crossing.TollgateID)

	bytes, err := json.Marshal(crossing)
	if err != nil {
		log.Errorf("crossing marshal error: %q", err)
		return
	}
	payload := string(bytes)
	err = c.pubSubClient.Publish(ctx, channel, payload).Err()

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
