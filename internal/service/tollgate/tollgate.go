package tollgate

import (
	"fmt"

	"context"
)

// Tollgate represents a tollgate the Subject passes through.
// Crossed detects if subject Movement has travelled through the tollgate.
type Tollgate interface {
	DetectCrossing(ctx context.Context, movement *Movement) (*Crossing, error)
}

// Crossing represents detected fact of passing through the tollgate by SubjectID.
type Crossing struct {
	TollgateID string
	SubjectID  string
	Location   *LocationXY
	Direction  Direction
}

// Movement represents a moving SubjectID from one LocationXY to another.
type Movement struct {
	SubjectID string
	From      *LocationXY
	To        *LocationXY
}

// LocationXY is longitude latitude corresponding to linear algebra X Y axis.
type LocationXY struct {
	LongitudeX float64
	LatitudeY  float64
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
	fromX := m.From.LongitudeX
	fromY := m.From.LatitudeY
	toX := m.To.LongitudeX
	toY := m.To.LatitudeY

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
