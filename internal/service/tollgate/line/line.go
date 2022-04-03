package line

import (
	"math"

	"github.com/openmarketplaceengine/openmarketplaceengine/internal/service/tollgate"
)

// DetectCrossing function detects if Line was crossed by tollgate.Movement,
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
// Movement - two points representing Movement line, from previous to current Location
// precision - float greater than 0, i.e. 0.001. for Lat/Long should be 0.0000001
// returns nil if no Crossing, otherwise the location at which Crossing detected.
func DetectCrossing(tollgateID string, line *tollgate.Line, movement *tollgate.Movement, precision float64) *tollgate.Crossing {
	//Tollgate-representing line
	tx1 := line.Lon1
	ty1 := line.Lat1
	tx2 := line.Lon2
	ty2 := line.Lat2

	A1 := ty2 - ty1
	B1 := tx1 - tx2
	C1 := ty1*tx2 - tx1*ty2

	//Movement-representing line
	mx1 := movement.From.Lon
	my1 := movement.From.Lat
	mx2 := movement.To.Lon
	my2 := movement.To.Lat

	A2 := my2 - my1
	B2 := mx1 - mx2
	C2 := my1*mx2 - mx1*my2

	v := A1*B2 - B1*A2
	if math.Abs(v) > precision {
		x := -(C1*B2 - B1*C2) / v
		y := -(A1*C2 - C1*A2) / v
		return &tollgate.Crossing{
			SubjectID:  movement.SubjectID,
			TollgateID: tollgateID,
			Location: &tollgate.Location{
				Lon: x,
				Lat: y,
			},
			Direction: movement.Direction(),
			Alg:       tollgate.LineAlg,
		}
	}
	return nil
}
