package line

import (
	"context"
	"math"

	"github.com/openmarketplaceengine/openmarketplaceengine/internal/service/tollgate"
)

func NewTollgate(id string, point1 *tollgate.LocationXY, point2 *tollgate.LocationXY) *Tollgate {
	return &Tollgate{
		ID:     id,
		Point1: point1,
		Point2: point2,
	}
}

// Tollgate represents a two points line that Subject crosses.
type Tollgate struct {
	ID     string
	Point1 *tollgate.LocationXY
	Point2 *tollgate.LocationXY
}

// Detect detects tollgate crossing
// returns nil if no Crossing, otherwise the LocationXY and Direction at which Crossing detected.
func (t *Tollgate) Detect(ctx context.Context, movement *tollgate.Movement) (*tollgate.Crossing, error) {
	_ = ctx
	return detectCrossing(t, movement, 0.0000001), nil
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
func detectCrossing(lineTollgate *Tollgate, movement *tollgate.Movement, precision float64) *tollgate.Crossing {
	//Tollgate-representing line
	tx1 := lineTollgate.Point1.LongitudeX
	ty1 := lineTollgate.Point1.LatitudeY
	tx2 := lineTollgate.Point2.LongitudeX
	ty2 := lineTollgate.Point2.LatitudeY

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
		return &tollgate.Crossing{
			SubjectID:  movement.SubjectID,
			TollgateID: lineTollgate.ID,
			Location: &tollgate.LocationXY{
				LongitudeX: x,
				LatitudeY:  y,
			},
			Direction: movement.Direction(),
		}
	}
	return nil
}
