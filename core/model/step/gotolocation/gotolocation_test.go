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

	t.Run("testMovingToNear", func(t *testing.T) {
		testMovingToNear(t)
	})

	t.Run("testTransitionError", func(t *testing.T) {
		testTransitionError(t)
	})
}

func testMovingToNear(t *testing.T) {
	ctx := context.Background()
	id := uuid.New().String()
	jobID := uuid.New().String()
	gtl, err := New(ctx, id, jobID)
	require.NoError(t, err)

	prevState := gtl.State
	prevUpdatedAt := gtl.UpdatedAt
	err = gtl.Handle(NearAction)
	require.NoError(t, err)

	require.NoError(t, err)
	assert.Equal(t, Moving, prevState)
	assert.Equal(t, gtl.StepID(), gtl.StepID())
	assert.Equal(t, Near, gtl.State)
	require.NoError(t, err)
	updatedAt0, _ := time.Parse(time.RFC3339Nano, prevUpdatedAt)
	updatedAt1, _ := time.Parse(time.RFC3339Nano, gtl.UpdatedAt)
	assert.Less(t, updatedAt0.UnixNano(), updatedAt1.UnixNano())
}

func testTransitionError(t *testing.T) {
	ctx := context.Background()
	id := uuid.New().String()
	jobID := uuid.New().String()
	gtl, err := New(ctx, id, jobID)
	require.NoError(t, err)

	err = gtl.Handle(NearAction)
	require.NoError(t, err)

	prevState := gtl.State
	prevUpdatedAt := gtl.UpdatedAt
	err = gtl.Handle(NearAction)
	require.Error(t, err)

	assert.Equal(t, prevState, gtl.State)
	updatedAt0, _ := time.Parse(time.RFC3339Nano, prevUpdatedAt)
	updatedAt1, _ := time.Parse(time.RFC3339Nano, gtl.UpdatedAt)
	assert.Equal(t, updatedAt0.UnixNano(), updatedAt1.UnixNano())
}
