package gotolocation

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"testing"
	"time"
)

func TestStorage(t *testing.T) {

	storage := newStorage(1 * time.Second)

	t.Run("testStoreAndRetrieve", func(t *testing.T) {
		testStoreAndRetrieve(t, storage)
	})
}

func testStoreAndRetrieve(t *testing.T, storage *Storage) {
	ctx := context.Background()
	key := uuid.New().String()
	stateIn := GoToLocation{
		DriverID:                "driver-1",
		DestinationLatitude:     7,
		DestinationLongitude:    8,
		LastModifiedAt:          time.Now(),
		LastModifiedAtLatitude:  9,
		LastModifiedAtLongitude: 10,
		State:                   Moving,
	}
	err := storage.Store(ctx, key, stateIn)

	require.NoError(t, err)

	stateOut, err := storage.Retrieve(ctx, key)
	require.NoError(t, err)
	assert.Equal(t, stateIn.DriverID, stateOut.DriverID)
	assert.Equal(t, stateIn.DestinationLatitude, stateOut.DestinationLatitude)
	assert.Equal(t, stateIn.DestinationLongitude, stateOut.DestinationLongitude)
	assert.Equal(t, stateIn.State, stateOut.State)
	assert.Equal(t, stateIn.LastModifiedAt.UnixMilli(), stateOut.LastModifiedAt.UnixMilli())
}
