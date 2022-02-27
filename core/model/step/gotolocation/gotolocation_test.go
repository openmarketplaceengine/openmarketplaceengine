package gotolocation

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGoToLocation(t *testing.T) {
	err := cfg.Load()
	require.NoError(t, err)

	storage = newStorage(30 * time.Second)

	ctx := context.Background()
	id := uuid.New().String()
	jobID := uuid.New().String()
	gtl, err := New(ctx, id, jobID)
	require.NoError(t, err)

	t.Run("testRetrieveNil", func(t *testing.T) {
		testRetrieveNil(t)
	})

	t.Run("testStorage", func(t *testing.T) {
		testStorage(t, gtl)
	})

	t.Run("testMovingToNear", func(t *testing.T) {
		testMovingToNear(t, gtl)
	})

	t.Run("testTransitionError", func(t *testing.T) {
		testTransitionError(t, gtl)
	})
}

func testRetrieveNil(t *testing.T) {
	ctx := context.Background()
	id := uuid.New().String()
	retrieved, err := storage.Retrieve(ctx, id)
	require.NoError(t, err)
	require.Nil(t, retrieved)
}

func testStorage(t *testing.T, gtl *GoToLocation) {
	ctx := context.Background()

	retrieved, err := storage.Retrieve(ctx, gtl.StepID())
	require.NoError(t, err)
	assert.Equal(t, Moving, gtl.State)
	assert.Equal(t, gtl.State, retrieved.State)
	updatedAt0, _ := time.Parse(time.RFC3339Nano, gtl.UpdatedAt)
	updatedAt1, _ := time.Parse(time.RFC3339Nano, retrieved.UpdatedAt)
	assert.Equal(t, updatedAt0.UnixNano(), updatedAt1.UnixNano())
}

func testMovingToNear(t *testing.T, gtl *GoToLocation) {
	ctx := context.Background()
	prevState := gtl.State
	prevUpdatedAt := gtl.UpdatedAt
	err := gtl.Handle(NearAction)
	require.NoError(t, err)

	movingGTL, err := storage.Retrieve(ctx, gtl.StepID())
	require.NoError(t, err)
	assert.Equal(t, Moving, prevState)
	assert.Equal(t, gtl.StepID(), movingGTL.StepID())
	assert.Equal(t, Near, movingGTL.State)
	require.NoError(t, err)
	updatedAt0, _ := time.Parse(time.RFC3339Nano, prevUpdatedAt)
	updatedAt1, _ := time.Parse(time.RFC3339Nano, movingGTL.UpdatedAt)
	assert.Less(t, updatedAt0.UnixNano(), updatedAt1.UnixNano())
}

func testTransitionError(t *testing.T, gtl *GoToLocation) {
	ctx := context.Background()
	err := gtl.Handle(ArriveAction)
	require.NoError(t, err)

	prevState := gtl.State
	prevUpdatedAt := gtl.UpdatedAt

	err = gtl.Handle(ArriveAction)
	require.Error(t, err)

	retrieved, err := storage.Retrieve(ctx, gtl.StepID())
	require.NoError(t, err)
	assert.Equal(t, prevState, retrieved.State)
	updatedAt0, _ := time.Parse(time.RFC3339Nano, prevUpdatedAt)
	updatedAt1, _ := time.Parse(time.RFC3339Nano, retrieved.UpdatedAt)
	assert.Equal(t, updatedAt0.UnixNano(), updatedAt1.UnixNano())
}
