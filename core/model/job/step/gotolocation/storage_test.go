package gotolocation

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
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
	stateIn := GoToLocation{
		DriverID:             driverID,
		DestinationLatitude:  7,
		DestinationLongitude: 8,
		UpdatedAt:            time.Now().Format(time.RFC3339),
		UpdatedAtLatitude:    9,
		UpdatedAtLongitude:   10,
		State:                Moving,
	}
	err := storage.Store(ctx, stateIn)

	require.NoError(t, err)

	stateOut, err := storage.Retrieve(ctx, driverID)
	require.NoError(t, err)
	assert.Equal(t, stateIn.DriverID, stateOut.DriverID)
	assert.Equal(t, stateIn.DestinationLatitude, stateOut.DestinationLatitude)
	assert.Equal(t, stateIn.DestinationLongitude, stateOut.DestinationLongitude)
	assert.Equal(t, stateIn.State, stateOut.State)
	assert.Equal(t, stateIn.UpdatedAt, stateOut.UpdatedAt)
}
