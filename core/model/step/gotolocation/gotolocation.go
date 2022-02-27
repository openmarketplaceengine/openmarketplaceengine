package gotolocation

import (
	"context"
	"fmt"
	"time"

	"github.com/cocoonspace/fsm"
	"github.com/openmarketplaceengine/openmarketplaceengine/core/model/step"
)

const (
	Moving   step.State = "Moving"
	Near     step.State = "Near"
	Arrived  step.State = "Arrived"
	Canceled step.State = "Canceled"
)

const (
	NearAction   step.Action = "NearAction"
	ArriveAction step.Action = "ArriveAction"
	CancelAction step.Action = "CancelAction"
)

type GoToLocation struct {
	ID        string     `json:",string"`
	JobID     string     `json:",string"`
	UpdatedAt string     `json:",string"`
	State     step.State `json:",string"`
	fsm       *fsm.FSM
}

func (gtl *GoToLocation) StepID() string {
	return gtl.ID
}

func (gtl *GoToLocation) CurrentState() step.State {
	state := gtl.fsm.Current()
	return stateToStatus[state]
}

func (gtl *GoToLocation) AvailableActions() []step.Action {
	return statusToAvailableActions[gtl.State]
}

func (gtl *GoToLocation) Handle(action step.Action) error {
	ok := gtl.fsm.Event(actionToEvent[action])
	if !ok {
		return fmt.Errorf("illegal transition from status=%v by action=%v", gtl.State, action)
	}
	status := stateToStatus[gtl.fsm.Current()]
	err := gtl.updateStatus(context.Background(), status)
	if err != nil {
		return err
	}
	return nil
}

func New(ctx context.Context, stepID string, jobID string) (*GoToLocation, error) {
	// retrieve existing from database or create and store if not exists.
	_ = ctx

	gtl := &GoToLocation{
		ID:        stepID,
		JobID:     jobID,
		UpdatedAt: time.Now().Format(time.RFC3339Nano),
		State:     Moving,
		fsm:       newFsm(Moving),
	}

	return gtl, nil
}

func (gtl *GoToLocation) updateStatus(ctx context.Context, status step.State) error {
	gtl.UpdatedAt = time.Now().Format(time.RFC3339Nano)
	gtl.State = status
	// persist in database
	_ = ctx
	return nil
}
