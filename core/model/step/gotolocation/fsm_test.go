package gotolocation

import (
	"fmt"
	"testing"

	"github.com/openmarketplaceengine/openmarketplaceengine/core/model/step"
	"github.com/stretchr/testify/assert"
)

func TestGTLTransitions(t *testing.T) {
	assert.NoError(t, checkTransition(Moving, NearAction))
	assert.NoError(t, checkTransition(Near, ArriveAction))

	assert.NoError(t, checkTransition(Moving, CancelAction))
	assert.NoError(t, checkTransition(Near, CancelAction))
}

func TestGTLIllegalTransitions(t *testing.T) {
	assert.EqualError(t, checkTransition(Canceled, NearAction), "illegal transition from status=Canceled by action=NearAction")
	assert.EqualError(t, checkTransition(Moving, ArriveAction), "illegal transition from status=Moving by action=ArriveAction")
	assert.EqualError(t, checkTransition(Arrived, CancelAction), "illegal transition from status=Arrived by action=CancelAction")
}

func checkTransition(current step.State, action step.Action) error {
	f := newFsm(current)
	ok := f.Event(actionToEvent[action])
	if !ok {
		return fmt.Errorf("illegal transition from status=%v by action=%v", current, action)
	}
	return nil
}
