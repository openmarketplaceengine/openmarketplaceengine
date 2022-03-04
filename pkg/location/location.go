package location

import (
	"time"
)

// QueryLocation represents geographic coordinates of principal relative to querying coordinates.
type QueryLocation struct {
	WorkerID      string
	Longitude     float64
	Latitude      float64
	Distance      float64
	FromLongitude float64
	FromLatitude  float64
	LastSeen      time.Time
}

// Location represents geographic coordinates of WorkerID.
type Location struct {
	WorkerID  string
	Longitude float64
	Latitude  float64
}
