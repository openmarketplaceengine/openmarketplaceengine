package detector

import (
	"context"
	"fmt"

	"github.com/google/uuid"
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
func detectCrossingBBox(ctx context.Context, storage Storage, tollgateID string, workerID string, bBoxes []*BBox, required int32, movement *Movement) (*Crossing, error) {
	key := storageKey(tollgateID, workerID)
	size := len(bBoxes)
	for i, box := range bBoxes {
		inb := inBoundary(box, movement.To)
		if inb {
			err := storage.Visit(ctx, key, size, i)
			if err != nil {
				return nil, fmt.Errorf("set visited err=%v", err)
			}
			visits, err := storage.Visits(ctx, key, size)
			if err != nil {
				return nil, fmt.Errorf("get visits err=%v", err)
			}
			done := checkVisits(visits, required)
			if done {
				err = storage.Del(ctx, key)
				if err != nil {
					return nil, fmt.Errorf("delete visits err=%v", err)
				}
				return &Crossing{
					ID:         uuid.NewString(),
					TollgateID: tollgateID,
					WorkerID:   workerID,
					Movement:   movement,
					Direction:  movement.Direction(),
					Alg:        BBoxAlg,
				}, nil
			}
		}
	}
	return nil, nil
}

func storageKey(tollgateID, workerID string) string {
	return fmt.Sprintf("toll-bbox-%s-%s", tollgateID, workerID)
}

func checkVisits(visits []int, required int32) bool {
	var count int32
	for i := 0; i < len(visits); i++ {
		count = count + int32(visits[i])
		if required <= count {
			return true
		}
	}
	return false
}

func inBoundary(box *BBox, location *Location) bool {
	x := location.Longitude
	y := location.Latitude
	return box.LonMin < x && x < box.LonMax && box.LatMin < y && y < box.LatMax
}
