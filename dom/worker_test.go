package dom

import (
	"testing"
	"time"

	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
	"github.com/openmarketplaceengine/openmarketplaceengine/dao"
	"github.com/stretchr/testify/require"
)

//-----------------------------------------------------------------------------

func TestWorker_Persist(t *testing.T) {
	WillTest(t, "test", true)
	ctx := cfg.Context()
	for i := 0; i < 100; i++ {
		wrk := genWorker()
		require.NoError(t, wrk.Persist(ctx))
	}
}

//-----------------------------------------------------------------------------

func TestWorkerVehicle_Persist(t *testing.T) {
	WillTest(t, "test", true)
	ctx := cfg.Context()
	for i := 0; i < 100; i++ {
		wv := genWorkerVehicle()
		require.NoError(t, wv.Persist(ctx))
	}
}

//-----------------------------------------------------------------------------

func TestWorker_InsertConstraint(t *testing.T) {
	WillTest(t, "test", true)
	wrk := genWorker()
	require.NoError(t, wrk.Persist(cfg.Context()))
	err := wrk.Persist(cfg.Context())
	require.True(t, dao.ErrUniqueViolation(err))
}

//-----------------------------------------------------------------------------

func TestGetWorker(t *testing.T) {
	WillTest(t, "test", true)
	ctx := cfg.Context()
	for i := 0; i < 100; i++ {
		wput := genWorker()
		require.NoError(t, wput.Persist(ctx))
		wget, err := GetWorker(ctx, wput.ID)
		require.NoError(t, err)
		require.NotNil(t, wget)
		wput.Created.Reset()
		wput.Updated.Reset()
		wget.Created.Reset()
		wget.Updated.Reset()
		require.Equal(t, wput, wget)
		testGetWorkerStatus(t, ctx, wput)
	}
}

//-----------------------------------------------------------------------------

func TestSetWorkerStatus(t *testing.T) {
	WillTest(t, "test", true)
	wrk := genWorker()
	ctx := cfg.Context()
	require.NoError(t, wrk.Persist(ctx))
	testGetWorkerStatus(t, ctx, wrk)
	wrk.Status = WorkerStatus(mockEnum(WorkerDisabled))
	require.NoError(t, SetWorkerStatus(ctx, wrk.ID, wrk.Status))
	testGetWorkerStatus(t, ctx, wrk)
}

//-----------------------------------------------------------------------------

func TestWorker_RowsAffected(t *testing.T) {
	WillTest(t, "test", true)
	max := 100
	ctx := cfg.Context()
	for i := 0; i < max; i++ {
		wrk := genWorker()
		require.NoError(t, wrk.Persist(ctx))
	}
	sql := dao.Update(workerTable).Set("updated", time.Now())
	require.NoError(t, dao.ExecTX(ctx, sql))
	require.Equal(t, max, int(sql.RowsAffected()))
}

//-----------------------------------------------------------------------------
// Helpers
//-----------------------------------------------------------------------------

func testGetWorkerStatus(t *testing.T, ctx Context, wput *Worker) {
	status, er := GetWorkerStatus(ctx, wput.ID)
	require.NoError(t, er)
	require.Equal(t, wput.Status, status)
}

//-----------------------------------------------------------------------------

func genWorker() *Worker {
	stamp := mockStamp()
	return &Worker{
		ID:        mockUUID("drv"),
		Status:    WorkerStatus(mockEnum(WorkerDisabled)),
		Rating:    int32(mockIntn(100)),
		Jobs:      mockIntn(1_000),
		FirstName: mockFirstName(),
		LastName:  mockLastName(),
		Vehicle:   mockUUID("car"),
		Created:   stamp,
		Updated:   stamp,
	}
}

//-----------------------------------------------------------------------------

func genWorkerVehicle() *WorkerVehicle {
	return &WorkerVehicle{
		Worker:  mockUUID("drv"),
		Vehicle: mockUUID("car"),
	}
}
