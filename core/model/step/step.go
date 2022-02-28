package step

// Action defines command that changes Step state.
type Action string

// State is a description of Step status.
// Step in certain State is eligible to execute certain transitions.
type State string

// Step interface exposes current State and Action list for eligible transitions.
// Handle(Action) function performs state transition.
// AvailableActions (i.e. Move, Arrive, PickUp, DropOff), empty list means entity has reached its final State.
type Step interface {
	StepID() string
	CurrentState() State
	AvailableActions() []Action
	Handle(action Action) error
}
