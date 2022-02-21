package gotolocation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGTLTransitions(t *testing.T) {
	assert.NoError(t, checkTransition(New, MoveEvent))
	assert.NoError(t, checkTransition(Moving, NearEvent))
	assert.NoError(t, checkTransition(Near, ArrivedEvent))

	assert.NoError(t, checkTransition(New, CancelledEvent))
	assert.NoError(t, checkTransition(Moving, CancelledEvent))
	assert.NoError(t, checkTransition(Near, CancelledEvent))
	assert.NoError(t, checkTransition(Arrived, CancelledEvent))
}

func TestGTLIllegalTransitions(t *testing.T) {
	assert.EqualError(t, checkTransition(New, NearEvent), "illegal transition from state=New by event=NearEvent")
	assert.EqualError(t, checkTransition(New, ArrivedEvent), "illegal transition from state=New by event=ArrivedEvent")
	assert.EqualError(t, checkTransition(Moving, ArrivedEvent), "illegal transition from state=Moving by event=ArrivedEvent")
	assert.EqualError(t, checkTransition(Arrived, MoveEvent), "illegal transition from state=Arrived by event=MoveEvent")
}
