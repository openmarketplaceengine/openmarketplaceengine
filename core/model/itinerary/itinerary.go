package itinerary

import (
	"fmt"
	"time"

	"github.com/openmarketplaceengine/openmarketplaceengine/core/model/job"
	"github.com/openmarketplaceengine/openmarketplaceengine/core/model/step"
)

type Status int

const (
	Scheduled Status = iota
	InProgress
	Completed
	Canceled
)

// Itinerary is a planned journey for a job.Job array.
type Itinerary struct {
	ID          string
	Jobs        []*job.Job
	Steps       []*step.Step
	CurrentStep int
	StartTime   time.Time
	WorkerID    string
	Status      Status
}

func NewItinerary(id string, steps []*step.Step) *Itinerary {
	return &Itinerary{
		ID:          id,
		Jobs:        nil,
		Steps:       steps,
		CurrentStep: 0,
		StartTime:   time.Time{},
		WorkerID:    "",
	}
}

func (it *Itinerary) GetFirstStep() (*step.Step, error) {
	if len(it.Steps) == 0 {
		return nil, fmt.Errorf("itinerary %s has no steps", it.ID)
	}
	s := it.Steps[0]
	return s, nil
}

func (it *Itinerary) AddStep(step *step.Step) {
	it.Steps = append(it.Steps, step)
}

func (it *Itinerary) RemoveStep(stepID string) {
	i := it.GetStepIndex(stepID)
	if i > -1 {
		copy(it.Steps[i:], it.Steps[i+1:])
		n := len(it.Steps)
		it.Steps[n-1] = nil
		it.Steps = it.Steps[:n-1]
	}
}

func (it *Itinerary) GetStep(stepID string) *step.Step {
	for _, s := range it.Steps {
		if s.ID == stepID {
			return s
		}
	}
	return nil
}

func (it *Itinerary) GetStepIndex(stepID string) int {
	for i, s := range it.Steps {
		if s.ID == stepID {
			return i
		}
	}
	return -1
}
