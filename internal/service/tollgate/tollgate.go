package tollgate

import "math"

// LocationXY is longitude latitude corresponding to linear algebra X Y axis.
type LocationXY struct {
	longitudeX float64
	latitudeY  float64
}

// Tollgate represents a two points line that driver crosses.
type Tollgate struct {
	ID     string
	Point1 LocationXY
	Point2 LocationXY
}

// Movement represents the latest trace driver moved in form of two points line from previous to current location.
type Movement struct {
	DriverID     string
	PrevLocation LocationXY
	CurrLocation LocationXY
}

// Crossing represents location at which Tollgate was crossed.
type Crossing struct {
	DriverID   string
	TollgateID string
	Location   LocationXY
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
// precision - float greater than 0, i.e. 0.001.
// returns nil if no crossing, otherwise the location at which crossing detected.
func detectCrossing(tollgate *Tollgate, movement *Movement, precision float64) *Crossing {
	//Tollgate-representing line
	tx1 := tollgate.Point1.longitudeX
	ty1 := tollgate.Point1.latitudeY
	tx2 := tollgate.Point2.longitudeX
	ty2 := tollgate.Point2.latitudeY

	A1 := ty2 - ty1
	B1 := tx1 - tx2
	C1 := ty1*tx2 - tx1*ty2

	//Movement-representing line
	mx1 := movement.PrevLocation.longitudeX
	my1 := movement.PrevLocation.latitudeY
	mx2 := movement.CurrLocation.longitudeX
	my2 := movement.CurrLocation.latitudeY

	A2 := my2 - my1
	B2 := mx1 - mx2
	C2 := my1*mx2 - mx1*my2

	v := A1*B2 - B1*A2
	if math.Abs(v) > precision {
		x := -(C1*B2 - B1*C2) / v
		y := -(A1*C2 - C1*A2) / v
		return &Crossing{
			DriverID:   movement.DriverID,
			TollgateID: tollgate.ID,
			Location: LocationXY{
				longitudeX: x,
				latitudeY:  y,
			},
		}
	}
	return nil
}
