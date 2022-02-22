package itinerary

import (
	"fmt"
	"time"

	"github.com/openmarketplaceengine/openmarketplaceengine/core/model/job"
	"github.com/openmarketplaceengine/openmarketplaceengine/core/model/step"
)

type ID string

type Status int

const (
	Scheduled Status = iota
	InProgress
	Completed
	Canceled
)

// Itinerary is a planned journey for a job.Job array.
type Itinerary struct {
	ID          ID
	Jobs        []*job.Job
	Steps       []*step.Step
	CurrentStep int
	StartTime   time.Time
	WorkerID    string
	Status      Status
}

// NewItinerary builds planned slice of step.Step for a slice of job.Job.
// Each job.Transportation consists of predefined set of steps,
// i.e. step.GoToLocation, step.Pickup, step.DropOff, etc.
func NewItinerary(id ID, jobs []*job.Job) *Itinerary {

	i := &Itinerary{
		ID:          id,
		Jobs:        jobs,
		Steps:       []*step.Step{},
		CurrentStep: 0,
		StartTime:   time.Time{},
		WorkerID:    "",
	}

	for _, j := range jobs {
		step1 := &step.Step{
			ID:     step.ID(fmt.Sprintf("%s-%s", j.ID, "step-1")),
			JobID:  j.ID,
			Action: step.GoToLocation,
		}
		i.AddStep(step1)

		step2 := &step.Step{
			ID:     step.ID(fmt.Sprintf("%s-%s", j.ID, "step-2")),
			JobID:  j.ID,
			Action: step.Pickup,
		}
		i.AddStep(step2)

		step3 := &step.Step{
			ID:     step.ID(fmt.Sprintf("%s-%s", j.ID, "step-3")),
			JobID:  j.ID,
			Action: step.DropOff,
		}
		i.AddStep(step3)
	}

	return i
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

func (it *Itinerary) RemoveStep(id step.ID) {
	i := it.GetStepIndex(id)
	if i > -1 {
		copy(it.Steps[i:], it.Steps[i+1:])
		n := len(it.Steps)
		it.Steps[n-1] = nil
		it.Steps = it.Steps[:n-1]
	}
}

func (it *Itinerary) GetStep(id step.ID) *step.Step {
	for _, s := range it.Steps {
		if s.ID == id {
			return s
		}
	}
	return nil
}

func (it *Itinerary) GetStepIndex(id step.ID) int {
	for i, s := range it.Steps {
		if s.ID == id {
			return i
		}
	}
	return -1
}
