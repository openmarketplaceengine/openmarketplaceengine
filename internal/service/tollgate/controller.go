package tollgate

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
	"github.com/openmarketplaceengine/openmarketplaceengine/dom/tollgate"
	"github.com/openmarketplaceengine/openmarketplaceengine/internal/api/tollgate/v1beta1"
	"github.com/openmarketplaceengine/openmarketplaceengine/srv"
	svcTollgate "github.com/openmarketplaceengine/openmarketplaceengine/svc/tollgate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

		err := svcTollgate.Load(cfg.Context())
		if err != nil {
			return fmt.Errorf("load tollgates error: %w", err)
		}

		return nil
	})
}

func (c *Controller) GetTollgate(ctx context.Context, request *v1beta1.GetTollgateRequest) (*v1beta1.GetTollgateResponse, error) {
	toll, err := tollgate.QueryOne(ctx, request.TollgateId)
	if err != nil {
		if err == sql.ErrNoRows {
			st := status.New(codes.NotFound, "Tollgate not found")
			st, innerErr := st.WithDetails(request)
			if innerErr != nil {
				panic(fmt.Errorf("enrich grpc status with details error: %w", innerErr))
			}
			return nil, st.Err()
		}
		return nil, status.Errorf(codes.Internal, "query tollgate error: %s", err)
	}
	return &v1beta1.GetTollgateResponse{
		Tollgate: transform(toll),
	}, nil
}

func (c *Controller) ListTollgates(ctx context.Context, request *v1beta1.ListTollgatesRequest) (*v1beta1.ListTollgatesResponse, error) {
	all, err := tollgate.QueryAll(ctx, 100)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "QueryAll error: %s", err)
	}
	return &v1beta1.ListTollgatesResponse{
		Tollgates:     transformAll(all),
		NextPageToken: "",
	}, nil
}

func transform(t *tollgate.Tollgate) *v1beta1.Tollgate {
	var bBoxes *v1beta1.BBoxes
	var gLine *v1beta1.GateLine

	if t.GateLine != nil {
		gLine = &v1beta1.GateLine{
			Lon1: t.GateLine.Line.Lon1,
			Lat1: t.GateLine.Line.Lat1,
			Lon2: t.GateLine.Line.Lon2,
			Lat2: t.GateLine.Line.Lat2,
		}
	}

	if t.BBoxes != nil {
		bb := make([]*v1beta1.BBox, 0)
		for _, box := range t.BBoxes.BBoxes {
			bb = append(bb, &v1beta1.BBox{
				LonMin: box.LonMin,
				LatMin: box.LatMin,
				LonMax: box.LonMax,
				LatMax: box.LatMax,
			})
		}
		bBoxes = &v1beta1.BBoxes{
			BBoxes:   bb,
			Required: t.BBoxes.Required,
		}
	}

	return &v1beta1.Tollgate{
		Id:       t.ID,
		Name:     t.Name,
		BBoxes:   bBoxes,
		GateLine: gLine,
		Created:  timestamppb.New(t.Created.Time),
		Updated:  timestamppb.New(t.Updated.Time),
	}
}

func transformAll(tollgates []*tollgate.Tollgate) []*v1beta1.Tollgate {
	result := make([]*v1beta1.Tollgate, 0)
	for _, t := range tollgates {
		result = append(result, transform(t))
	}
	return result
}
