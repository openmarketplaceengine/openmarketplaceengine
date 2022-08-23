package job

import (
	"encoding/json"
	"time"
)

type Estimate struct {
	DistanceMeters int
	Duration       time.Duration
}

type Job struct {
	ID      string
	Pickup  LatLon
	DropOff LatLon
}

type EstimatedJob struct {
	ID              string
	ToPickup        Estimate
	PickupToDropOff Estimate
	WorkerLocation  Location
	Pickup          Location
	DropOff         Location
}

func (j EstimatedJob) MarshalBinary() ([]byte, error) {
	return json.Marshal(j)
}

func (j *EstimatedJob) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &j)
}

type Location struct {
	Address string
	Lat     float64
	Lon     float64
}

type RankBy int

const (
	TimeAwaited       RankBy = iota
	DistanceToPickup  RankBy = iota
	QuickestToDeliver RankBy = iota
)
