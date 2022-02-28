package pickup

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPickup(t *testing.T) {
	t.Run("testReadyToCompleted", func(t *testing.T) {
		testReadyToCompleted(t)
	})

	t.Run("testTransitionError", func(t *testing.T) {
		testTransitionError(t)
	})
}

func testReadyToCompleted(t *testing.T) {
	ctx := context.Background()
	id := uuid.New().String()
	jobID := uuid.New().String()
	pickup, err := New(ctx, id, jobID)
	require.NoError(t, err)

	prevState := pickup.State
	prevUpdatedAt := pickup.UpdatedAt
	err = pickup.Handle(CompleteAction)
	require.NoError(t, err)
	assert.Equal(t, Ready, prevState)
	assert.Equal(t, Completed, pickup.State)
	prev, err := time.Parse(time.RFC3339Nano, prevUpdatedAt)
	require.NoError(t, err)
	last, err := time.Parse(time.RFC3339Nano, pickup.UpdatedAt)
	require.NoError(t, err)
	assert.Less(t, prev.UnixNano(), last.UnixNano())
}

func testTransitionError(t *testing.T) {
	ctx := context.Background()
	id := uuid.New().String()
	jobID := uuid.New().String()
	pickup, err := New(ctx, id, jobID)
	require.NoError(t, err)

	err = pickup.Handle(CompleteAction)
	require.NoError(t, err)

	prevState := pickup.State
	prevUpdatedAt := pickup.UpdatedAt
	assert.Equal(t, Completed, pickup.State)
	err = pickup.Handle(CancelAction)
	require.Error(t, err)
	assert.Equal(t, prevState, pickup.State)
	assert.Equal(t, prevUpdatedAt, pickup.UpdatedAt)
}
