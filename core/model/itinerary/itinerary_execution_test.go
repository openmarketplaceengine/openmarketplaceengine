package itinerary

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/openmarketplaceengine/openmarketplaceengine/core/model/step/gotolocation"

	"github.com/openmarketplaceengine/openmarketplaceengine/core/model/job"

	"github.com/openmarketplaceengine/openmarketplaceengine/core/model/step"
	"github.com/stretchr/testify/require"
)

func TestItineraryExecution(t *testing.T) {
	t.Run("testItineraryExecutionOneStep", func(t *testing.T) {
		testItineraryExecutionOneStep(t)
	})

	t.Run("testItineraryExecutionAllSteps", func(t *testing.T) {
		testItineraryExecutionAllSteps(t)
	})
}

func testItineraryExecutionOneStep(t *testing.T) {
	itinerary, err := NewItinerary("test-flow-1", []*job.Job{createJob()})
	require.NoError(t, err)

	require.Len(t, itinerary.Steps, 3)
	currentStep := itinerary.GetCurrentStep()
	require.Equal(t, step.GoToLocation, currentStep.Atom)

	actionable := currentStep.Actionable
	actions0 := actionable.AvailableActions()
	require.Len(t, actions0, 2)
	require.ElementsMatch(t, actions0, []step.Action{gotolocation.Move, gotolocation.Cancel})
	require.Equal(t, gotolocation.New, actionable.CurrentStatus())

	err = actionable.Handle(actions0[0])
	require.NoError(t, err)

	require.Equal(t, gotolocation.Moving, actionable.CurrentStatus())
	actions1 := actionable.AvailableActions()
	require.Len(t, actions1, 2)
	require.ElementsMatch(t, actions1, []step.Action{gotolocation.Near, gotolocation.Cancel})

	err = actionable.Handle(actions1[0])
	require.NoError(t, err)

	require.Equal(t, gotolocation.NearBy, actionable.CurrentStatus())
	actions2 := actionable.AvailableActions()
	require.Len(t, actions2, 2)
	require.ElementsMatch(t, actions2, []step.Action{gotolocation.Arrive, gotolocation.Cancel})

	err = actionable.Handle(actions2[0])
	require.NoError(t, err)

	require.Equal(t, gotolocation.Arrived, actionable.CurrentStatus())
	actions3 := actionable.AvailableActions()
	require.Len(t, actions3, 0)
}

func testItineraryExecutionAllSteps(t *testing.T) {
	itinerary, err := NewItinerary("test-flow-1", []*job.Job{createJob()})
	require.NoError(t, err)

	require.Len(t, itinerary.Steps, 3)
	currentStep := itinerary.GetCurrentStep()
	counter := 0
	for len(currentStep.AvailableActions()) != 0 {
		actionable := currentStep.Actionable
		actions := actionable.AvailableActions()

		for len(actions) != 0 {
			err = itinerary.Handle(actions[0])
			require.NoError(t, err)
			counter++
			currentStep = itinerary.GetCurrentStep()
			break
		}
	}

	require.Equal(t, 9, counter)
}

func createJob() *job.Job {
	return &job.Job{
		ID: job.ID(uuid.New().String()),
		Transportation: job.Transportation{
			PickupLocation: job.Location{
				Longitude: 1,
				Latitude:  1,
				Name:      "pickup-at",
				Address:   job.Address{},
			},
			DropOffLocation: job.Location{
				Longitude: 1,
				Latitude:  1,
				Name:      "drop-at",
				Address:   job.Address{},
			},
			SubjectID:          "passenger-1",
			RequestedTime:      time.Now(),
			RequestedStartTime: time.Time{},
		},
		StartTime: time.Now(),
		EndTime:   time.Time{},
	}
}
