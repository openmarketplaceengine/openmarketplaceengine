package yaml

import (
	"bytes"
	_ "embed"
	"io"

	"gopkg.in/yaml.v2"
)

//go:embed tollgates.yaml
var embedded []byte

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

func ReadYaml(r io.Reader) ([]Tollgate, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r)
	if err != nil {
		return nil, err
	}

	var tollgates []Tollgate
	err = yaml.Unmarshal(buf.Bytes(), &tollgates)

	if err != nil {
		return nil, err
	}

	return tollgates, nil
}

func ReadEmbedded() ([]Tollgate, error) {
	return ReadYaml(bytes.NewReader(embedded))
}
