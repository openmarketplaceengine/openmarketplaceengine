package provider

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

	t.Run("testInitialState", func(t *testing.T) {
		testInitialState(t)
	})

	t.Run("testAllTransitions", func(t *testing.T) {
		testAllTransitions(t, state)
	})
}

func testGoOnline(t *testing.T, state *StateController) {
	state.fcm.Reset()
	err := state.GoOnline()
	require.NoError(t, err)
	require.Equal(t, Online, State(state.fcm.Current()))

	err = state.GoOnline()
	assert.Error(t, err)
}

func testGoOffline(t *testing.T, state *StateController) {
	state.fcm.Reset()
	err := state.GoOffline()
	require.Error(t, err)
	require.Equal(t, Offline, State(state.fcm.Current()))

	err = state.GoOnline()
	assert.NoError(t, err)

	err = state.GoOffline()
	require.NoError(t, err)
}

func testInitialState(t *testing.T) {
	state := NewStateController(Delivering)
	err := state.GoOffline()
	require.Error(t, err)
	require.Equal(t, Delivering, State(state.fcm.Current()))

	err = state.CompleteDelivery()
	assert.NoError(t, err)
	require.Equal(t, Online, State(state.fcm.Current()))
}

func testAllTransitions(t *testing.T, state *StateController) {
	state.fcm.Reset()
	err := state.GoOnline()
	assert.NoError(t, err)
	require.Equal(t, Online, State(state.fcm.Current()))

	err = state.PickUp()
	require.NoError(t, err)
	require.Equal(t, PickingUp, State(state.fcm.Current()))

	err = state.Deliver()
	require.NoError(t, err)
	require.Equal(t, Delivering, State(state.fcm.Current()))

	err = state.CompleteDelivery()
	require.NoError(t, err)
	require.Equal(t, Online, State(state.fcm.Current()))

	err = state.GoOffline()
	require.NoError(t, err)
	require.Equal(t, Offline, State(state.fcm.Current()))
}
