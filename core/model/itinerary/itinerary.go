package itinerary

import (
	"fmt"
	"time"

	"github.com/cocoonspace/fsm"

	"github.com/driverscooperative/geosrv/core/model/step"
)

// Itinerary is a planned journey represented as step.Step array.
type Itinerary struct {
	ID        string
	Steps     []*step.Step
	step      int // current step
	StartTime time.Time
	WorkerID  string
}

// New builds planned slice of step.Step for a slice of job.Job.
// Each job.Request consists of predefined set of steps,
// i.e. step.GoToLocation, step.Pickup, step.DropOff, etc.
func New(id string, steps []*step.Step) *Itinerary {
	it := &Itinerary{
		ID:        id,
		Steps:     steps,
		step:      0,
		StartTime: time.Time{},
		WorkerID:  "",
	}
	return it
}

func (it *Itinerary) CurrentStep() (*step.Step, error) {
	if len(it.Steps) == 0 {
		return nil, fmt.Errorf("no current step")
	}
	return it.Steps[it.step], nil
}

func (it *Itinerary) AvailableEvents() ([]fsm.Event, error) {
	currentStep, err := it.CurrentStep()
	if err != nil {
		return nil, err
	}
	return currentStep.AvailableEvents(), nil
}

func (it *Itinerary) Handle(event fsm.Event) error {
	currentStep, err := it.CurrentStep()
	if err != nil {
		return err
	}
	err = currentStep.Handle(event)
	if err != nil {
		return err
	}

	if len(currentStep.AvailableEvents()) == 0 {
		if it.step < len(it.Steps)-1 {
			it.step++
		}
	}

	return nil
}

func (it *Itinerary) AddStep(step *step.Step) {
	it.Steps = append(it.Steps, step)
}

func (it *Itinerary) RemoveStep(id string) {
	i := it.stepIndex(id)
	if i > -1 {
		copy(it.Steps[i:], it.Steps[i+1:])
		n := len(it.Steps)
		it.Steps[n-1] = nil
		it.Steps = it.Steps[:n-1]
	}
}

func (it *Itinerary) GetStep(id string) *step.Step {
	for _, s := range it.Steps {
		if s.StepID() == id {
			return s
		}
	}
	return nil
}

func (it *Itinerary) stepIndex(id string) int {
	for i, s := range it.Steps {
		if s.StepID() == id {
			return i
		}
	}
	return -1
}
