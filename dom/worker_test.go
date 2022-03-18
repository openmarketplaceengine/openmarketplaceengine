package dom

import (
	"testing"

	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
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
