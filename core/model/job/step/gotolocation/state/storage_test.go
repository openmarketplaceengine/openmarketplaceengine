package state

import (
	"github.com/google/uuid"
	"github.com/openmarketplaceengine/openmarketplaceengine/core/model/job/step/gotolocation/fsm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"testing"
	"time"
)

func TestStorage(t *testing.T) {

	storage := NewStorage(1 * time.Second)

	t.Run("testStoreAndRetrieve", func(t *testing.T) {
		testStoreAndRetrieve(t, storage)
	})
}

func testStoreAndRetrieve(t *testing.T, storage *Storage) {
	ctx := context.Background()
	key := uuid.New().String()
	stateIn := State{
		DriverID:                "driver-1",
		PassengerIDs:            []string{"pass-1", "pass-2"},
		DestinationLatitude:     7,
		DestinationLongitude:    8,
		CreatedAt:               time.Now(),
		LastModifiedAt:          time.Now(),
		LastModifiedAtLatitude:  9,
		LastModifiedAtLongitude: 10,
		LastEvent:               fsm.Near,
		LastState:               fsm.NearState,
	}
	err := storage.Store(ctx, key, stateIn)

	require.NoError(t, err)

	stateOut, err := storage.Retrieve(ctx, key)
	require.NoError(t, err)
	assert.Equal(t, stateIn.DriverID, stateOut.DriverID)
	assert.Equal(t, stateIn.PassengerIDs, stateOut.PassengerIDs)
	assert.ElementsMatch(t, stateIn.PassengerIDs, stateOut.PassengerIDs)
	assert.Equal(t, stateIn.DestinationLatitude, stateOut.DestinationLatitude)
	assert.Equal(t, stateIn.DestinationLongitude, stateOut.DestinationLongitude)
	assert.Equal(t, stateIn.LastEvent, stateOut.LastEvent)
	assert.Equal(t, stateIn.LastState, stateOut.LastState)
	assert.Equal(t, stateIn.CreatedAt.UnixMilli(), stateOut.CreatedAt.UnixMilli())
	assert.Equal(t, stateIn.LastModifiedAt.UnixMilli(), stateOut.LastModifiedAt.UnixMilli())
}
