package flow

import (
	"fmt"
)

type Step interface {
	StepID() string
}

type Flow struct {
	ID    string
	Steps []Step
}

func NewFlow(id string, steps []Step) *Flow {
	return &Flow{
		ID:    id,
		Steps: steps,
	}
}

func (f *Flow) GetFirstStep() (Step, error) {
	if len(f.Steps) == 0 {
		return nil, fmt.Errorf("flow %s has no steps", f.ID)
	}
	step := f.Steps[0]
	return step, nil
}

func (f *Flow) AddStep(step Step) {
	f.Steps = append(f.Steps, step)
}

func (f *Flow) RemoveStep(stepID string) {
	i := f.GetStepIndex(stepID)
	if i > -1 {
		f.Steps = append(f.Steps[:i], f.Steps[i+1:]...)
	}
}

func (f *Flow) GetStep(stepID string) Step {
	for _, step := range f.Steps {
		if step.StepID() == stepID {
			return step
		}
	}
	return nil
}

func (f *Flow) GetStepIndex(stepID string) int {
	for i, step := range f.Steps {
		if step.StepID() == stepID {
			return i
		}
	}
	return -1
}
