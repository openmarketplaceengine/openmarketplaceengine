package location

import (
	"time"
)

// QueryLocation represents geographic coordinates of principal relative to querying coordinates.
type QueryLocation struct {
	PrincipalID string
	Longitude   float64
	Latitude    float64
	Distance    float64
	LastSeen    time.Time
}

// UpdateLocation represents geographic coordinates of principal.
type UpdateLocation struct {
	PrincipalID string
	Longitude   float64
	Latitude    float64
}
