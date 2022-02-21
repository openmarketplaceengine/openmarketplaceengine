package job

import "time"

// Address as defined in geocoding
type Address struct {
	ShortName string
}

// Location to visit
type Location struct {
	Longitude float64
	Latitude  float64
	Name      string
	Address   Address
}

// RideRequest represents ride (and delivery?) pickup/drop-off request
// SubjectID refers to either passenger or package
type RideRequest struct {
	PickupLocation     Location
	DropOffLocation    Location
	SubjectID          string
	RequestedTime      time.Time
	RequestedStartTime time.Time
}

// Job represents activities assigned to Worker
// WorkerID refers to worker.Worker
type Job struct {
	RideRequest RideRequest
	StartTime   time.Time
	EndTime     time.Time
	WorkerID    string
}

// Itinerary is a planned journey for a Job array
type Itinerary struct {
	Jobs        []Job
	Steps       []Step
	CurrentStep Step
	StartTime   time.Time
	WorkerId    string
}

// Action defines primitive action constituting a Job
type Action int

const (
	GoToLocation Action = iota
	Pickup
	DropOff
	CollectCache
	CollectVoucher
	CallPhone
)

// Step is a part of Job execution
// JobID refers to Job step belongs to
type Step struct {
	JobID  string
	Action Action
}
