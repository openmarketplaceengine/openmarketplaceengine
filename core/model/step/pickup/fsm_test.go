package pickup

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPickupFSM(t *testing.T) {
	t.Run("testTransitions", func(t *testing.T) {
		testTransitions(t)
	})
	t.Run("testIllegalTransitions", func(t *testing.T) {
		testIllegalTransitions(t)
	})
}

func testTransitions(t *testing.T) {
	assert.NoError(t, checkTransition(New, ReadyEvent))
	assert.NoError(t, checkTransition(Ready, CompletedEvent))
	assert.NoError(t, checkTransition(Ready, CancelledEvent))
}

func testIllegalTransitions(t *testing.T) {
	assert.EqualError(t, checkTransition(New, CompletedEvent), "illegal transition from state=New by event=CompletedEvent")
	assert.EqualError(t, checkTransition(Completed, CancelledEvent), "illegal transition from state=Completed by event=CancelledEvent")
}
