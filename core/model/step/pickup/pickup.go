package pickup

import (
	"context"
	"fmt"
	"time"

	"github.com/cocoonspace/fsm"
	"github.com/openmarketplaceengine/openmarketplaceengine/core/model/step"
)

const (
	Ready     step.State = "Ready"
	Completed step.State = "Completed"
	Canceled  step.State = "Canceled"
)

const (
	CompleteAction step.Action = "CompleteAction"
	CancelAction   step.Action = "CancelAction"
)

type Pickup struct {
	ID        string     `json:",string"`
	JobID     string     `json:",string"`
	UpdatedAt string     `json:",string"`
	State     step.State `json:",string"`
	fsm       *fsm.FSM
}

func (p *Pickup) StepID() string {
	return p.ID
}

func (p *Pickup) CurrentState() step.State {
	state := p.fsm.Current()
	return stateToStatus[state]
}

func (p *Pickup) AvailableActions() []step.Action {
	return statusToAvailableActions[p.State]
}

func (p *Pickup) Handle(action step.Action) error {
	ok := p.fsm.Event(actionToEvent[action])
	if !ok {
		return fmt.Errorf("illegal transition from status=%v by action=%v", p.State, action)
	}
	status := stateToStatus[p.fsm.Current()]
	err := p.updateState(context.Background(), status)
	if err != nil {
		return err
	}
	return nil
}

func New(ctx context.Context, stepID string, jobID string) (pickup *Pickup, err error) {
	// retrieve existing from database or create and store if not exists.
	_ = ctx
	pickup = &Pickup{
		ID:        stepID,
		JobID:     jobID,
		UpdatedAt: time.Now().Format(time.RFC3339Nano),
		State:     Ready,
		fsm:       newFsm(Ready),
	}

	return
}

func (p *Pickup) updateState(ctx context.Context, state step.State) error {
	p.UpdatedAt = time.Now().Format(time.RFC3339Nano)
	p.State = state

	// persist in database
	_ = ctx
	return nil
}
