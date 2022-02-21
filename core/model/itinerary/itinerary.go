package itinerary

import (
	"github.com/openmarketplaceengine/openmarketplaceengine/core/model/job"
	"github.com/openmarketplaceengine/openmarketplaceengine/core/model/step"
	"time"
)

// Itinerary is a planned journey for a job.Job array
type Itinerary struct {
	Jobs        []job.Job
	Steps       []step.Step
	CurrentStep step.Step
	StartTime   time.Time
	WorkerId    string
}

func NewItinerary() *Itinerary {
	return &Itinerary{
		Jobs:        nil,
		Steps:       nil,
		CurrentStep: step.Step{},
		StartTime:   time.Time{},
		WorkerId:    "",
	}
}
