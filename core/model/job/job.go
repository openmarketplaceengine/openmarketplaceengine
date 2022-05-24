package job

import (
	"time"
)

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

// Request represents the action of ride or delivery
// SubjectID refers to either passenger or package.
type Request struct {
	PickupLocation     Location
	DropOffLocation    Location
	SubjectID          string
	RequestedByID      string
	RequestedTime      time.Time
	RequestedStartTime time.Time
}

// Job represents activities assigned to Worker.
type Job struct {
	ID        string
	Request   Request
	StartTime time.Time
	EndTime   time.Time
}

func NewJob(id string, request Request, startTime time.Time) *Job {
	return &Job{
		ID:        id,
		Request:   request,
		StartTime: startTime,
		EndTime:   time.Time{},
	}
}
