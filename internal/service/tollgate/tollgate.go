package tollgate

import (
	"fmt"

	"context"
)

// Detector represents a tollgate the Subject passes through.
// Crossed detects if subject Movement has travelled through the tollgate.
type Detector interface {
	DetectCrossing(ctx context.Context, movement *Movement) (*Crossing, error)
}

// Line represents a two points line that Subject crosses.
// longitude latitude corresponding to linear algebra X Y axis.
type Line struct {
	Lon1 float64
	Lat1 float64
	Lon2 float64
	Lat2 float64
}

// BBox is a bounded box, represents area defined by two longitudes and two latitudes
// left,bottom,right,top - lonMin, latMin, lonMax, latMax
// Longitude and Latitude correspond to linear algebra X and Y axis.
type BBox struct {
	LonMin float64
	LatMin float64
	LonMax float64
	LatMax float64
}

// Alg represents algorithm.
type Alg uint8

const (
	LineAlg Alg = iota
	BBoxAlg
	VectorAlg
)

// Crossing represents detected fact of passing through the tollgate by DriverID.
type Crossing struct {
	TollgateID string
	DriverID   string
	Movement   *Movement
	Direction  Direction
	Alg        Alg
}

// Movement represents a moving SubjectID from one Location to another.
type Movement struct {
	SubjectID string
	From      *Location
	To        *Location
}

// Location is longitude, latitude corresponding to linear algebra X, Y axis.
type Location struct {
	Lon float64
	Lat float64
}

//Direction to North, South, East or West in form of N, S, E, W, NE, NW, SE, SW.
type Direction string

// Direction of movement
// When moving to North - latitude increases until 90.
// When moving to South - latitude decreases until -90.
// When moving to East - longitude increases until 180.
// When moving to West - longitude decreases until -180.
// If 90/180 limit crossed it jumps to -90/-180 and vice versa.
// Movement represents a moving subject
// returns Direction in form of N, S, E, W, NE, NW, SE, SW.
func (m *Movement) Direction() Direction {
	fromX := m.From.Lon
	fromY := m.From.Lat
	toX := m.To.Lon
	toY := m.To.Lat

	var pole string
	var side string
	if fromY < toY {
		pole = "N"
	}
	if fromY > toY {
		pole = "S"
	}

	if fromX < toX {
		side = "E"
	}
	if fromX > toX {
		side = "W"
	}

	return Direction(fmt.Sprintf("%s%s", pole, side))
}
