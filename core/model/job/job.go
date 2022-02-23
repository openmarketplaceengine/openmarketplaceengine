package job

import (
	"time"
)

type ID string

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

// Job represents activities assigned to Worker.
type Job struct {
	ID             ID
	Transportation Transportation
	StartTime      time.Time
	EndTime        time.Time
}

func NewJob(id ID, transportation Transportation, startTime time.Time) *Job {
	return &Job{
		ID:             id,
		Transportation: transportation,
		StartTime:      startTime,
		EndTime:        time.Time{},
	}
}
