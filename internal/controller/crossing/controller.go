package crossing

import (
	"context"
	"fmt"

	"github.com/openmarketplaceengine/openmarketplaceengine/dom/crossing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/openmarketplaceengine/openmarketplaceengine/srv"
	"google.golang.org/grpc"

	"github.com/openmarketplaceengine/openmarketplaceengine/internal/api/crossing/v1beta1"
	typeV1beta1 "github.com/openmarketplaceengine/openmarketplaceengine/internal/api/type/v1beta1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type controller struct {
	v1beta1.UnimplementedCrossingServiceServer
}

func init() {
	srv.Grpc.Register(func(s *grpc.Server) error {
		srv.Grpc.Infof("registering: %s", v1beta1.CrossingService_ServiceDesc.ServiceName)
		v1beta1.RegisterCrossingServiceServer(s, &controller{})
		return nil
	})
}

func (c *controller) ListCrossings(ctx context.Context, req *v1beta1.ListCrossingsRequest) (*v1beta1.ListCrossingsResponse, error) {
	err := req.ValidateAll()

	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	wheres := make([]crossing.Where, 0)
	if req.TollgateId != "" {
		wheres = append(wheres, crossing.Where{
			Expr: "tollgate_id = ?",
			Args: []interface{}{req.TollgateId},
		})
	}

	if req.WorkerId != "" {
		wheres = append(wheres, crossing.Where{
			Expr: "worker_id = ?",
			Args: []interface{}{req.WorkerId},
		})
	}

	orderBy := []string{"created desc"}
	crossings, err := crossing.QueryBy(ctx, wheres, orderBy, 100)
	if err != nil {
		st := status.Newf(codes.Internal, "query tollgate crossing error: %s", err)
		st, err = st.WithDetails(req)
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
	for _, cross := range crossings {
		result = append(result, &typeV1beta1.Crossing{
			Id:         cross.ID,
			TollgateId: cross.TollgateID,
			WorkerId:   cross.WorkerID,
			Direction:  string(cross.Crossing.Crossing.Direction),
			Alg:        fmt.Sprintf("%v", cross.Crossing.Crossing.Alg),
			Movement: &typeV1beta1.Movement{
				From: &typeV1beta1.Location{
					Latitude:  cross.Crossing.Crossing.Movement.From.Latitude,
					Longitude: cross.Crossing.Crossing.Movement.From.Longitude,
				},
				To: &typeV1beta1.Location{
					Latitude:  cross.Crossing.Crossing.Movement.To.Latitude,
					Longitude: cross.Crossing.Crossing.Movement.To.Longitude,
				},
			},
			CreateTime: timestamppb.New(cross.Created.Time),
		})
	}
	return result
}
