package job

import "time"

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

// RideRequest represents ride (and delivery?) pickup/drop-off request
// SubjectID refers to either passenger or package.
type RideRequest struct {
	PickupLocation     Location
	DropOffLocation    Location
	SubjectID          string
	RequestedTime      time.Time
	RequestedStartTime time.Time
}

// Job represents activities assigned to Worker
// WorkerID refers to worker.Worker.
type Job struct {
	RideRequest RideRequest
	StartTime   time.Time
	EndTime     time.Time
}

func NewJob(rideRequest RideRequest, startTime time.Time) *Job {
	return &Job{
		RideRequest: rideRequest,
		StartTime:   startTime,
		EndTime:     time.Time{},
	}
}
