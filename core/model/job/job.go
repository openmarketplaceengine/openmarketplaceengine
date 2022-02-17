package job

import "time"

//type step int

//const (
//	GoToLocation step = iota
//	PickupPassenger
//	DropOffPassenger
//	CollectCache
//	CollectVoucher
//	CallPhone
//)

// TimeWindow constraint in hours(?)
type TimeWindow struct {
	From int
	To   int
}

// VehicleCapacity constraint in passengers limit or cargo volume.
type VehicleCapacity string

// Location to visit.
type Location struct {
	Longitude  float64
	Latitude   float64
	Name       string
	Address    string
	TimeWindow TimeWindow
}

// Ride represents pickup/drop-off passengers.
type Ride struct {
	PickupLocation  Location
	DropOffLocation Location
	PassengerID     string
}

// Delivery represents pickup/drop-off packages.
type Delivery struct {
	PickupLocation  Location
	DropOffLocation Location
	PackageID       string
}

// Job represents activities assigned to driver.
type Job struct {
	Rides           []Ride
	Deliveries      []Delivery
	DriverID        string
	VehicleCapacity VehicleCapacity
	StartAt         time.Time
}

// Route for Ride or Delivery.
type Route struct {
}

type Routing struct {
	//distance matrix taken from Google Distance Matrix API
	DistanceMatrix [][]int
	Routes         []Route
}
