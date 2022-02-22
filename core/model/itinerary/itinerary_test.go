package itinerary

import (
	"fmt"
	"github.com/openmarketplaceengine/openmarketplaceengine/core/model/job"
	"testing"

	"github.com/openmarketplaceengine/openmarketplaceengine/core/model/step"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestItinerary(t *testing.T) {
	t.Run("testAddStep", func(t *testing.T) {
		testAddStep(t)
	})

	t.Run("testGetStep", func(t *testing.T) {
		testGetStep(t)
	})

	t.Run("testRemoveStep", func(t *testing.T) {
		testRemoveStep(t)
	})

	t.Run("testGetStepIndex", func(t *testing.T) {
		testGetStepIndex(t)
	})
}

func testAddStep(t *testing.T) {
	itinerary := NewItinerary("test-flow-1", []*job.Job{})

	_, err := itinerary.GetFirstStep()
	require.EqualError(t, err, "itinerary test-flow-1 has no steps")

	assert.Len(t, itinerary.Steps, 0)
	itinerary.AddStep(newStep("step1"))
	itinerary.AddStep(newStep("step2"))
	itinerary.AddStep(newStep("step3"))
	assert.Len(t, itinerary.Steps, 3)
}

func testGetStep(t *testing.T) {
	itinerary := NewItinerary("test-flow-1", []*job.Job{})

	step1 := newStep("step1")
	step2 := newStep("step2")

	itinerary.AddStep(step1)
	itinerary.AddStep(step2)

	assert.Equal(t, step2, itinerary.GetStep("step2"))
	firstStep, err := itinerary.GetFirstStep()
	require.NoError(t, err)
	assert.Equal(t, step1, firstStep)
}

func testRemoveStep(t *testing.T) {
	itinerary := NewItinerary("test-flow-1", []*job.Job{})

	step1 := newStep("step1")
	step2 := newStep("step2")
	step3 := newStep("step3")
	itinerary.AddStep(step1)
	itinerary.AddStep(step2)
	itinerary.AddStep(step3)

	assert.ElementsMatch(t, itinerary.Steps, []*step.Step{step1, step2, step3})

	itinerary.RemoveStep("step2")
	assert.ElementsMatch(t, itinerary.Steps, []*step.Step{step1, step3})

	itinerary.RemoveStep("step1")
	assert.ElementsMatch(t, itinerary.Steps, []*step.Step{step3})

	itinerary.RemoveStep("step3")
	assert.ElementsMatch(t, itinerary.Steps, []*step.Step{})
}

func testGetStepIndex(t *testing.T) {
	itinerary := NewItinerary("test-flow-1", []*job.Job{})

	step1 := newStep("step1")
	step2 := newStep("step2")
	step3 := newStep("step3")
	itinerary.AddStep(step1)
	itinerary.AddStep(step2)
	itinerary.AddStep(step3)

	i := itinerary.GetStepIndex("none")
	assert.Equal(t, -1, i)

	i1 := itinerary.GetStepIndex("step1")
	assert.Equal(t, 0, i1)

	i2 := itinerary.GetStepIndex("step2")
	assert.Equal(t, 1, i2)

	i3 := itinerary.GetStepIndex("step3")
	assert.Equal(t, 2, i3)
}

func (it *Itinerary) dump() {
	fmt.Printf("Itinerary: %s\n", it.ID)
	for i, s := range it.Steps {
		fmt.Printf("%v-%+v ->", i, s.ID)
	}

	fmt.Println()
}

func newStep(id step.ID) *step.Step {
	return &step.Step{ID: id}
}
