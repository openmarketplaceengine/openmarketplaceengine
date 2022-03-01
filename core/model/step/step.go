package step

import (
	"time"

	"golang.org/x/net/context"
)

// Type defines Step types i.e. GoToLocation, Pickup, DropOff.
type Type string

const (
	GoToLocation Type = "GoToLocation"
	Pickup       Type = "Pickup"
	DropOff      Type = "DropOff"
)

// Event defines command that changes State.
type Event string

// State is a description of Step status.
// Step in certain State is eligible to execute certain transitions.
type State string

// StateMachine interface exposes current State and Event list for eligible transitions.
// Handle(Event) function performs state transition.
// AvailableEvents (i.e. Move, Arrive, PickUp, DropOff), empty list means entity has reached its final State.
type StateMachine interface {
	CurrentState() State
	AvailableEvents() []Event
	Handle(event Event) error
}

// Step represents itinerary step.
type Step struct {
	ID           string
	Type         Type
	State        State
	StateMachine StateMachine
	JobID        string
	UpdatedAt    string
}

func New(id string, jobID string, stateAware StateMachine, stepType Type) *Step {
	state := stateAware.CurrentState()
	return &Step{
		ID:           id,
		JobID:        jobID,
		UpdatedAt:    time.Now().String(),
		State:        state,
		Type:         stepType,
		StateMachine: stateAware,
	}
}

func (s *Step) StepID() string {
	return s.ID
}

func (s *Step) CurrentState() State {
	return s.StateMachine.CurrentState()
}

func (s *Step) AvailableEvents() []Event {
	return s.StateMachine.AvailableEvents()
}

func (s *Step) Handle(event Event) error {
	return s.StateMachine.Handle(event)
}

func (s *Step) updateState(ctx context.Context, state State) error {
	s.UpdatedAt = time.Now().Format(time.RFC3339Nano)
	s.State = state
	// persist in database type, state etc.
	_ = ctx
	panic("implement me")
}
