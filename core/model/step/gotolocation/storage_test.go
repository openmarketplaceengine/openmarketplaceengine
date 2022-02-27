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
	id := uuid.New().String()
	stateIn := GoToLocation{
		ID:        id,
		UpdatedAt: time.Now().Format(time.RFC3339Nano),
		State:     Moving,
	}
	err := storage.Store(ctx, stateIn)

	require.NoError(t, err)

	stateOut, err := storage.Retrieve(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, stateIn.StepID(), stateOut.StepID())
	assert.Equal(t, stateIn.State, stateOut.State)
	updatedAtIn, _ := time.Parse(time.RFC3339Nano, stateIn.UpdatedAt)
	updatedAtOut, _ := time.Parse(time.RFC3339Nano, stateOut.UpdatedAt)
	assert.Equal(t, updatedAtIn.UnixNano(), updatedAtOut.UnixNano())
}
