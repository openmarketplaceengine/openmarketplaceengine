package pickup

import (
	"fmt"
	"testing"

	"github.com/openmarketplaceengine/openmarketplaceengine/core/model/step"

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
	assert.NoError(t, checkTransition(Ready, CompleteAction))
	assert.NoError(t, checkTransition(Ready, CancelAction))
}

func testIllegalTransitions(t *testing.T) {
	assert.EqualError(t, checkTransition(Canceled, CompleteAction), "illegal transition from status=Canceled by action=CompleteAction")
	assert.EqualError(t, checkTransition(Canceled, CancelAction), "illegal transition from status=Canceled by action=CancelAction")
	assert.EqualError(t, checkTransition(Completed, CancelAction), "illegal transition from status=Completed by action=CancelAction")
	assert.EqualError(t, checkTransition(Completed, CompleteAction), "illegal transition from status=Completed by action=CompleteAction")
}

func checkTransition(current step.State, action step.Action) error {
	f := newFsm(current)
	ok := f.Event(actionToEvent[action])
	if !ok {
		return fmt.Errorf("illegal transition from status=%v by action=%v", current, action)
	}
	return nil
}
