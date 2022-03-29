package bbox

import (
	"context"
	"fmt"

	"github.com/openmarketplaceengine/openmarketplaceengine/internal/service/tollgate"
)

// Tollgate detects if subject LocationXY has travelled through Required number of BBoxes.
// Algorithm requires storing 'moving' state in Storage.
// Required is count of visits required for successful pass through, Required <= []BBox size.
type Tollgate struct {
	TollgateID string
	BBoxes     []*BBox // Overlapping bounded boxes to form a route
	Required   int     // Number of visited boxes required to pass
	storage    Storage
}

// BBox is a bounded box, represents area defined by two longitudes and two latitudes (left,bottom,right,top).
// Longitude and Latitude correspond to linear algebra X and Y axis.
type BBox struct {
	Left   float64
	Bottom float64
	Right  float64
	Top    float64
}

func NewTollgate(tollgateID string, boxes []*BBox, required int, storage Storage) (*Tollgate, error) {
	if required > len(boxes) {
		return nil, fmt.Errorf("required visits must be less than number of boxes")
	}

	return &Tollgate{
		TollgateID: tollgateID,
		BBoxes:     boxes,
		Required:   required,
		storage:    storage,
	}, nil
}

// DetectCrossing checks if LocationXY hits a required number of BBoxes.
// If the criteria is met - a successful tollgate.Crossing will be returned.
// Consumer of DetectCrossing should check returned value for not null, meaning detected fact of crossing the tollgate.
func (d *Tollgate) DetectCrossing(ctx context.Context, movement *tollgate.Movement) (*tollgate.Crossing, error) {
	for i, box := range d.BBoxes {
		inb := inBoundary(box, movement.To)
		if inb {
			err := d.storage.Visit(ctx, d.TollgateID, movement.SubjectID, i)
			if err != nil {
				return nil, fmt.Errorf("set visited err=%v", err)
			}
			visits, err := d.storage.Visits(ctx, d.TollgateID, movement.SubjectID, len(d.BBoxes))
			if err != nil {
				return nil, fmt.Errorf("get visits err=%v", err)
			}
			done := checkVisits(visits, d.Required)
			if done {
				err := d.storage.Del(ctx, d.TollgateID, movement.SubjectID)
				if err != nil {
					return nil, fmt.Errorf("delete visits err=%v", err)
				}
				return &tollgate.Crossing{
					TollgateID: d.TollgateID,
					SubjectID:  movement.SubjectID,
					Location:   movement.To,
					Direction:  movement.Direction(),
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

func inBoundary(box *BBox, location *tollgate.LocationXY) bool {
	x := location.LongitudeX
	y := location.LatitudeY
	return box.Left < x && x < box.Right && box.Bottom < y && y < box.Top
}
