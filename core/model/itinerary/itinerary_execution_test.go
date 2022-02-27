package itinerary

import (
	"testing"

	"github.com/google/uuid"
	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
	"github.com/openmarketplaceengine/openmarketplaceengine/core/model/step"
	"github.com/openmarketplaceengine/openmarketplaceengine/core/model/step/gotolocation"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
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
	ctx := context.Background()
	jobID := uuid.New().String()
	step1, _ := gotolocation.New(ctx, uuid.New().String(), jobID)
	step2, _ := gotolocation.New(ctx, uuid.New().String(), jobID)
	step3, _ := gotolocation.New(ctx, uuid.New().String(), jobID)

	itinerary := New(uuid.New().String(), []step.Step{step1, step2, step3})

	require.Len(t, itinerary.Steps, 3)
	currentStep, _ := itinerary.CurrentStep()
	require.Equal(t, step1, currentStep)

	actions0 := currentStep.AvailableActions()
	require.Len(t, actions0, 2)
	require.ElementsMatch(t, actions0, []step.Action{gotolocation.NearAction, gotolocation.CancelAction})
	require.Equal(t, gotolocation.Moving, step1.CurrentState())

	err := itinerary.Handle(actions0[0])
	require.NoError(t, err)

	nextStep, _ := itinerary.CurrentStep()
	require.Equal(t, gotolocation.Near, nextStep.CurrentState())
	actions1 := nextStep.AvailableActions()
	require.Len(t, actions1, 2)
	require.ElementsMatch(t, actions1, []step.Action{gotolocation.ArriveAction, gotolocation.CancelAction})

	err = nextStep.Handle(actions1[0])
	require.NoError(t, err)

	nextStep, _ = itinerary.CurrentStep()
	require.Equal(t, gotolocation.Arrived, nextStep.CurrentState())
	actions2, _ := itinerary.AvailableActions()
	require.Len(t, actions2, 0)
}

func testItineraryExecutionAllSteps(t *testing.T) {
	ctx := context.Background()
	jobID := uuid.New().String()
	step1, _ := gotolocation.New(ctx, uuid.New().String(), jobID)
	step2, _ := gotolocation.New(ctx, uuid.New().String(), jobID)
	step3, _ := gotolocation.New(ctx, uuid.New().String(), jobID)

	itinerary := New(uuid.New().String(), []step.Step{step1, step2, step3})

	require.Len(t, itinerary.Steps, 3)
	currentStep, _ := itinerary.CurrentStep()
	counter := 0
	for len(currentStep.AvailableActions()) != 0 {
		actions, _ := itinerary.AvailableActions()

		for len(actions) != 0 {
			err := itinerary.Handle(actions[0])
			require.NoError(t, err)
			counter++
			currentStep, _ = itinerary.CurrentStep()
			break
		}
	}

	require.Equal(t, 6, counter)
}
