package google

import (
	"context"
	"github.com/openmarketplaceengine/openmarketplaceengine/dispatch/metrics"
	"time"

	"googlemaps.github.io/maps"
)

func Matrix(ctx context.Context, c *maps.Client, input MatrixPointsInput) (*MatrixOutput, error) {
	defer func(begin time.Time) {
		metrics.MatrixApiHits.Inc()
		milliseconds := time.Since(begin).Milliseconds()
		value := float64(milliseconds)
		metrics.MatrixApiCallDuration.Observe(value)
	}(time.Now())

	// Batch reverse-geocode all locations
	geocoder := NewGeocoder(c)
	parallelizationFactor := 10
	geocodeOut, err := BatchReverseGeocode(
		ctx,
		geocoder,
		append(input.Origins, input.Destinations...),
		parallelizationFactor)
	if err != nil {
		return nil, err
	}

	var origins []string
	var destinations []string
	for idx, e := range geocodeOut {
		if idx < len(input.Origins) {
			origins = append(origins, e.PlaceID)
		} else {
			destinations = append(destinations, e.PlaceID)
		}
	}

	matrix, err := MatrixFromPlaces(ctx, c, MatrixPlacesInput{
		Origins:      origins,
		Destinations: destinations,
	})
	if err != nil {
		return nil, err
	}

	return matrix, nil
}

func MatrixFromPlaces(ctx context.Context, c *maps.Client, input MatrixPlacesInput) (*MatrixOutput, error) {
	defer func(begin time.Time) {
		metrics.MatrixApiHits.Inc()
		milliseconds := time.Since(begin).Milliseconds()
		value := float64(milliseconds)
		metrics.MatrixApiCallDuration.Observe(value)
	}(time.Now())

	// nolint:prealloc
	var origins []string
	for _, placeID := range input.Origins {
		origins = append(origins, "place_id:"+placeID)
	}
	// nolint:prealloc
	var destinations []string
	for _, placeID := range input.Destinations {
		destinations = append(destinations, "place_id:"+placeID)
	}
	matrix, err := c.DistanceMatrix(ctx, &maps.DistanceMatrixRequest{
		Origins:      origins,
		Destinations: destinations,
	})
	if err != nil {
		return nil, err
	}
	return toMatrixOutput(matrix), err
}

func toMatrixOutput(response *maps.DistanceMatrixResponse) *MatrixOutput {
	// nolint:prealloc
	var rows []MatrixElementsRow
	for i := range response.Rows {
		row := response.Rows[i]
		var elements []MatrixElement
		for j := range row.Elements {
			elem := row.Elements[j]
			elements = append(elements, toElem(elem))
		}
		rows = append(rows, MatrixElementsRow{Elements: elements})
	}
	return &MatrixOutput{
		OriginAddresses:      response.OriginAddresses,
		DestinationAddresses: response.DestinationAddresses,
		Rows:                 rows,
	}
}

func toElem(element *maps.DistanceMatrixElement) MatrixElement {
	return MatrixElement{
		Status:            element.Status,
		Duration:          element.Duration,
		DurationInTraffic: element.DurationInTraffic,
		Distance:          element.Distance.Meters,
	}
}
