package pickup

import (
	"testing"

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
	pickup := New(Ready)

	prevState := pickup.CurrentState()
	err := pickup.Handle(Complete)
	require.NoError(t, err)
	assert.Equal(t, Ready, prevState)
	assert.Equal(t, Completed, pickup.CurrentState())
}

func testTransitionError(t *testing.T) {
	pickup := New(Ready)
	err := pickup.Handle(Complete)
	require.NoError(t, err)

	prevState := pickup.CurrentState()
	assert.Equal(t, Completed, prevState)
	err = pickup.Handle(Cancel)
	require.EqualError(t, err, "illegal transition from state=Completed by event=Cancel")
	assert.Equal(t, prevState, pickup.CurrentState())
}
