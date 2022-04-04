package line

import (
	"math"

	"github.com/openmarketplaceengine/openmarketplaceengine/internal/util"

	"github.com/openmarketplaceengine/openmarketplaceengine/internal/service/tollgate"
)

func DetectCrossing(tollgateID string, line *tollgate.Line, movement *tollgate.Movement) *tollgate.Crossing {
	return detectCrossingVector(tollgateID, line, movement)
}

//NB THIS LINE CROSSING ALGORITHM CANNOT BE USED, BECAUSE IT DETECTS LINE CROSSING,
//BUT WE NEED TO DETECT IF INTERVALS/SEGMENTS BETWEEN TWO POINTS INTERSECT.
// See detectCrossingVector for proper implementation.
//
// detectCrossingLine function detects if Line was crossed by tollgate.Movement,
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
// precision - float greater than 0, i.e. 0.001. for Lat/Long should be 0.00001
// returns nil if no Crossing, otherwise the location at which Crossing detected.
func detectCrossingLine(tollgateID string, line *tollgate.Line, movement *tollgate.Movement, precision float64) *tollgate.Crossing {
	//Tollgate-representing line
	tx1 := util.Round6(line.Lon1)
	ty1 := util.Round6(line.Lat1)
	tx2 := util.Round6(line.Lon2)
	ty2 := util.Round6(line.Lat2)

	A1 := ty2 - ty1
	B1 := tx1 - tx2
	//C1 := ty1*tx2 - tx1*ty2

	//Movement-representing line
	mx1 := util.Round6(movement.From.Lon)
	my1 := util.Round6(movement.From.Lat)
	mx2 := util.Round6(movement.To.Lon)
	my2 := util.Round6(movement.To.Lat)

	A2 := my2 - my1
	B2 := mx1 - mx2
	//C2 := util.Round6(my1*mx2 - mx1*my2)

	v := A1*B2 - B1*A2
	if math.Abs(v) > precision {
		//x := -(C1*B2 - B1*C2) / v
		//y := -(A1*C2 - C1*A2) / v
		return &tollgate.Crossing{
			WorkerID:   movement.SubjectID,
			TollgateID: tollgateID,
			Movement:   movement,
			Direction:  movement.Direction(),
			Alg:        tollgate.LineAlg,
		}
	}
	return nil
}

// detectCrossingVector detects two line segment crossing using vector math.
func detectCrossingVector(tollgateID string, line *tollgate.Line, movement *tollgate.Movement) *tollgate.Crossing {
	//Tollgate-representing line
	tx1 := util.Round6(line.Lon1)
	ty1 := util.Round6(line.Lat1)
	tx2 := util.Round6(line.Lon2)
	ty2 := util.Round6(line.Lat2)

	//Movement-representing line
	mx1 := util.Round6(movement.From.Lon)
	my1 := util.Round6(movement.From.Lat)
	mx2 := util.Round6(movement.To.Lon)
	my2 := util.Round6(movement.To.Lat)

	s1 := segment{
		start: point{
			x: tx1,
			y: ty1,
		},
		end: point{
			x: tx2,
			y: ty2,
		},
	}
	s2 := segment{
		start: point{
			x: mx1,
			y: my1,
		},
		end: point{
			x: mx2,
			y: my2,
		},
	}
	intersect := intersects(s1, s2)
	if intersect {
		return &tollgate.Crossing{
			WorkerID:   movement.SubjectID,
			TollgateID: tollgateID,
			Movement:   movement,
			Direction:  movement.Direction(),
			Alg:        tollgate.VectorAlg,
		}
	}
	return nil
}

type point struct {
	x float64
	y float64
}

type vector struct {
	start      point
	end        point
	xComponent float64
	yComponent float64
}

func (v *vector) calcComponent() {
	v.xComponent = v.end.x - v.start.x
	v.yComponent = v.end.y - v.start.y
}

type segment struct {
	start point
	end   point
}

func crossProduct(v1, v2 vector) float64 {
	v1.calcComponent()
	v2.calcComponent()
	return v1.xComponent*v2.yComponent - v2.xComponent*v1.yComponent
}

func rangeIntersection(a, b, c, d float64) bool {
	if a > b {
		t := a
		a = b
		b = t
	}
	if c > d {
		t := c
		c = d
		d = t
	}
	return math.Max(a, c) <= math.Min(b, d)
}

func boundingBox(ab, cd segment) bool {
	xRangeIntersection := rangeIntersection(ab.start.x, ab.end.x, cd.start.x, cd.end.x)
	yRangeIntersection := rangeIntersection(ab.start.y, ab.end.y, cd.start.y, cd.end.y)
	return xRangeIntersection && yRangeIntersection
}

func intersects(ab, cd segment) bool {
	if !boundingBox(ab, cd) {
		return false
	}
	vAB := vector{start: ab.start, end: ab.end}
	vAC := vector{start: ab.start, end: cd.start}
	vAD := vector{start: ab.start, end: cd.end}
	vCD := vector{start: cd.start, end: cd.end}
	vCA := vector{start: cd.start, end: ab.start}
	vCB := vector{start: cd.start, end: ab.end}

	d1 := crossProduct(vAB, vAC)
	d2 := crossProduct(vAB, vAD)
	d3 := crossProduct(vCD, vCA)
	d4 := crossProduct(vCD, vCB)
	return ((d1 <= 0 && d2 >= 0) || (d1 >= 0 && d2 <= 0)) && ((d3 <= 0 && d4 >= 0) || (d3 >= 0 && d4 <= 0))
}
