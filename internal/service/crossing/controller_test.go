package crossing

import (
	"context"
	"log"
	"math/rand"
	"net"
	"testing"
	"time"

	"github.com/openmarketplaceengine/openmarketplaceengine/dom"

	"github.com/google/uuid"
	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
	"github.com/openmarketplaceengine/openmarketplaceengine/internal/omeapi/crossing/v1beta1"
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
	client := v1beta1.NewCrossingServiceClient(conn)

	t.Run("testQuery", func(t *testing.T) {
		testQuery(t, client)
	})
}

func dialer() func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)

	server := grpc.NewServer()
	controller := newController()
	v1beta1.RegisterCrossingServiceServer(server, controller)

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func testQuery(t *testing.T, client v1beta1.CrossingServiceClient) {
	ctx := cfg.Context()
	r := rand.New(rand.NewSource(time.Now().UnixMilli()))

	toll := newRandomTollgate(r, "testCreate")

	err := toll.Insert(ctx)
	require.NoError(t, err)

	workerID1 := uuid.NewString()
	workerID2 := uuid.NewString()
	tollgateID := toll.ID

	for i := 0; i < 5; i++ {
		x1 := newRandomCrossing(r, tollgateID, workerID1)
		err = x1.Insert(ctx)
		require.NoError(t, err)
		x2 := newRandomCrossing(r, tollgateID, workerID2)
		err = x2.Insert(ctx)
		require.NoError(t, err)
	}

	req1 := &v1beta1.ListCrossingsRequest{
		TollgateId: tollgateID,
		WorkerId:   workerID1,
		PageSize:   0,
		PageToken:  "",
	}

	res1, err := client.ListCrossings(ctx, req1)
	require.NoError(t, err)
	require.Len(t, res1.Crossings, 5)
	require.Equal(t, tollgateID, res1.Crossings[0].TollgateId)
	require.Equal(t, workerID1, res1.Crossings[0].WorkerId)
	require.NotEqual(t, res1.Crossings[0].Movement.To.Lon, float64(0))
	require.NotEqual(t, res1.Crossings[0].Movement.To.Lat, float64(0))

	req2 := &v1beta1.ListCrossingsRequest{
		TollgateId: tollgateID,
	}

	res2, err := client.ListCrossings(ctx, req2)
	require.NoError(t, err)
	require.Len(t, res2.Crossings, 10)
}