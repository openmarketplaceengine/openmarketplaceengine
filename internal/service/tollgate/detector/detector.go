package detector

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

// Alg represents algorithm.
type Alg uint8

const (
	LineAlg Alg = iota
	BBoxAlg
	VectorAlg
)

// Crossing represents detected fact of passing through the tollgate by WorkerID.
type Crossing struct {
	TollgateID string
	WorkerID   string
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

type Tollgate struct {
	ID             string
	Line           *Line
	BBoxes         []*BBox
	BBoxesRequired int
}

type Detector struct {
	tollgates []*Tollgate
	storage   *storage
}

func NewDetector(tollgates []*Tollgate, client *redis.Client) *Detector {
	return &Detector{
		tollgates: tollgates,
		storage:   newStorage(client),
	}
}

// DetectCrossing detects if subject Movement has travelled through the tollgate.
func (d *Detector) DetectCrossing(ctx context.Context, movement *Movement) (*Crossing, error) {
	for _, t := range d.tollgates {
		if t.Line != nil {
			crossing := DetectCrossingLine(t.ID, t.Line, movement)
			if crossing != nil {
				return crossing, nil
			}
		}

		if len(t.BBoxes) > 0 {
			crossing, err := DetectCrossingBBox(ctx, d.storage, t.ID, t.BBoxes, t.BBoxesRequired, movement)
			if err != nil {
				return nil, err
			}
			if crossing != nil {
				return crossing, nil
			}
		}
	}
	return nil, nil
}
