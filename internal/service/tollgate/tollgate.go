package tollgate

import (
	"fmt"
	"math"
)

type Detector struct {
	Tollgates []*Tollgate
}

func New(tollgates []*Tollgate) *Detector {
	return &Detector{
		Tollgates: tollgates,
	}
}

func NewNoOp() *Detector {
	var tollgates []*Tollgate
	return &Detector{
		Tollgates: tollgates,
	}
}

// LocationXY is longitude latitude corresponding to linear algebra X Y axis.
type LocationXY struct {
	LongitudeX float64
	LatitudeY  float64
}

// Tollgate represents a two points line that Subject crosses.
type Tollgate struct {
	ID     string
	Point1 LocationXY
	Point2 LocationXY
}

// Movement represents a moving SubjectID from one LocationXY to another.
type Movement struct {
	SubjectID string
	From      LocationXY
	To        LocationXY
}

// Crossing represents location at which Tollgate was crossed.
type Crossing struct {
	SubjectID  string
	TollgateID string
	Location   LocationXY
	Direction  Direction
}

//Direction to North, South, East or West in form of N, S, E, W, NE, NW, SE, SW.
type Direction string

// Detect detects tollgate crossing
// returns nil if no Crossing, otherwise the LocationXY and Direction at which Crossing detected.
func (d *Detector) Detect(movement *Movement) *Crossing {
	for _, d := range d.Tollgates {
		c := detectCrossing(d, movement, 0.0000001)
		if c != nil {
			d := detectDirection(movement)
			c.Direction = d
			return c
		}
	}
	return nil
}

// A little of linear algebra math on lines Crossing.
// Line equation is ax + by + c = 0.
// Line equation can be derived from equation by two points: (y - y1)/(y2 - y1) = (x - x1)/(x2 - x1).
// i.e. x(y2-y1) + y(x1-x2) - x1y2 + y1x2 = 0
// A=y2-y1
// B=x1-x2
// C=y1x2-x1y2
// Point of two lines Crossing can be found by solving a system:
// A1x + B1y + C1 = 0
// A2x + B2y + C2 = 0
// Cramer rule to get point of intersection:
// x = -(C1B2 - B1C2)/(A1B2 - B1A2)
// y = -(A1C2 - C1A2)/(A1B2 - B1A2)
//
// Tollgate - two points representing Tollgate line
// Movement - two points representing Movement line, from previous to current LocationXY
// precision - float greater than 0, i.e. 0.001. for Lat/Long should be 0.0000001
// returns nil if no Crossing, otherwise the location at which Crossing detected.
func detectCrossing(tollgate *Tollgate, movement *Movement, precision float64) *Crossing {
	//Tollgate-representing line
	tx1 := tollgate.Point1.LongitudeX
	ty1 := tollgate.Point1.LatitudeY
	tx2 := tollgate.Point2.LongitudeX
	ty2 := tollgate.Point2.LatitudeY

	A1 := ty2 - ty1
	B1 := tx1 - tx2
	C1 := ty1*tx2 - tx1*ty2

	//Movement-representing line
	mx1 := movement.From.LongitudeX
	my1 := movement.From.LatitudeY
	mx2 := movement.To.LongitudeX
	my2 := movement.To.LatitudeY

	A2 := my2 - my1
	B2 := mx1 - mx2
	C2 := my1*mx2 - mx1*my2

	v := A1*B2 - B1*A2
	if math.Abs(v) > precision {
		x := -(C1*B2 - B1*C2) / v
		y := -(A1*C2 - C1*A2) / v
		return &Crossing{
			SubjectID:  movement.SubjectID,
			TollgateID: tollgate.ID,
			Location: LocationXY{
				LongitudeX: x,
				LatitudeY:  y,
			},
		}
	}
	return nil
}

// When moving to North - latitude increases until 90.
// When moving to South - latitude decreases until -90.
// When moving to East - longitude increases until 180.
// When moving to West - longitude decreases until -180.
// If 90/180 limit crossed it jumps to -90/-180 and vice versa.
// Movement represents a moving subject
// returns Direction in form of N, S, E, W, NE, NW, SE, SW.
func detectDirection(movement *Movement) Direction {
	fromX := movement.From.LongitudeX
	fromY := movement.From.LatitudeY
	toX := movement.To.LongitudeX
	toY := movement.To.LatitudeY

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
