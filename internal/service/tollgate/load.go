package tollgate

import (
	"context"
	"fmt"

	"github.com/openmarketplaceengine/openmarketplaceengine/dom/tollgate"

	"github.com/openmarketplaceengine/openmarketplaceengine/dom/tollgate"

	"github.com/openmarketplaceengine/openmarketplaceengine/pkg/detector"
	"github.com/openmarketplaceengine/openmarketplaceengine/pkg/tollgate/yaml"
)

func Load(ctx context.Context) error {
	tollgates, err := yaml.ReadEmbedded()
	if err != nil {
		return fmt.Errorf("read embedded tollgates error: %w", err)
	}
	for _, t := range tollgates {
		_, err := tollgate.CreateIfNotExists(ctx, &tollgate.Tollgate{
			ID:       t.ID,
			Name:     t.Name,
			BBoxes:   transformBoxes(t.BBoxes),
			GateLine: transformLine(t.Line),
		})

		if err != nil {
			return fmt.Errorf("load tollgate=%v error: %w", t, err)
		}
	}

	return nil
}

func transformBoxes(bBoxes yaml.BBoxes) *tollgate.BBoxes {
	boxes := make([]*detector.BBox, 0)
	for _, b := range bBoxes.Boxes {
		boxes = append(boxes, &detector.BBox{
			LonMin: b.LonMin,
			LatMin: b.LatMin,
			LonMax: b.LonMax,
			LatMax: b.LatMax,
		})
	}
	return &tollgate.BBoxes{
		BBoxes:   boxes,
		Required: bBoxes.BoxesRequired,
	}
}

func transformLine(l yaml.Line) *tollgate.GateLine {
	return &tollgate.GateLine{
		Line: &detector.Line{
			Lon1: l.Lon1,
			Lat1: l.Lat1,
			Lon2: l.Lon2,
			Lat2: l.Lat2,
		},
	}
}
