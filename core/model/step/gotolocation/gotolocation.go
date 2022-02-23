package gotolocation

import (
	"context"
	"fmt"
	"time"

	"github.com/openmarketplaceengine/openmarketplaceengine/core/model/job"

	"github.com/cocoonspace/fsm"
	"github.com/openmarketplaceengine/openmarketplaceengine/core/model/step"
)

const (
	New      step.Status = "New"
	Moving   step.Status = "Moving"
	NearBy   step.Status = "NearBy"
	Arrived  step.Status = "Arrived"
	Canceled step.Status = "Canceled"
)

const (
	Move   step.Action = "Move"
	Near   step.Action = "Near"
	Arrive step.Action = "Arrive"
	Cancel step.Action = "Cancel"
)

type GoToLocation struct {
	StepID    string      `json:",string"`
	UpdatedAt string      `json:",string"`
	Status    step.Status `json:",string"`
	fsm       *fsm.FSM
}

func (gtl *GoToLocation) CurrentStatus() step.Status {
	state := gtl.fsm.Current()
	return stateToStatus[state]
}

func (gtl *GoToLocation) AvailableActions() []step.Action {
	return statusToAvailableActions[gtl.Status]
}

func (gtl *GoToLocation) Handle(action step.Action) error {
	ok := gtl.fsm.Event(actionToEvent[action])
	if !ok {
		return fmt.Errorf("illegal transition from status=%v by action=%v", gtl.Status, action)
	}
	status := stateToStatus[gtl.fsm.Current()]
	err := gtl.updateStatus(context.Background(), status)
	if err != nil {
		return err
	}
	return nil
}

func NewStep(ctx context.Context, stepID step.ID, jobID job.ID) (*step.Step, error) {
	gtl, err := retrieveOrCreate(ctx, stepID)
	if err != nil {
		return nil, err
	}
	return &step.Step{
		ID:         stepID,
		JobID:      jobID,
		Actionable: gtl,
		Atom:       step.GoToLocation,
	}, nil
}

func retrieveOrCreate(ctx context.Context, stepID step.ID) (gtl *GoToLocation, err error) {
	existing, err := storage.Retrieve(ctx, stepID)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		existing.fsm = newFsm(existing.Status)
		return existing, nil
	}

	gtl = &GoToLocation{
		StepID:    string(stepID),
		UpdatedAt: time.Now().Format(time.RFC3339Nano),
		Status:    New,
		fsm:       newFsm(New),
	}

	err = storage.Store(ctx, *gtl)
	if err != nil {
		return
	}
	return
}

func (gtl *GoToLocation) updateStatus(ctx context.Context, status step.Status) error {
	gtl.UpdatedAt = time.Now().Format(time.RFC3339Nano)
	gtl.Status = status
	err := storage.Store(ctx, *gtl)
	if err != nil {
		return err
	}
	return nil
}
