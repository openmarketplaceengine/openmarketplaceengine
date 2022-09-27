package osrm

import (
	"fmt"
	"net/http"
	"time"

	"github.com/driverscooperative/geosrv/geo/distance"
	"github.com/driverscooperative/geosrv/geo/osrm"
	"github.com/driverscooperative/geosrv/geo/osrm/table"
)

func Matrix(c *http.Client, input distance.MatrixPointsInput) (*distance.MatrixOutput, error) {
	req := toTableRequest(input)
	res, err := table.Table(c, req)
	if err != nil {
		return nil, fmt.Errorf("OSRM table request error: %w", err)
	}
	return toMatrixOutput(res), nil
}

func toTableRequest(input distance.MatrixPointsInput) table.Request {
	var coordinates = make([]osrm.LngLat, 0)
	for _, origin := range input.Origins {
		coordinates = append(coordinates, osrm.LngLat{origin.Lng, origin.Lat})
	}
	for _, destination := range input.Destinations {
		coordinates = append(coordinates, osrm.LngLat{destination.Lng, destination.Lat})
	}
	origins := makeRange(0, len(input.Origins)-1)
	destinations := makeRange(len(input.Origins), len(input.Origins)+len(input.Destinations)-1)
	return table.Request{
		Coordinates:  coordinates,
		Origins:      origins,
		Destinations: destinations,
		Annotations:  table.DurationDistance,
	}
}

func makeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}

func toMatrixOutput(response *table.Response) *distance.MatrixOutput {
	srclen := len(response.Sources)

	rows := make([]distance.MatrixElementsRow, 0, srclen)

	originAddresses := make([]string, 0, srclen)
	destinationAddresses := make([]string, 0, srclen)

	for i, source := range response.Sources {
		originAddresses = append(originAddresses, source.Location.Textual())
		var elements []distance.MatrixElement
		for j, destination := range response.Destinations {
			destinationAddresses = append(destinationAddresses, destination.Location.Textual())
			elements = append(elements, distance.MatrixElement{
				Status:            "",
				Duration:          time.Duration(response.Durations[i][j]) * time.Second,
				DurationInTraffic: 0,
				Distance:          int(response.Distances[i][j] * 1000),
			})
		}
		rows = append(rows, distance.MatrixElementsRow{Elements: elements})
	}

	return &distance.MatrixOutput{
		OriginAddresses:      originAddresses,
		DestinationAddresses: destinationAddresses,
		Rows:                 rows,
	}
}
