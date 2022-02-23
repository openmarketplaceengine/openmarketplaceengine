package step

import (
	"fmt"

	"github.com/openmarketplaceengine/openmarketplaceengine/core/model/job"
)

// ID is a unique identifier of Step.
type ID string

// Action defines command that changes Step state.
type Action string

// Status is a description of Step status.
// Step in certain Status is eligible to execute certain transitions.
type Status string

// Atom represents an individual itinerary.Itinerary Step type.
type Atom int

const (
	GoToLocation Atom = iota
	Pickup
	DropOff
	CollectCache
	CollectVoucher
	CallPhone
)

func (a Atom) String() string {
	switch a {
	case GoToLocation:
		return "GoToLocation"
	case Pickup:
		return "Pickup"
	case DropOff:
		return "DropOff"
	default:
		return fmt.Sprintf("%d", a)
	}
}

// Actionable interface exposes current Status and Action list for eligible transitions.
// Handle(Action) function performs state transition.
// AvailableActions empty list means entity has reached its final Status.
type Actionable interface {
	CurrentStatus() Status
	AvailableActions() []Action
	Handle(action Action) error
}

// Step is a part of Job execution
// JobID refers to job.Job step belongs to.
type Step struct {
	ID         ID
	JobID      job.ID
	Atom       Atom
	Actionable Actionable
}

func (s *Step) CurrentStatus() Status {
	return s.Actionable.CurrentStatus()
}

func (s *Step) AvailableActions() []Action {
	return s.Actionable.AvailableActions()
}

func (s *Step) Handle(action Action) error {
	return s.Actionable.Handle(action)
}
