package itinerary

import "github.com/openmarketplaceengine/openmarketplaceengine/core/model/job/flow"

type Itinerary struct {
	ID        string
	StartTime string
	flow      *flow.Flow
}

func NewItinerary(id string, startTime string, steps []flow.Step) *Itinerary {
	return &Itinerary{
		ID:        id,
		StartTime: startTime,
		flow:      flow.NewFlow(id, steps),
	}
}
