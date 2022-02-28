package itinerary

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
	"github.com/openmarketplaceengine/openmarketplaceengine/core/model/step/gotolocation"

	"github.com/openmarketplaceengine/openmarketplaceengine/core/model/step"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestItinerary(t *testing.T) {
	err := cfg.Load()
	require.NoError(t, err)

	t.Run("testAddStep", func(t *testing.T) {
		testAddStep(t)
	})

	t.Run("testRemoveStep", func(t *testing.T) {
		testRemoveStep(t)
	})

	t.Run("testStepIndex", func(t *testing.T) {
		testStepIndex(t)
	})
}

func testAddStep(t *testing.T) {
	ctx := context.Background()
	jobID := uuid.New().String()

	step1, err := gotolocation.New(ctx, uuid.New().String(), jobID)
	require.NoError(t, err)
	step2, err := gotolocation.New(ctx, uuid.New().String(), jobID)
	require.NoError(t, err)
	step3, err := gotolocation.New(ctx, uuid.New().String(), jobID)
	require.NoError(t, err)

	itinerary := New(uuid.New().String(), []step.Step{})

	s, err := itinerary.CurrentStep()
	require.Error(t, err)
	require.Nil(t, s)

	assert.Len(t, itinerary.Steps, 0)

	itinerary.AddStep(step1)
	itinerary.AddStep(step2)
	itinerary.AddStep(step3)

	assert.Len(t, itinerary.Steps, 3)
	currentStep, err := itinerary.CurrentStep()
	require.NoError(t, err)
	assert.Equal(t, step1, currentStep)
}

func testRemoveStep(t *testing.T) {
	ctx := context.Background()
	jobID := uuid.New().String()

	step1, err := gotolocation.New(ctx, uuid.New().String(), jobID)
	require.NoError(t, err)
	step2, err := gotolocation.New(ctx, uuid.New().String(), jobID)
	require.NoError(t, err)
	step3, err := gotolocation.New(ctx, uuid.New().String(), jobID)
	require.NoError(t, err)
	itinerary := New(uuid.New().String(), []step.Step{step1, step2, step3})

	assert.ElementsMatch(t, itinerary.Steps, []step.Step{step1, step2, step3})

	itinerary.RemoveStep(step2.StepID())
	assert.ElementsMatch(t, itinerary.Steps, []step.Step{step1, step3})

	itinerary.RemoveStep(step1.StepID())
	assert.ElementsMatch(t, itinerary.Steps, []step.Step{step3})

	itinerary.RemoveStep(step3.StepID())
	assert.ElementsMatch(t, itinerary.Steps, []step.Step{})
}

func testStepIndex(t *testing.T) {
	ctx := context.Background()
	jobID := uuid.New().String()

	step1, err := gotolocation.New(ctx, uuid.New().String(), jobID)
	require.NoError(t, err)
	step2, err := gotolocation.New(ctx, uuid.New().String(), jobID)
	require.NoError(t, err)
	step3, err := gotolocation.New(ctx, uuid.New().String(), jobID)
	require.NoError(t, err)

	itinerary := New(uuid.New().String(), []step.Step{})

	itinerary.AddStep(step1)
	itinerary.AddStep(step2)
	itinerary.AddStep(step3)

	i := itinerary.stepIndex("none")
	assert.Equal(t, -1, i)

	i1 := itinerary.stepIndex(step1.StepID())
	assert.Equal(t, 0, i1)

	i2 := itinerary.stepIndex(step2.StepID())
	assert.Equal(t, 1, i2)

	i3 := itinerary.stepIndex(step3.StepID())
	assert.Equal(t, 2, i3)
}

func (it *Itinerary) dump() {
	fmt.Printf("Itinerary: %s\n", it.ID)
	for i, s := range it.Steps {
		fmt.Printf("%v-%+v ->", i, s.StepID())
	}

	fmt.Println()
}
