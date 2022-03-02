package itinerary

import (
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
	jobID := uuid.New().String()

	goToLocation1 := gotolocation.New(gotolocation.Moving)
	goToLocation2 := gotolocation.New(gotolocation.Moving)

	itinerary := New(uuid.New().String(), []*step.Step{})
	step1 := step.New(uuid.NewString(), jobID, goToLocation1, step.GoToLocation)
	step2 := step.New(uuid.NewString(), jobID, goToLocation2, step.GoToLocation)

	s, err := itinerary.CurrentStep()
	require.Error(t, err)
	require.Nil(t, s)

	assert.Len(t, itinerary.Steps, 0)

	itinerary.AddStep(step1)
	itinerary.AddStep(step2)

	assert.Len(t, itinerary.Steps, 2)
	currentStep, err := itinerary.CurrentStep()
	require.NoError(t, err)
	assert.Equal(t, step1, currentStep)
}

func testRemoveStep(t *testing.T) {
	jobID := uuid.New().String()

	goToLocation1 := gotolocation.New(gotolocation.Moving)
	goToLocation2 := gotolocation.New(gotolocation.Moving)

	itinerary := New(uuid.New().String(), []*step.Step{})

	step1 := step.New(uuid.NewString(), jobID, goToLocation1, step.GoToLocation)
	itinerary.AddStep(step1)
	step2 := step.New(uuid.NewString(), jobID, goToLocation2, step.GoToLocation)
	itinerary.AddStep(step2)

	assert.ElementsMatch(t, itinerary.Steps, []*step.Step{step1, step2})

	itinerary.RemoveStep(step2.StepID())
	assert.ElementsMatch(t, itinerary.Steps, []*step.Step{step1})

	itinerary.RemoveStep(step1.StepID())
	assert.ElementsMatch(t, itinerary.Steps, []*step.Step{})
}

func testStepIndex(t *testing.T) {
	jobID := uuid.NewString()

	goToLocation1 := gotolocation.New(gotolocation.Moving)
	goToLocation2 := gotolocation.New(gotolocation.Moving)

	itinerary := New(uuid.New().String(), []*step.Step{})

	step1 := step.New(uuid.NewString(), jobID, goToLocation1, step.GoToLocation)
	itinerary.AddStep(step1)
	step2 := step.New(uuid.NewString(), jobID, goToLocation2, step.GoToLocation)
	itinerary.AddStep(step2)

	i := itinerary.stepIndex("none")
	assert.Equal(t, -1, i)

	i1 := itinerary.stepIndex(step1.StepID())
	assert.Equal(t, 0, i1)

	i2 := itinerary.stepIndex(step2.StepID())
	assert.Equal(t, 1, i2)
}

func (it *Itinerary) dump() {
	fmt.Printf("Itinerary: %s\n", it.ID)
	for i, s := range it.Steps {
		fmt.Printf("%v-%+v ->", i, s.StepID())
	}

	fmt.Println()
}
