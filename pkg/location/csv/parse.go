package csv

import (
	"strings"
	"time"

	"github.com/openmarketplaceengine/openmarketplaceengine/pkg/util"
)

type Location struct {
	DriverID  string    `csv:"driver_id"`
	Latitude  float64   `csv:"lat"`
	Longitude float64   `csv:"lng"`
	Timestamp time.Time `csv:"bronx_timestamp"`
}

func Parse(line string) (*Location, error) {
	split := strings.Split(line, ",")

	latitude, err := util.ParseLatitude(split[1])
	if err != nil {
		return nil, err
	}

	longitude, err := util.ParseLongitude(split[2])
	if err != nil {
		return nil, err
	}

	ts, err := time.Parse(time.RFC3339, split[3])
	if err != nil {
		return nil, err
	}
	return &Location{
		DriverID:  split[0],
		Latitude:  latitude,
		Longitude: longitude,
		Timestamp: ts,
	}, nil
}
