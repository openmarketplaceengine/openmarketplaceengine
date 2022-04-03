package bbox

import (
	"context"
	"fmt"

	"github.com/openmarketplaceengine/openmarketplaceengine/internal/service/tollgate"
)

// DetectCrossing checks if Location hits a required number of BBoxes,
// meaning subject has travelled through Required number of BBoxes.
// If the criteria is met - a successful tollgate.Crossing will be returned.
// Consumer of DetectCrossing should check returned value for not null, meaning detected fact of crossing the tollgate.
// Algorithm requires storing 'moving' state in Storage.
// required is count of visits for successful pass through, required <= []BBox size.
func DetectCrossing(ctx context.Context, storage Storage, tollgateID string, bBoxes []*tollgate.BBox, required int, movement *tollgate.Movement) (*tollgate.Crossing, error) {
	for i, box := range bBoxes {
		inb := inBoundary(box, movement.To)
		if inb {
			err := storage.Visit(ctx, tollgateID, movement.SubjectID, i)
			if err != nil {
				return nil, fmt.Errorf("set visited err=%v", err)
			}
			visits, err := storage.Visits(ctx, tollgateID, movement.SubjectID, len(bBoxes))
			if err != nil {
				return nil, fmt.Errorf("get visits err=%v", err)
			}
			done := checkVisits(visits, required)
			if done {
				err := storage.Del(ctx, tollgateID, movement.SubjectID)
				if err != nil {
					return nil, fmt.Errorf("delete visits err=%v", err)
				}
				return &tollgate.Crossing{
					TollgateID: tollgateID,
					SubjectID:  movement.SubjectID,
					Location:   movement.To,
					Direction:  movement.Direction(),
					Alg:        tollgate.BBoxAlg,
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

func inBoundary(box *tollgate.BBox, location *tollgate.Location) bool {
	x := location.Lon
	y := location.Lat
	return box.LonMin < x && x < box.LonMax && box.LatMin < y && y < box.LatMax
}
