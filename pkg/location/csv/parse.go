package csv

import (
	"strings"
	"time"

	"github.com/openmarketplaceengine/openmarketplaceengine/pkg/util"
)

type Location struct {
	DriverID  string    `csv:"driver_id"`
	Lat       float64   `csv:"lat"`
	Lon       float64   `csv:"lng"`
	Timestamp time.Time `csv:"bronx_timestamp"`
}

func Parse(line string) (*Location, error) {
	split := strings.Split(line, ",")

	lat, err := util.ParseLatitude(split[1])
	if err != nil {
		return nil, err
	}

	lon, err := util.ParseLongitude(split[2])
	if err != nil {
		return nil, err
	}

	ts, err := time.Parse(time.RFC3339, split[3])
	if err != nil {
		return nil, err
	}
	return &Location{
		DriverID:  split[0],
		Lat:       lat,
		Lon:       lon,
		Timestamp: ts,
	}, nil
}
