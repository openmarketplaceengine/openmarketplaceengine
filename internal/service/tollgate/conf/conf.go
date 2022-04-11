package conf

import (
	"context"
	"embed"

	"github.com/openmarketplaceengine/openmarketplaceengine/internal/service/tollgate/detector"

	"github.com/openmarketplaceengine/openmarketplaceengine/internal/service/tollgate"

	"github.com/openmarketplaceengine/openmarketplaceengine/log"

	"gopkg.in/yaml.v2"
)

//go:embed tollgates.yaml
var tollgatesFile []byte

var _ embed.FS

type BBoxes struct {
	Boxes         []BBox `yaml:"boxes"`
	BoxesRequired int32  `yaml:"min_boxes_required"`
}

type BBox struct {
	LonMin float64 `yaml:"lon_min"`
	LatMin float64 `yaml:"lat_min"`
	LonMax float64 `yaml:"lon_max"`
	LatMax float64 `yaml:"lat_max"`
}

type Line struct {
	Lon1 float64 `yaml:"lon1"`
	Lat1 float64 `yaml:"lat1"`
	Lon2 float64 `yaml:"lon2"`
	Lat2 float64 `yaml:"lat2"`
}

type Tollgate struct {
	ID     string `yaml:"id"`
	Name   string `yaml:"name"`
	BBoxes BBoxes `yaml:"bounding_boxes"`
	Line   Line   `yaml:"gate_line"`
}

func LoadTollgates(ctx context.Context) error {
	var tollgates []Tollgate
	err := yaml.Unmarshal(tollgatesFile, &tollgates)

	if err != nil {
		return err
	}

	for _, t := range tollgates {
		_, err := tollgate.CreateIfNotExists(ctx, &tollgate.Tollgate{
			ID:       t.ID,
			Name:     t.Name,
			BBoxes:   transformBoxes(t.BBoxes),
			GateLine: transformLine(t.Line),
		})

		if err != nil {
			log.Errorf("LoadTollgates load tollgate=%v error: %w", t, err)
		}
	}

	return nil
}

func transformBoxes(bBoxes BBoxes) *tollgate.BBoxes {
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

func transformLine(l Line) *tollgate.GateLine {
	return &tollgate.GateLine{
		Line: &detector.Line{
			Lon1: l.Lon1,
			Lat1: l.Lat1,
			Lon2: l.Lon2,
			Lat2: l.Lat2,
		},
	}
}
