package tollgate

import (
	"context"
	"log"
	"math/rand"
	"net"
	"testing"
	"time"

	"github.com/openmarketplaceengine/openmarketplaceengine/dom"

	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
	"github.com/openmarketplaceengine/openmarketplaceengine/internal/omeapi/tollgate/v1beta1"
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
	controller := newController()
	v1beta1.RegisterTollgateServiceServer(server, controller)

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
	r := rand.New(rand.NewSource(time.Now().UnixMilli()))

	toll := newRandomTollgate(r, "testCreate")

	err := toll.Insert(ctx)
	require.NoError(t, err)

	tollgateID := toll.ID

	res, err := client.GetTollgate(ctx, &v1beta1.GetTollgateRequest{TollgateId: tollgateID})
	require.NoError(t, err)
	require.Equal(t, tollgateID, res.Tollgate.Id)
}
