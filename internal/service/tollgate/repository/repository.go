package repository

import (
	"embed"

	"gopkg.in/yaml.v2"
)

//go:embed tollgates.yaml
var tollgatesFile []byte

var _ embed.FS

type BBox struct {
	ID            string `yaml:"id"`
	Name          string `yaml:"name"`
	Boxes         []Box  `yaml:"boxes"`
	BoxesRequired int    `yaml:"boxes_required"`
}

type Box struct {
	LatMax float64 `yaml:"lat_max"`
	LonMax float64 `yaml:"lon_max"`
	LatMin float64 `yaml:"lat_min"`
	LonMin float64 `yaml:"lon_min"`
}

type Line struct {
	ID     string  `yaml:"id"`
	Name   string  `yaml:"name"`
	LatMax float64 `yaml:"lat_max"`
	LonMax float64 `yaml:"lon_max"`
	LatMin float64 `yaml:"lat_min"`
	LonMin float64 `yaml:"lon_min"`
}

type Tollgates struct {
	BboxTollgates []BBox `yaml:"bbox-tollgates"`
	LineTollgates []Line `yaml:"line-tollgates"`
}

func FindAll() (*Tollgates, error) {
	var tollgates Tollgates
	err := yaml.Unmarshal(tollgatesFile, &tollgates)

	if err != nil {
		return nil, err
	}
	return &tollgates, nil
}
