package detector

import (
	"context"
	"fmt"
)

// BBox is a bounded box, represents area defined by two longitudes and two latitudes
// left,bottom,right,top - lonMin, latMin, lonMax, latMax
// Longitude and Latitude correspond to linear algebra X and Y axis.
type BBox struct {
	LonMin float64
	LatMin float64
	LonMax float64
	LatMax float64
}

// detectCrossingBBox checks if Location hits a required number of BBoxes,
// meaning subject has travelled through Required number of BBoxes.
// If the criteria is met - a successful tollgate.Crossing will be returned.
// Consumer of detectCrossingBBox should check returned value for not null, meaning detected fact of crossing the tollgate.
// Algorithm requires storing 'moving' state in Storage.
// required is count of visits for successful pass through, required <= []BBox size.
func detectCrossingBBox(ctx context.Context, storage *storage, tollgateID string, bBoxes []*BBox, required int, movement *Movement) (*Crossing, error) {
	for i, box := range bBoxes {
		inb := inBoundary(box, movement.To)
		if inb {
			err := storage.visit(ctx, tollgateID, movement.SubjectID, i)
			if err != nil {
				return nil, fmt.Errorf("set visited err=%v", err)
			}
			visits, err := storage.visits(ctx, tollgateID, movement.SubjectID, len(bBoxes))
			if err != nil {
				return nil, fmt.Errorf("get visits err=%v", err)
			}
			done := checkVisits(visits, required)
			if done {
				err := storage.del(ctx, tollgateID, movement.SubjectID)
				if err != nil {
					return nil, fmt.Errorf("delete visits err=%v", err)
				}
				return &Crossing{
					TollgateID: tollgateID,
					WorkerID:   movement.SubjectID,
					Movement:   movement,
					Direction:  movement.Direction(),
					Alg:        BBoxAlg,
				}, nil
			}
		}
	}
	return nil, nil
}

func checkVisits(visits []int64, required int) bool {
	count := 0
	for i := 0; i < len(visits); i++ {
		count = count + int(visits[i])
		if required <= count {
			return true
		}
	}
	return false
}

func inBoundary(box *BBox, location *Location) bool {
	x := location.Lon
	y := location.Lat
	return box.LonMin < x && x < box.LonMax && box.LatMin < y && y < box.LatMax
}
