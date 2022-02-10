package itinerary

type Stepf func()

type Passenger struct {
	ID   string
	Name string
}

type Itinerary struct {
	ID        string
	StartTime string
	flow      *flow
}

func (i *Itinerary) AddGoToLocationStep() {
	i.flow.addStep("abc", goToLocation)
}

func NewItinerary(id string, startTime string) *Itinerary {
	return &Itinerary{
		ID:        id,
		StartTime: startTime,
		flow:      newFlow(),
	}
}
