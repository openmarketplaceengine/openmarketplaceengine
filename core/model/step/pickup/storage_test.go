package pickup

import (
	"context"
	"github.com/google/uuid"
	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestStorage(t *testing.T) {
	err := cfg.Load()
	require.NoError(t, err)
	storage := newStorage(1 * time.Second)

	t.Run("testStoreAndRetrieve", func(t *testing.T) {
		testStoreAndRetrieve(t, storage)
	})
}

func testStoreAndRetrieve(t *testing.T, storage *Storage) {
	ctx := context.Background()
	driverID := uuid.New().String()
	stateIn := Pickup{
		DriverID:        driverID,
		PickupLatitude:  3,
		PickupLongitude: 4,
		PassengerIds:    "s1,s2",
		UpdatedAt:       time.Now().Format(time.RFC3339),
		State:           Ready,
	}
	err := storage.Store(ctx, stateIn)

	require.NoError(t, err)

	stateOut, err := storage.Retrieve(ctx, driverID)
	require.NoError(t, err)
	assert.Equal(t, stateIn.DriverID, stateOut.DriverID)
	assert.Equal(t, stateIn.PickupLatitude, stateOut.PickupLatitude)
	assert.Equal(t, stateIn.PickupLongitude, stateOut.PickupLongitude)
	assert.Equal(t, stateIn.State, stateOut.State)
	assert.Equal(t, stateIn.UpdatedAt, stateOut.UpdatedAt)
	assert.Equal(t, stateIn.PassengerIds, stateOut.PassengerIds)
}
