package flow

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFlow(t *testing.T) {
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
	f := NewFlow("test-flow-1", []Step{})

	_, err := f.GetFirstStep()
	require.EqualError(t, err, "flow test-flow-1 has no steps")

	assert.Len(t, f.Steps, 0)
	f.AddStep(newStep("step1"))
	f.AddStep(newStep("step2"))
	f.AddStep(newStep("step3"))
	assert.Len(t, f.Steps, 3)
}

func testGetStep(t *testing.T) {
	f := NewFlow("test-flow-1", []Step{})

	step1 := newStep("step1")
	step2 := newStep("step2")

	f.AddStep(step1)
	f.AddStep(step2)

	assert.Equal(t, step2, f.GetStep("step2"))
	step, err := f.GetFirstStep()
	require.NoError(t, err)
	assert.Equal(t, step1, step)
}

func testRemoveStep(t *testing.T) {
	f := NewFlow("test-flow-1", []Step{})

	step1 := newStep("step1")
	step2 := newStep("step2")
	step3 := newStep("step3")
	f.AddStep(step1)
	f.AddStep(step2)
	f.AddStep(step3)

	assert.ElementsMatch(t, f.Steps, []Step{step1, step2, step3})

	f.RemoveStep("step2")
	assert.ElementsMatch(t, f.Steps, []Step{step1, step3})

	f.RemoveStep("step1")
	assert.ElementsMatch(t, f.Steps, []Step{step3})

	f.RemoveStep("step3")
	assert.ElementsMatch(t, f.Steps, []Step{})
}

func testGetStepIndex(t *testing.T) {
	f := NewFlow("test-flow-1", []Step{})

	step1 := newStep("step1")
	step2 := newStep("step2")
	step3 := newStep("step3")
	f.AddStep(step1)
	f.AddStep(step2)
	f.AddStep(step3)

	i := f.GetStepIndex("none")
	assert.Equal(t, -1, i)

	i1 := f.GetStepIndex("step1")
	assert.Equal(t, 0, i1)

	i2 := f.GetStepIndex("step2")
	assert.Equal(t, 1, i2)

	i3 := f.GetStepIndex("step3")
	assert.Equal(t, 2, i3)
}

func (f *Flow) dump() {
	fmt.Printf("Flow: %s\n", f.ID)
	for i, step := range f.Steps {
		fmt.Printf("%v-%+v ->", i, step.StepID())
	}

	fmt.Println()
}

func newStep(id string) testStep {
	return testStep{id: id}
}

type testStep struct {
	id string
}

func (t testStep) StepID() string {
	return t.id
}
