package detector

import (
	"context"

	"github.com/openmarketplaceengine/openmarketplaceengine/internal/service/tollgate/bbox"
	"github.com/openmarketplaceengine/openmarketplaceengine/internal/service/tollgate/line"
	"github.com/openmarketplaceengine/openmarketplaceengine/internal/service/tollgate/model"

	"github.com/openmarketplaceengine/openmarketplaceengine/internal/service/tollgate"
	"github.com/openmarketplaceengine/openmarketplaceengine/internal/service/tollgate/conf"
)

type Detector struct {
	tollgates []*model.Tollgate
	storage   bbox.Storage
}

func NewDetector(ctx context.Context, storage bbox.Storage) (*Detector, error) {
	err := conf.LoadTollgates(ctx)
	if err != nil {
		return nil, err
	}
	tollgates, err := model.QueryAll(ctx, 100)
	if err != nil {
		return nil, err
	}
	return &Detector{
		tollgates: tollgates,
		storage:   storage,
	}, nil
}

func (d *Detector) DetectCrossing(ctx context.Context, movement *tollgate.Movement) (*tollgate.Crossing, error) {
	for _, t := range d.tollgates {
		if t.GateLine != nil {
			crossing := line.DetectCrossing(t.ID, &t.GateLine.Line, movement, 0.0000001)
			if crossing != nil {
				return crossing, nil
			}
		}

		if t.BBoxes != nil {
			required := t.BBoxes.Required
			boxes := t.BBoxes.BBoxes
			crossing, err := bbox.DetectCrossing(ctx, d.storage, t.ID, boxes, required, movement)
			if err != nil {
				return nil, err
			}
			if crossing != nil {
				return crossing, nil
			}
		}
	}
	return nil, nil
}
