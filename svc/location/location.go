package location

import (
	"time"
)

// RangeLocation represents geographic coordinates of principal relative to querying coordinates.
type RangeLocation struct {
	WorkerID     string
	Longitude    float64
	Latitude     float64
	Distance     float64
	FromLon      float64
	FromLat      float64
	LastSeenTime time.Time
}

// WorkerLocation represents geographic coordinates of worker with LastSeenTime.
type WorkerLocation struct {
	WorkerID     string
	Longitude    float64
	Latitude     float64
	LastSeenTime time.Time
}

// Location represents geographic coordinates of WorkerID.
type Location struct {
	WorkerID  string
	Longitude float64
	Latitude  float64
}
