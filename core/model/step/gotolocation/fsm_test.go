package gotolocation

import (
	"fmt"
	"testing"

	"github.com/openmarketplaceengine/openmarketplaceengine/core/model/step"
	"github.com/stretchr/testify/assert"
)

func TestGTLTransitions(t *testing.T) {
	assert.NoError(t, checkTransition(New, Move))
	assert.NoError(t, checkTransition(Moving, Near))
	assert.NoError(t, checkTransition(NearBy, Arrive))

	assert.NoError(t, checkTransition(New, Cancel))
	assert.NoError(t, checkTransition(Moving, Cancel))
	assert.NoError(t, checkTransition(NearBy, Cancel))
	assert.NoError(t, checkTransition(Arrived, Cancel))
}

func TestGTLIllegalTransitions(t *testing.T) {
	assert.EqualError(t, checkTransition(New, Near), "illegal transition from status=New by action=Near")
	assert.EqualError(t, checkTransition(New, Arrive), "illegal transition from status=New by action=Arrive")
	assert.EqualError(t, checkTransition(Moving, Arrive), "illegal transition from status=Moving by action=Arrive")
	assert.EqualError(t, checkTransition(Arrived, Move), "illegal transition from status=Arrived by action=Move")
}

func checkTransition(current step.Status, action step.Action) error {
	f := newFsm(current)
	ok := f.Event(actionToEvent[action])
	if !ok {
		return fmt.Errorf("illegal transition from status=%v by action=%v", current, action)
	}
	return nil
}
