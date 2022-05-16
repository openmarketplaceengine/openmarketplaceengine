package crossing

import (
	"context"
	"fmt"

	"github.com/openmarketplaceengine/openmarketplaceengine/dom/crossing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/openmarketplaceengine/openmarketplaceengine/srv"
	"google.golang.org/grpc"

	"github.com/openmarketplaceengine/openmarketplaceengine/internal/omeapi/crossing/v1beta1"
	typeV1beta1 "github.com/openmarketplaceengine/openmarketplaceengine/internal/omeapi/type/v1beta1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Controller struct {
	v1beta1.UnimplementedCrossingServiceServer
}

func newController() *Controller {
	return &Controller{}
}

func GrpcRegister() {
	srv.Grpc.Register(func(srv *grpc.Server) error {
		controller := newController()
		v1beta1.RegisterCrossingServiceServer(srv, controller)
		return nil
	})
}

func (c *Controller) ListCrossings(ctx context.Context, request *v1beta1.ListCrossingsRequest) (*v1beta1.ListCrossingsResponse, error) {
	wheres := make([]crossing.Where, 0)
	if request.TollgateId != "" {
		wheres = append(wheres, crossing.Where{
			Expr: "tollgate_id = ?",
			Args: []interface{}{request.TollgateId},
		})
	}

	if request.WorkerId != "" {
		wheres = append(wheres, crossing.Where{
			Expr: "worker_id = ?",
			Args: []interface{}{request.WorkerId},
		})
	}

	orderBy := []string{"created desc"}
	crossings, err := crossing.QueryBy(ctx, wheres, orderBy, 100)
	if err != nil {
		st := status.Newf(codes.Internal, "query tollgate crossing error: %s", err)
		st, err = st.WithDetails(request)
		if err != nil {
			panic(fmt.Errorf("enrich grpc status with details error: %w", err))
		}
		return nil, st.Err()
	}

	return &v1beta1.ListCrossingsResponse{
		Crossings:     transform(crossings),
		NextPageToken: "",
	}, nil
}

func transform(crossings []*crossing.TollgateCrossing) []*typeV1beta1.Crossing {
	result := make([]*typeV1beta1.Crossing, 0)
	for _, crossing := range crossings {
		result = append(result, &typeV1beta1.Crossing{
			Id:         crossing.ID,
			TollgateId: crossing.TollgateID,
			WorkerId:   crossing.WorkerID,
			Direction:  string(crossing.Crossing.Crossing.Direction),
			Alg:        fmt.Sprintf("%v", crossing.Crossing.Crossing.Alg),
			Movement: &typeV1beta1.Movement{
				From: &typeV1beta1.Location{
					Lat: crossing.Crossing.Crossing.Movement.From.Lat,
					Lon: crossing.Crossing.Crossing.Movement.From.Lon,
				},
				To: &typeV1beta1.Location{
					Lat: crossing.Crossing.Crossing.Movement.To.Lat,
					Lon: crossing.Crossing.Crossing.Movement.To.Lon,
				},
			},
			CreateTime: timestamppb.New(crossing.Created.Time),
		})
	}
	return result
}
