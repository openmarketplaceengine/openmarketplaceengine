package tollgate

import (
	"context"
	"log"
	"net"
	"testing"

	"github.com/driverscooperative/geosrv/dom"
	"github.com/driverscooperative/geosrv/dom/tollgate"
	"github.com/driverscooperative/geosrv/pkg/detector"
	"github.com/google/uuid"

	"github.com/driverscooperative/geosrv/cfg"
	"github.com/driverscooperative/geosrv/internal/idl/api/tollgate/v1beta1"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

func TestController(t *testing.T) {
	err := cfg.Load()
	require.NoError(t, err)

	dom.WillTest(t, "test", false)
	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, "", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(dialer()))
	defer func(conn *grpc.ClientConn) {
		innerErr := conn.Close()
		if innerErr != nil {
			log.Fatal(innerErr)
		}
	}(conn)

	require.NoError(t, err)
	client := v1beta1.NewTollgateServiceClient(conn)

	t.Run("testQueryAll", func(t *testing.T) {
		testQueryOne(t, client)
	})
}

func dialer() func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)

	server := grpc.NewServer()
	v1beta1.RegisterTollgateServiceServer(server, &controller{})

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func testQueryOne(t *testing.T, client v1beta1.TollgateServiceClient) {
	ctx := cfg.Context()
	toll := newTollgate(uuid.NewString(), "testCreate")

	err := toll.Insert(ctx)
	require.NoError(t, err)

	tollgateID := toll.ID

	res, err := client.GetTollgate(ctx, &v1beta1.GetTollgateRequest{TollgateId: tollgateID})
	require.NoError(t, err)
	require.Equal(t, tollgateID, res.Tollgate.Id)
}

func newTollgate(id string, name string) *tollgate.Tollgate {
	return &tollgate.Tollgate{
		ID:   id,
		Name: name,
		BBoxes: &tollgate.BBoxes{
			BBoxes: []*detector.BBox{{
				LonMin: 0,
				LatMin: 0,
				LonMax: 0,
				LatMax: 0,
			}},
			Required: 2,
		},
		GateLine: &tollgate.GateLine{
			Line: &detector.Line{
				Lon1: 0,
				Lat1: 0,
				Lon2: 0,
				Lat2: 0,
			},
		},
	}
}
