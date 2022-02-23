package gotolocation

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
	"github.com/openmarketplaceengine/openmarketplaceengine/core/model/step"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGoToLocation(t *testing.T) {
	err := cfg.Load()
	require.NoError(t, err)

	storage = newStorage(30 * time.Second)

	ctx := context.Background()
	id := step.ID(uuid.New().String())
	gtl, err := RetrieveOrCreate(ctx, id)
	require.NoError(t, err)

	t.Run("testRetrieveNil", func(t *testing.T) {
		testRetrieveNil(t)
	})

	t.Run("testNew", func(t *testing.T) {
		testNew(t, gtl)
	})

	t.Run("testNewToMoving", func(t *testing.T) {
		testNewToMoving(t, gtl)
	})

	t.Run("testNewToArrived", func(t *testing.T) {
		testNewToArrived(t, gtl)
	})
}

func testRetrieveNil(t *testing.T) {
	ctx := context.Background()
	id := step.ID(uuid.New().String())
	retrieved, err := storage.Retrieve(ctx, id)
	require.NoError(t, err)
	require.Nil(t, retrieved)
}

func testNew(t *testing.T, gtl *GoToLocation) {
	ctx := context.Background()

	retrieved, err := storage.Retrieve(ctx, step.ID(gtl.StepID))
	require.NoError(t, err)
	assert.Equal(t, New, gtl.Status)
	assert.Equal(t, gtl.Status, retrieved.Status)
	updatedAt0, _ := time.Parse(time.RFC3339Nano, gtl.UpdatedAt)
	updatedAt1, _ := time.Parse(time.RFC3339Nano, retrieved.UpdatedAt)
	assert.Equal(t, updatedAt0.UnixNano(), updatedAt1.UnixNano())
}

func testNewToMoving(t *testing.T, gtl *GoToLocation) {
	ctx := context.Background()
	prevState := gtl.Status
	prevUpdatedAt := gtl.UpdatedAt
	err := gtl.Handle(Move)
	require.NoError(t, err)

	movingGTL, err := storage.Retrieve(ctx, step.ID(gtl.StepID))
	require.NoError(t, err)
	assert.Equal(t, New, prevState)
	assert.Equal(t, gtl.StepID, movingGTL.StepID)
	assert.Equal(t, Moving, movingGTL.Status)
	require.NoError(t, err)
	updatedAt0, _ := time.Parse(time.RFC3339Nano, prevUpdatedAt)
	updatedAt1, _ := time.Parse(time.RFC3339Nano, movingGTL.UpdatedAt)
	assert.Less(t, updatedAt0.UnixNano(), updatedAt1.UnixNano())
}

func testNewToArrived(t *testing.T, gtl *GoToLocation) {
	ctx := context.Background()
	prevState := gtl.Status
	prevUpdatedAt := gtl.UpdatedAt
	err := gtl.Handle(Arrive)
	require.Error(t, err)

	retrieved, err := storage.Retrieve(ctx, step.ID(gtl.StepID))
	require.NoError(t, err)
	assert.Equal(t, prevState, retrieved.Status)
	updatedAt0, _ := time.Parse(time.RFC3339Nano, prevUpdatedAt)
	updatedAt1, _ := time.Parse(time.RFC3339Nano, retrieved.UpdatedAt)
	assert.Equal(t, updatedAt0.UnixNano(), updatedAt1.UnixNano())
}
