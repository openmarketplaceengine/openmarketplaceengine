package crossing

import (
	"context"
	"fmt"

	"github.com/openmarketplaceengine/openmarketplaceengine/internal/omeapi/tollgate_crossing/v1beta1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Controller struct {
	v1beta1.UnimplementedTollgateCrossingServiceServer
}

func New() *Controller {
	return &Controller{}
}

func (c *Controller) QueryTollgateCrossings(ctx context.Context, request *v1beta1.QueryTollgateCrossingsRequest) (*v1beta1.QueryTollgateCrossingsResponse, error) {
	wheres := make([]Where, 0)
	if request.TollgateId != "" {
		wheres = append(wheres, Where{
			Expr: "tollgate_id = ?",
			Args: []interface{}{request.TollgateId},
		})
	}

	if request.DriverId != "" {
		wheres = append(wheres, Where{
			Expr: "driver_id = ?",
			Args: []interface{}{request.DriverId},
		})
	}

	orderBy := []string{"created_at desc"}
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
			DriverId:   crossing.DriverID,
			Direction:  string(crossing.Crossing.Crossing.Direction),
			Alg:        fmt.Sprintf("%v", crossing.Crossing.Crossing.Alg),
			Lon:        crossing.Crossing.Crossing.Location.Lon,
			Lat:        crossing.Crossing.Crossing.Location.Lat,
			CreatedTime: &timestamppb.Timestamp{
				Seconds: crossing.CreatedAt.Unix(),
				Nanos:   0,
			},
		})
	}
	return result
}
