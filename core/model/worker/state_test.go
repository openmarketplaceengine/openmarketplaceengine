package worker

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProviderState(t *testing.T) {
	state := NewStateController(Offline)

	t.Run("testGoOnline", func(t *testing.T) {
		testGoOnline(t, state)
	})

	t.Run("testGoOffline", func(t *testing.T) {
		testGoOffline(t, state)
	})

	t.Run("testCannotSignOffWhileDelivering", func(t *testing.T) {
		testCannotSignOffWhileDelivering(t)
	})

	t.Run("testAllTransitions", func(t *testing.T) {
		testAllTransitions(t, state)
	})
}

func testGoOnline(t *testing.T, state *StateController) {
	state.fcm.Reset()
	err := state.Ready()
	require.NoError(t, err)
	require.Equal(t, Idle, State(state.fcm.Current()))

	err = state.Ready()
	assert.Error(t, err)
}

func testGoOffline(t *testing.T, state *StateController) {
	state.fcm.Reset()
	err := state.SignOff()
	require.Error(t, err)
	require.Equal(t, Offline, State(state.fcm.Current()))

	err = state.Ready()
	assert.NoError(t, err)

	err = state.SignOff()
	require.NoError(t, err)
}

func testCannotSignOffWhileDelivering(t *testing.T) {
	state := NewStateController(Delivering)
	err := state.SignOff()
	require.Error(t, err)
	require.Equal(t, Delivering, State(state.fcm.Current()))

	err = state.DropOff()
	assert.NoError(t, err)
	require.Equal(t, DroppingOff, State(state.fcm.Current()))

	err = state.Ready()
	assert.NoError(t, err)
	require.Equal(t, Idle, State(state.fcm.Current()))
}

func testAllTransitions(t *testing.T, state *StateController) {
	state.fcm.Reset()
	err := state.Ready()
	assert.NoError(t, err)
	require.Equal(t, Idle, State(state.fcm.Current()))

	err = state.PickUp()
	require.NoError(t, err)
	require.Equal(t, PickingUp, State(state.fcm.Current()))

	err = state.Deliver()
	require.NoError(t, err)
	require.Equal(t, Delivering, State(state.fcm.Current()))

	err = state.DropOff()
	require.NoError(t, err)
	require.Equal(t, DroppingOff, State(state.fcm.Current()))

	err = state.Ready()
	assert.NoError(t, err)
	require.Equal(t, Idle, State(state.fcm.Current()))

	err = state.SignOff()
	require.NoError(t, err)
	require.Equal(t, Offline, State(state.fcm.Current()))
}
