package crossing

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/openmarketplaceengine/openmarketplaceengine/srv"
	"google.golang.org/grpc"

	"github.com/openmarketplaceengine/openmarketplaceengine/internal/omeapi/tollgate_crossing/v1beta1"
	typeV1beta1 "github.com/openmarketplaceengine/openmarketplaceengine/internal/omeapi/type/v1beta1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Controller struct {
	v1beta1.UnimplementedTollgateCrossingServiceServer
}

func newController() *Controller {
	return &Controller{}
}

func GrpcRegister() {
	srv.Grpc.Register(func(srv *grpc.Server) error {
		controller := newController()
		v1beta1.RegisterTollgateCrossingServiceServer(srv, controller)
		return nil
	})
}

func (c *Controller) ListTollgateCrossings(ctx context.Context, request *v1beta1.ListTollgateCrossingsRequest) (*v1beta1.ListTollgateCrossingsResponse, error) {
	wheres := make([]Where, 0)
	if request.TollgateId != "" {
		wheres = append(wheres, Where{
			Expr: "tollgate_id = ?",
			Args: []interface{}{request.TollgateId},
		})
	}

	if request.WorkerId != "" {
		wheres = append(wheres, Where{
			Expr: "worker_id = ?",
			Args: []interface{}{request.WorkerId},
		})
	}

	orderBy := []string{"created desc"}
	crossings, err := QueryBy(ctx, wheres, orderBy, 100)
	if err != nil {
		st := status.Newf(codes.Internal, "query tollgate crossing error: %s", err)
		st, err = st.WithDetails(request)
		if err != nil {
			panic(fmt.Errorf("enrich grpc status with details error: %w", err))
		}
		return nil, st.Err()
	}

	return &v1beta1.ListTollgateCrossingsResponse{
		Crossings:     transform(crossings),
		NextPageToken: "",
	}, nil
}

func transform(crossings []*TollgateCrossing) []*typeV1beta1.TollgateCrossing {
	result := make([]*typeV1beta1.TollgateCrossing, 0)
	for _, crossing := range crossings {
		result = append(result, &typeV1beta1.TollgateCrossing{
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
			CreateTime: &timestamppb.Timestamp{
				Seconds: crossing.Created.Unix(),
				Nanos:   0,
			},
		})
	}
	return result
}
