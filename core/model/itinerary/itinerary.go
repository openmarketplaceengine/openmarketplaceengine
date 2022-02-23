package itinerary

import (
	"context"
	"fmt"
	"time"

	"github.com/openmarketplaceengine/openmarketplaceengine/core/model/step/gotolocation"

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
func NewItinerary(id ID, jobs []*job.Job) (*Itinerary, error) {
	it := &Itinerary{
		ID:          id,
		Jobs:        jobs,
		Steps:       []*step.Step{},
		CurrentStep: 0,
		StartTime:   time.Time{},
		WorkerID:    "",
	}

	for _, j := range jobs {
		goTo, err := newStep(context.Background(), step.ID(fmt.Sprintf("%s-%s", j.ID, "step-1")), j.ID, step.GoToLocation)
		if err != nil {
			return nil, err
		}
		it.AddStep(goTo)

		pickup, err := newStep(context.Background(), step.ID(fmt.Sprintf("%s-%s", j.ID, "step-2")), j.ID, step.Pickup)
		if err != nil {
			return nil, err
		}
		it.AddStep(pickup)

		dropOff, err := newStep(context.Background(), step.ID(fmt.Sprintf("%s-%s", j.ID, "step-2")), j.ID, step.DropOff)
		if err != nil {
			return nil, err
		}
		it.AddStep(dropOff)
	}

	return it, nil
}

func (it *Itinerary) GetCurrentStep() *step.Step {
	if len(it.Steps) == 0 {
		return nil
	}
	return it.Steps[it.CurrentStep]
}

func (it *Itinerary) CurrentStatus() step.Status {
	currentStep := it.GetCurrentStep()
	return currentStep.CurrentStatus()
}

func (it *Itinerary) AvailableActions() []step.Action {
	currentStep := it.GetCurrentStep()
	return currentStep.AvailableActions()
}

func (it *Itinerary) Handle(action step.Action) error {
	currentStep := it.GetCurrentStep()
	err := currentStep.Handle(action)
	if err != nil {
		return err
	}

	if len(currentStep.AvailableActions()) == 0 {
		if it.CurrentStep < len(it.Steps)-1 {
			it.CurrentStep++
		}
	}

	return nil
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

func newStep(ctx context.Context, stepID step.ID, jobID job.ID, atom step.Atom) (*step.Step, error) {
	switch atom {
	case step.GoToLocation:
		return gotolocation.NewStep(ctx, stepID, jobID)
	case step.Pickup:
		return gotolocation.NewStep(ctx, stepID, jobID)
	case step.DropOff:
		return gotolocation.NewStep(ctx, stepID, jobID)
	default:
		return nil, fmt.Errorf("unsupported step atom %q", atom)
	}
}
