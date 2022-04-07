package crossing

import (
	"context"
	"fmt"

	"github.com/openmarketplaceengine/openmarketplaceengine/srv"
	"google.golang.org/grpc"

	"github.com/openmarketplaceengine/openmarketplaceengine/internal/omeapi/tollgate_crossing/v1beta1"
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

func (c *Controller) QueryTollgateCrossings(ctx context.Context, request *v1beta1.QueryTollgateCrossingsRequest) (*v1beta1.QueryTollgateCrossingsResponse, error) {
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
		return nil, err
	}

	return &v1beta1.QueryTollgateCrossingsResponse{
		Tollgate: transform(crossings),
	}, nil
}

func transform(crossings []*TollgateCrossing) []*v1beta1.TollgateCrossing {
	result := make([]*v1beta1.TollgateCrossing, 0)
	for _, crossing := range crossings {
		result = append(result, &v1beta1.TollgateCrossing{
			Id:         crossing.ID,
			TollgateId: crossing.TollgateID,
			WorkerId:   crossing.WorkerID,
			Direction:  string(crossing.Crossing.Crossing.Direction),
			Alg:        fmt.Sprintf("%v", crossing.Crossing.Crossing.Alg),
			Movement: &v1beta1.Movement{
				FromLon: crossing.Crossing.Crossing.Movement.From.Lon,
				FromLat: crossing.Crossing.Crossing.Movement.From.Lat,
				ToLon:   crossing.Crossing.Crossing.Movement.To.Lon,
				ToLat:   crossing.Crossing.Crossing.Movement.To.Lat,
			},
			CreatedTime: &timestamppb.Timestamp{
				Seconds: crossing.Created.Unix(),
				Nanos:   0,
			},
		})
	}
	return result
}
