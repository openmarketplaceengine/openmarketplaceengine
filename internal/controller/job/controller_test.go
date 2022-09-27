package job

import (
	"context"
	"log"
	"net"
	"testing"

	"github.com/driverscooperative/geosrv/cfg"
	"github.com/driverscooperative/geosrv/dao"
	"github.com/driverscooperative/geosrv/dom"
	"github.com/driverscooperative/geosrv/dom/job"
	"github.com/driverscooperative/geosrv/internal/idl/api/job/v1beta1"
	"github.com/driverscooperative/geosrv/pkg/detector"
	svcJob "github.com/driverscooperative/geosrv/svc/job"
	"github.com/driverscooperative/geosrv/svc/location"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

func TestController(t *testing.T) {
	err := cfg.Load()
	require.NoError(t, err)

	dom.WillTest(t, "test", false)
	if !dao.Reds.State.Running() {
		require.NoError(t, dao.Reds.Boot())
	}

	ctx := context.Background()

	storeClient := dao.Reds.StoreClient
	noOp := detector.NewDetectorNoOp()
	storage := location.NewStorage(storeClient)
	tracker := location.NewTracker(storage, noOp)
	service := svcJob.NewService(tracker)

	conn, err := grpc.DialContext(ctx, "", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(dialer(service)))
	defer func(conn *grpc.ClientConn) {
		innerErr := conn.Close()
		if innerErr != nil {
			log.Fatal(innerErr)
		}
	}(conn)

	require.NoError(t, err)
	client := v1beta1.NewJobServiceClient(conn)

	t.Run("testGetAvailableJobs", func(t *testing.T) {
		testGetAvailableJobs(t, client, tracker)
	})
}

func dialer(service *svcJob.Service) func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)

	server := grpc.NewServer()

	v1beta1.RegisterJobServiceServer(server, &controller{jobService: service})

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func testGetAvailableJobs(t *testing.T, client v1beta1.JobServiceClient, tracker *location.Tracker) {
	ctx := cfg.Context()
	err := job.DeleteAll(ctx)
	require.NoError(t, err)
	id := uuid.NewString()
	fromLat := 40.633650
	fromLon := -74.143650

	job1 := &job.Job{
		ID:         uuid.NewString(),
		WorkerID:   "",
		State:      "AVAILABLE",
		PickupLat:  40.636916,
		PickupLon:  -74.195995,
		DropoffLat: 40.634408,
		DropoffLon: -74.198356,
	}
	job2 := &job.Job{
		ID:         uuid.NewString(),
		WorkerID:   "",
		State:      "AVAILABLE",
		PickupLat:  40.634408,
		PickupLon:  -74.198356,
		DropoffLat: 40.636916,
		DropoffLon: -74.195995,
	}

	for _, j := range []*job.Job{job1, job2} {
		_, _, innerErr := j.Upsert(ctx)
		require.NoError(t, innerErr)
	}
	areaKey := "test-tracker"
	_, err = tracker.TrackLocation(ctx, areaKey, id, fromLon, fromLat)
	require.NoError(t, err)

	req1 := &v1beta1.GetAvailableJobsRequest{
		AreaKey:      areaKey,
		WorkerId:     id,
		RadiusMeters: 10000,
		Limit:        24,
	}

	res1, err := client.GetAvailableJobs(ctx, req1)
	require.NoError(t, err)
	require.Len(t, res1.Jobs, 2)
}
