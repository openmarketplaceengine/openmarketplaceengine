package location

import (
	"context"
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	"github.com/openmarketplaceengine/openmarketplaceengine/internal/service/tollgate/detector"

	lineTollgate "github.com/openmarketplaceengine/openmarketplaceengine/internal/service/tollgate/line"

	"github.com/google/uuid"
	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
	locationV1beta1 "github.com/openmarketplaceengine/openmarketplaceengine/internal/omeapi/location/v1beta1"
	"github.com/openmarketplaceengine/openmarketplaceengine/internal/service/tollgate"
	redisClient "github.com/openmarketplaceengine/openmarketplaceengine/redis/client"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const areaKey = "san_fran"
const tollgateID = "tollgate-123"

func TestController(t *testing.T) {
	err := cfg.Load()
	require.NoError(t, err)

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(dialer()))
	defer func(conn *grpc.ClientConn) {
		innerErr := conn.Close()
		if innerErr != nil {
			log.Fatal(innerErr)
		}
	}(conn)

	require.NoError(t, err)
	client := locationV1beta1.NewLocationServiceClient(conn)

	t.Run("testUpdateLocation", func(t *testing.T) {
		testUpdateLocation(t, client)
	})
	t.Run("testQueryLocation", func(t *testing.T) {
		testQueryLocation(t, client)
	})
	t.Run("testTollgateCrossing", func(t *testing.T) {
		testTollgateCrossing(t, client)
	})
}

func dialer() func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)

	server := grpc.NewServer()
	tg := lineTollgate.NewTollgate(
		tollgateID,
		&tollgate.LocationXY{
			LongitudeX: -79.870262,
			LatitudeY:  41.198497,
		},
		&tollgate.LocationXY{
			LongitudeX: -79.870218,
			LatitudeY:  41.200268,
		},
	)
	tollgateDetector := detector.NewDetector()
	tollgateDetector.AddTollgate(tg)
	controller := New(redisClient.NewStoreClient(), redisClient.NewPubSubClient(), areaKey, tollgateDetector)
	locationV1beta1.RegisterLocationServiceServer(server, controller)

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func testUpdateLocation(t *testing.T, client locationV1beta1.LocationServiceClient) {
	id := uuid.NewString()
	request := &locationV1beta1.UpdateLocationRequest{
		WorkerId:  id,
		Longitude: 12.000001966953278,
		Latitude:  13.000001966953278,
	}
	response, err := client.UpdateLocation(context.Background(), request)
	require.NoError(t, err)
	require.Equal(t, request.WorkerId, response.WorkerId)

	location, err := client.QueryLocation(context.Background(), &locationV1beta1.QueryLocationRequest{
		WorkerId: id,
	})
	require.NoError(t, err)
	require.Equal(t, request.WorkerId, location.WorkerId)
	require.InDelta(t, request.Longitude, location.Longitude, 0.001)
	require.InDelta(t, request.Latitude, location.Latitude, 0.001)
}

func testQueryLocation(t *testing.T, client locationV1beta1.LocationServiceClient) {
	id := uuid.NewString()
	request := &locationV1beta1.QueryLocationRequest{
		WorkerId: id,
	}

	ctx := context.Background()

	_, err := client.QueryLocation(ctx, request)
	require.Error(t, err)
	require.Contains(t, err.Error(), fmt.Sprintf("location not found for WorkerId=%s", request.WorkerId))

	response, err := client.UpdateLocation(ctx, &locationV1beta1.UpdateLocationRequest{
		WorkerId:  id,
		Longitude: 12,
		Latitude:  13,
	},
	)
	require.NoError(t, err)
	require.Equal(t, request.WorkerId, response.WorkerId)

	location, err := client.QueryLocation(ctx, request)
	require.NoError(t, err)
	require.Less(t, location.LastSeenTime.AsTime().UnixMilli(), time.Now().UnixMilli())
}

func testTollgateCrossing(t *testing.T, client locationV1beta1.LocationServiceClient) {
	ctx := context.Background()
	id := uuid.NewString()
	from := &locationV1beta1.UpdateLocationRequest{
		WorkerId:  id,
		Longitude: -79.871248,
		Latitude:  41.199493,
	}
	to := &locationV1beta1.UpdateLocationRequest{
		WorkerId:  id,
		Longitude: -79.867927,
		Latitude:  41.199329,
	}

	sync := make(chan string)
	var crossings <-chan tollgate.Crossing
	go func() {
		crossings = subscribe(crossingChannel(tollgateID))
		sync <- "done"
	}()

	select {
	case <-sync:
		break
	case <-time.After(5 * time.Second):
		require.Fail(t, "sync timeout")
	}

	response1, err := client.UpdateLocation(ctx, from)
	require.NoError(t, err)
	require.Equal(t, from.WorkerId, response1.WorkerId)
	require.Nil(t, response1.TollgateCrossing)

	response2, err := client.UpdateLocation(ctx, to)
	require.NoError(t, err)
	require.Equal(t, from.WorkerId, response2.WorkerId)
	require.NotNil(t, response2.TollgateCrossing)

	c := <-crossings
	require.Equal(t, tollgateID, c.TollgateID)
	require.Equal(t, tollgate.Direction("SE"), c.Direction)
	require.Equal(t, id, c.SubjectID)
	require.InDelta(t, to.Latitude, c.Location.LatitudeY, 0.003)
	require.InDelta(t, to.Longitude, c.Location.LongitudeX, 0.003)
}
