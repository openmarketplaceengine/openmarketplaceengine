package tollgate

import (
	"context"

	"github.com/openmarketplaceengine/openmarketplaceengine/srv"
	"google.golang.org/grpc"

	"github.com/openmarketplaceengine/openmarketplaceengine/internal/omeapi/tollgate/v1beta1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Controller struct {
	v1beta1.UnimplementedTollgateServiceServer
}

func newController() *Controller {
	return &Controller{}
}

func GrpcRegister() {
	srv.Grpc.Register(func(srv *grpc.Server) error {
		controller := newController()
		v1beta1.RegisterTollgateServiceServer(srv, controller)
		return nil
	})
}

func (c *Controller) QueryOne(ctx context.Context, request *v1beta1.QueryOneRequest) (*v1beta1.QueryOneResponse, error) {
	toll, err := QueryOne(ctx, request.TollgateId)
	if err != nil {
		return nil, err
	}
	return &v1beta1.QueryOneResponse{
		Tollgate: transform(toll),
	}, nil
}

func (c *Controller) QueryAll(ctx context.Context, request *v1beta1.QueryAllRequest) (*v1beta1.QueryAllResponse, error) {
	all, err := QueryAll(ctx, 100)
	if err != nil {
		return nil, err
	}
	return &v1beta1.QueryAllResponse{
		Tollgates: transformAll(all),
	}, nil
}

func transform(toll *Tollgate) *v1beta1.Tollgate {
	var bBoxes *v1beta1.BBoxes
	var gLine *v1beta1.GateLine

	if toll.GateLine != nil {
		gLine = &v1beta1.GateLine{
			Lon1: toll.GateLine.Line.Lon1,
			Lat1: toll.GateLine.Line.Lat1,
			Lon2: toll.GateLine.Line.Lon2,
			Lat2: toll.GateLine.Line.Lat2,
		}
	}

	if toll.BBoxes != nil {
		bb := make([]*v1beta1.BBox, 0)
		for _, box := range toll.BBoxes.BBoxes {
			bb = append(bb, &v1beta1.BBox{
				LonMin: box.LonMin,
				LatMin: box.LatMin,
				LonMax: box.LonMax,
				LatMax: box.LatMax,
			})
		}
		bBoxes = &v1beta1.BBoxes{
			BBoxes:   bb,
			Required: toll.BBoxes.Required,
		}
	}

	return &v1beta1.Tollgate{
		Id:       toll.ID,
		Name:     toll.Name,
		BBoxes:   bBoxes,
		GateLine: gLine,
		Created: &timestamppb.Timestamp{
			Seconds: toll.Created.Unix(),
			Nanos:   0,
		},
		Updated: &timestamppb.Timestamp{
			Seconds: toll.Updated.Unix(),
			Nanos:   0,
		},
	}
}

func transformAll(tollgates []*Tollgate) []*v1beta1.Tollgate {
	result := make([]*v1beta1.Tollgate, 0)
	for _, tollgate := range tollgates {
		result = append(result, transform(tollgate))
	}
	return result
}
