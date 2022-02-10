package itinerary

import "testing"

func TestFlow(t *testing.T) {
	f := newFlow()

	f.addStep("step1", goToLocation)
	f.addStep("step2", pickUpPassengers)
	f.addStep("step3", dropOffPassengers)

	f.removeStep("step2")
}
