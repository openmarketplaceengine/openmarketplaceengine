package state

import (
	"github.com/cocoonspace/fsm"
	"time"
)

type State struct {
	DriverID                string
	PassengerIDs            []string
	DestinationLatitude     float64
	DestinationLongitude    float64
	CreatedAt               time.Time
	LastModifiedAt          time.Time
	LastModifiedAtLatitude  float64
	LastModifiedAtLongitude float64
	LastEvent               fsm.Event
	LastState               fsm.State
}
