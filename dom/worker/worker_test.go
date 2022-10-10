package worker

import (
	"testing"
	"time"

	"github.com/openmarketplaceengine/openmarketplaceengine/dom"

	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
	"github.com/openmarketplaceengine/openmarketplaceengine/dao"
	"github.com/stretchr/testify/require"
)

func TestWorker(t *testing.T) {
	dom.WillTest(t, "test", true)

	t.Run("testPersist", func(t *testing.T) {
		testPersist(t)
	})
	t.Run("testVehiclePersist", func(t *testing.T) {
		testVehiclePersist(t)
	})
	t.Run("testWorkerInsertConstraint", func(t *testing.T) {
		testWorkerInsertConstraint(t)
	})
	t.Run("testGetWorker", func(t *testing.T) {
		testGetWorker(t)
	})
	t.Run("testQueryAll", func(t *testing.T) {
		testQueryAll(t)
	})
	t.Run("testWorkerStatus", func(t *testing.T) {
		testWorkerStatus(t)
	})
	t.Run("testRowsAffected", func(t *testing.T) {
		testRowsAffected(t)
	})
}

func testPersist(t *testing.T) {
	ctx := cfg.Context()
	for i := 0; i < 100; i++ {
		wrk := newWorker(randStatus())
		require.NoError(t, wrk.Insert(ctx))
	}
}

func testVehiclePersist(t *testing.T) {
	ctx := cfg.Context()
	for i := 0; i < 100; i++ {
		wv := newWorkerVehicle()
		require.NoError(t, wv.Insert(ctx))
	}
}

func testWorkerInsertConstraint(t *testing.T) {
	wrk := newWorker(randStatus())
	require.NoError(t, wrk.Insert(cfg.Context()))
	err := wrk.Insert(cfg.Context())
	require.True(t, dao.ErrUniqueViolation.Is(err))
}

func testGetWorker(t *testing.T) {
	ctx := cfg.Context()
	for i := 0; i < 100; i++ {
		wput := newWorker(randStatus())
		require.NoError(t, wput.Insert(ctx))
		wget, _, err := QueryOne(ctx, wput.ID)
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

func testQueryAll(t *testing.T) {
	ctx := cfg.Context()
	for i := 0; i < 20; i++ {
		worker := newWorker(randStatus())
		require.NoError(t, worker.Insert(ctx))
	}

	all0, err := QueryAll(ctx, nil, 10, 0)
	require.NoError(t, err)
	require.Len(t, all0, 10)

	all1, err := QueryAll(ctx, nil, 10, 10)
	require.NoError(t, err)
	require.Len(t, all1, 10)

	require.NotEqual(t, all0, all1)
}

func testWorkerStatus(t *testing.T) {
	dom.WillTest(t, "test", true)
	wrk := newWorker(randStatus())
	ctx := cfg.Context()
	require.NoError(t, wrk.Insert(ctx))
	testGetWorkerStatus(t, ctx, wrk)
	wrk.Status = randStatus()
	require.NoError(t, UpdateWorkerStatus(ctx, wrk.ID, wrk.Status))
	testGetWorkerStatus(t, ctx, wrk)
}

func testRowsAffected(t *testing.T) {
	dom.WillTest(t, "test", true)
	max := 100
	ctx := cfg.Context()
	for i := 0; i < max; i++ {
		wrk := newWorker(randStatus())
		require.NoError(t, wrk.Insert(ctx))
	}
	sql := dao.Update(workerTable).Set("updated", time.Now())
	require.NoError(t, dao.ExecTX(ctx, sql))
	require.Equal(t, max, int(sql.RowsAffected()))
}

func testGetWorkerStatus(t *testing.T, ctx dom.Context, wput *Worker) {
	status, _, er := QueryWorkerStatus(ctx, wput.ID)
	require.NoError(t, er)
	require.Equal(t, wput.Status, status)
}

func newWorker(status Status) *Worker {
	stamp := dom.Time{}
	stamp.Now()
	return &Worker{
		ID:        mockUUID("drv"),
		Status:    status,
		Rating:    int32(randInt(100)),
		Jobs:      randInt(1_000),
		FirstName: randFirstName(),
		LastName:  randLastName(),
		Vehicle:   mockUUID("car"),
		Created:   stamp,
		Updated:   stamp,
	}
}

func newWorkerVehicle() *WorkerVehicle {
	return &WorkerVehicle{
		Worker:  mockUUID("drv"),
		Vehicle: mockUUID("car"),
	}
}
