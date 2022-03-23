package bbox

import (
	"context"
	"fmt"
)

// Detector detects if subject LocationXY has travelled through Required number of BBoxes.
// Algorithm requires storing 'moving' state in Storage.
type Detector struct {
	TollgateID string
	BBoxes     []*BBox // Overlapping bounded boxes to form a route
	Required   int     // Number of visited boxes required to pass
	storage    Storage
}

// LocationXY is longitude latitude corresponding to linear algebra X Y axis.
type LocationXY struct {
	LongitudeX float64
	LatitudeY  float64
}

// BBox is a bounded box, represents area defined by two longitudes and two latitudes (left,bottom,right,top).
// Longitude and Latitude correspond to linear algebra X and Y axis.
type BBox struct {
	Left   float64
	Bottom float64
	Right  float64
	Top    float64
}

// PassedThrough represents detected fact of passing through the tollgate by SubjectID.
type PassedThrough struct {
	TollgateID string
	SubjectID  string
	Visits     []int64
}

func NewDetector(tollgateID string, boxes []*BBox, required int, storage Storage) (*Detector, error) {
	if required > len(boxes) {
		return nil, fmt.Errorf("required visits must be less than number of boxes")
	}

	return &Detector{
		TollgateID: tollgateID,
		BBoxes:     boxes,
		Required:   required,
		storage:    storage,
	}, nil
}

// Trace checks if LocationXY hits a required number of BBoxes.
// If the criteria is met - a successful PassedThrough will be returned.
// Consumer of Trace should check returned value for not null, meaning detected fact of passing through the tollgate.
func (d *Detector) Trace(ctx context.Context, subjectID string, location *LocationXY) (*PassedThrough, error) {
	for i, box := range d.BBoxes {
		inb := inBoundary(box, location)
		if inb {
			err := d.storage.Visit(ctx, d.TollgateID, subjectID, i)
			if err != nil {
				return nil, fmt.Errorf("set visited err=%v", err)
			}
			visits, err := d.storage.Visits(ctx, d.TollgateID, subjectID, len(d.BBoxes))
			if err != nil {
				return nil, fmt.Errorf("get visits err=%v", err)
			}
			done := checkVisits(visits, d.Required)
			if done {
				err := d.storage.Del(ctx, d.TollgateID, subjectID)
				if err != nil {
					return nil, fmt.Errorf("delete visits err=%v", err)
				}
				return &PassedThrough{
					TollgateID: d.TollgateID,
					SubjectID:  subjectID,
					Visits:     visits,
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

func inBoundary(box *BBox, location *LocationXY) bool {
	x := location.LongitudeX
	y := location.LatitudeY
	return box.Left < x && x < box.Right && box.Bottom < y && y < box.Top
}
