package job

import (
	"fmt"
	"time"
)

type ID string

type Status int

const (
	New Status = iota
	InProgress
	Completed
	Canceled
)

func (s Status) String() string {
	switch s {
	case New:
		return "New"
	case InProgress:
		return "InProgress"
	case Completed:
		return "Completed"
	case Canceled:
		return "Canceled"
	default:
		return fmt.Sprintf("%d", s)
	}
}

// Address as defined in geocoding.
type Address struct {
	ShortName string
}

// Location to visit.
type Location struct {
	Longitude float64
	Latitude  float64
	Name      string
	Address   Address
}

// Transportation represents the action of ride or delivery
// SubjectID refers to either passenger or package.
type Transportation struct {
	PickupLocation     Location
	DropOffLocation    Location
	SubjectID          string
	RequestedTime      time.Time
	RequestedStartTime time.Time
}

// Job represents activities assigned to Worker
// WorkerID refers to worker.Worker.
type Job struct {
	ID             ID
	Transportation Transportation
	Status         Status
	StartTime      time.Time
	EndTime        time.Time
}

func NewJob(id ID, transportation Transportation, startTime time.Time) *Job {
	return &Job{
		ID:             id,
		Transportation: transportation,
		Status:         New,
		StartTime:      startTime,
		EndTime:        time.Time{},
	}
}
