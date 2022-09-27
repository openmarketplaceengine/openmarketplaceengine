package gotolocation

import (
	"testing"

	"github.com/driverscooperative/geosrv/cfg"
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
	gtl := New(Moving)

	prevState := gtl.CurrentState()
	err := gtl.Handle(NearBy)
	require.NoError(t, err)

	require.NoError(t, err)
	assert.Equal(t, Moving, prevState)
	assert.Equal(t, Near, gtl.CurrentState())
	require.NoError(t, err)
}

func testTransitionError(t *testing.T) {
	gtl := New(Moving)
	err := gtl.Handle(NearBy)
	require.NoError(t, err)

	prevState := gtl.CurrentState()
	err = gtl.Handle(NearBy)
	require.EqualError(t, err, "illegal transition from state=1 by event=0")

	assert.Equal(t, prevState, gtl.CurrentState())
}
