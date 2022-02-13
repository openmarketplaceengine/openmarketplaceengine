package gotolocation

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
)

func TestGoToLocation(t *testing.T) {
	storage = newStorage(30 * time.Second)

	ctx := context.Background()
	driverID := uuid.New().String()
	newGTL, err := NewGoToLocation(ctx, driverID, 7, 8)
	require.NoError(t, err)

	t.Run("testRetrieveNil", func(t *testing.T) {
		testRetrieveNil(t)
	})

	t.Run("testNew", func(t *testing.T) {
		testNew(t, newGTL)
	})

	t.Run("testNewToMoving", func(t *testing.T) {
		testNewToMoving(t, newGTL)
	})

	t.Run("testNewToArrived", func(t *testing.T) {
		testNewToArrived(t, newGTL)
	})
}

func testRetrieveNil(t *testing.T) {
	ctx := context.Background()
	driverID := uuid.New().String()
	_, err := storage.Retrieve(ctx, driverID)
	require.Error(t, err)
}

func testNew(t *testing.T, newGTL *GoToLocation) {
	ctx := context.Background()

	retrieved, err := storage.Retrieve(ctx, newGTL.DriverID)
	require.NoError(t, err)
	assert.Equal(t, New, newGTL.State)
	assert.Equal(t, newGTL.DriverID, retrieved.DriverID)
	assert.Equal(t, newGTL.DestinationLatitude, retrieved.DestinationLatitude)
	assert.Equal(t, newGTL.DestinationLongitude, retrieved.DestinationLongitude)
	assert.Equal(t, newGTL.State, retrieved.State)
	assert.Equal(t, newGTL.LastModifiedAt, retrieved.LastModifiedAt)
}

func testNewToMoving(t *testing.T, newGTL *GoToLocation) {
	ctx := context.Background()
	prevState := newGTL.State
	prevLastModifiedAt := newGTL.LastModifiedAt
	err := newGTL.Moving(ctx, 7, 8)
	require.NoError(t, err)

	movingGTL, err := storage.Retrieve(ctx, newGTL.DriverID)
	require.NoError(t, err)
	assert.Equal(t, New, prevState)
	assert.Equal(t, newGTL.DriverID, movingGTL.DriverID)
	assert.Equal(t, Moving, movingGTL.State)
	prev, err := time.Parse(time.RFC3339Nano, prevLastModifiedAt)
	require.NoError(t, err)
	last, err := time.Parse(time.RFC3339Nano, movingGTL.LastModifiedAt)
	require.NoError(t, err)
	assert.Less(t, prev.UnixMilli(), last.UnixMilli())
}

func testNewToArrived(t *testing.T, newGTL *GoToLocation) {
	ctx := context.Background()
	prevState := newGTL.State
	prevLastModifiedAt := newGTL.LastModifiedAt
	err := newGTL.Arrived(ctx, 7, 8)
	require.Error(t, err)

	retrieved, err := storage.Retrieve(ctx, newGTL.DriverID)
	require.NoError(t, err)
	assert.Equal(t, prevState, retrieved.State)
	assert.Equal(t, prevLastModifiedAt, retrieved.LastModifiedAt)
}
