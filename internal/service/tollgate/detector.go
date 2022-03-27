package tollgate

import (
	"context"

	"github.com/openmarketplaceengine/openmarketplaceengine/log"
)

type Detector struct {
	tollgates []Tollgate
}

func NewDetector(tollgates []Tollgate) *Detector {
	return &Detector{
		tollgates: tollgates,
	}
}

func (d *Detector) DetectTollgateCrossing(ctx context.Context, movement *Movement) *Crossing {
	for _, tollgate := range d.tollgates {
		crossing, err := tollgate.DetectCrossing(ctx, movement)
		if err != nil {
			log.Errorf("detect tollgate crossing error: %q", err)
			continue
		}
		if crossing != nil {
			return crossing
		}
	}
	return nil
}
