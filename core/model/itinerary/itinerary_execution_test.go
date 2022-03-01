package itinerary

import (
	"testing"

	"github.com/cocoonspace/fsm"

	"github.com/google/uuid"
	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
	"github.com/openmarketplaceengine/openmarketplaceengine/core/model/step"
	"github.com/openmarketplaceengine/openmarketplaceengine/core/model/step/gotolocation"
	"github.com/stretchr/testify/require"
)

func TestItineraryExecution(t *testing.T) {
	err := cfg.Load()
	require.NoError(t, err)

	t.Run("testItineraryExecutionOneStep", func(t *testing.T) {
		testItineraryExecutionOneStep(t)
	})

	t.Run("testItineraryExecutionAllSteps", func(t *testing.T) {
		testItineraryExecutionAllSteps(t)
	})
}

func testItineraryExecutionOneStep(t *testing.T) {
	jobID := uuid.New().String()
	goToLocation1 := gotolocation.New(gotolocation.Moving)
	goToLocation2 := gotolocation.New(gotolocation.Moving)

	itinerary := New(uuid.New().String(), []*step.Step{})

	step1 := step.New(uuid.NewString(), jobID, goToLocation1, step.GoToLocation)
	itinerary.AddStep(step1)
	step2 := step.New(uuid.NewString(), jobID, goToLocation2, step.GoToLocation)
	itinerary.AddStep(step2)

	require.Len(t, itinerary.Steps, 2)
	currentStep, err := itinerary.CurrentStep()
	require.NoError(t, err)
	require.Equal(t, step1, currentStep)

	events0 := currentStep.AvailableEvents()
	require.Len(t, events0, 2)
	require.ElementsMatch(t, events0, []fsm.Event{gotolocation.NearBy, gotolocation.Cancel})
	require.Equal(t, gotolocation.Moving, step1.CurrentState())

	err = itinerary.Handle(events0[0])
	require.NoError(t, err)

	nextStep, _ := itinerary.CurrentStep()
	require.Equal(t, gotolocation.Near, nextStep.CurrentState())
	events1 := nextStep.AvailableEvents()
	require.Len(t, events1, 2)
	require.ElementsMatch(t, events1, []fsm.Event{gotolocation.Arrive, gotolocation.Cancel})

	err = nextStep.Handle(events1[0])
	require.NoError(t, err)

	nextStep, _ = itinerary.CurrentStep()
	require.Equal(t, gotolocation.Arrived, nextStep.CurrentState())
	events2, _ := itinerary.AvailableEvents()
	require.Len(t, events2, 0)
}

func testItineraryExecutionAllSteps(t *testing.T) {
	jobID := uuid.New().String()

	goToLocation1 := gotolocation.New(gotolocation.Moving)
	goToLocation2 := gotolocation.New(gotolocation.Moving)

	step1 := step.New(uuid.NewString(), jobID, goToLocation1, step.GoToLocation)
	step2 := step.New(uuid.NewString(), jobID, goToLocation2, step.GoToLocation)
	itinerary := New(uuid.New().String(), []*step.Step{step1, step2})

	require.Len(t, itinerary.Steps, 2)
	currentStep, _ := itinerary.CurrentStep()
	counter := 0
	for len(currentStep.AvailableEvents()) != 0 {
		events, _ := itinerary.AvailableEvents()

		for len(events) != 0 {
			err := itinerary.Handle(events[0])
			require.NoError(t, err)
			counter++
			currentStep, _ = itinerary.CurrentStep()
			break
		}
	}

	require.Equal(t, 4, counter)
}
