package dropoff

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDropOff(t *testing.T) {
	err := cfg.Load()
	require.NoError(t, err)

	storage = newStorage(30 * time.Second)

	ctx := context.Background()
	driverID := uuid.New().String()
	require.NoError(t, err)

	t.Run("testRetrieveNil", func(t *testing.T) {
		testRetrieveNil(t)
	})

	t.Run("testNew", func(t *testing.T) {
		testNew(ctx, t, driverID)
	})

	t.Run("testNewToReady", func(t *testing.T) {
		testNewToReady(ctx, t, driverID)
	})

	t.Run("testNewToCancelled", func(t *testing.T) {
		testNewToCancelled(ctx, t, driverID)
	})
}

func testRetrieveNil(t *testing.T) {
	ctx := context.Background()
	driverID := uuid.New().String()
	_, err := storage.Retrieve(ctx, driverID)
	require.Error(t, err)
}

func testNew(ctx context.Context, t *testing.T, driverID string) {
	dropOff, err := NewDropOff(ctx, driverID, 7, 8)
	require.NoError(t, err)

	retrieved, err := storage.Retrieve(ctx, dropOff.DriverID)
	require.NoError(t, err)
	assert.Equal(t, New, dropOff.State)
	assert.Equal(t, dropOff.DriverID, retrieved.DriverID)
	assert.Equal(t, dropOff.DropOffLatitude, retrieved.DropOffLatitude)
	assert.Equal(t, dropOff.DropOffLongitude, retrieved.DropOffLongitude)
	assert.Equal(t, dropOff.State, retrieved.State)
	assert.Equal(t, dropOff.UpdatedAt, retrieved.UpdatedAt)
}

func testNewToReady(ctx context.Context, t *testing.T, driverID string) {
	dropOff, err := NewDropOff(ctx, driverID, 7, 8)
	require.NoError(t, err)

	prevState := dropOff.State
	prevUpdatedAt := dropOff.UpdatedAt
	err = dropOff.Ready(ctx)
	require.NoError(t, err)

	ready, err := storage.Retrieve(ctx, dropOff.DriverID)
	require.NoError(t, err)
	assert.Equal(t, New, prevState)
	assert.Equal(t, dropOff.DriverID, ready.DriverID)
	assert.Equal(t, Ready, ready.State)
	prev, err := time.Parse(time.RFC3339Nano, prevUpdatedAt)
	require.NoError(t, err)
	last, err := time.Parse(time.RFC3339Nano, ready.UpdatedAt)
	require.NoError(t, err)
	assert.Less(t, prev.UnixNano(), last.UnixNano())
}

func testNewToCancelled(ctx context.Context, t *testing.T, driverID string) {
	dropOff, err := NewDropOff(ctx, driverID, 7, 8)
	require.NoError(t, err)

	prevState := dropOff.State
	prevUpdatedAt := dropOff.UpdatedAt
	assert.Equal(t, New, dropOff.State)
	err = dropOff.Cancel(ctx)
	require.Error(t, err)

	retrieved, err := storage.Retrieve(ctx, dropOff.DriverID)
	require.NoError(t, err)
	assert.Equal(t, prevState, retrieved.State)
	assert.Equal(t, prevUpdatedAt, retrieved.UpdatedAt)
}
