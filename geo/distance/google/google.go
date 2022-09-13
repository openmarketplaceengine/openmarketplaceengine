package google

import (
	"context"

	"github.com/openmarketplaceengine/openmarketplaceengine/geo/distance"
	"github.com/openmarketplaceengine/openmarketplaceengine/geo/geocode"
	"github.com/openmarketplaceengine/openmarketplaceengine/geo/geocode/google"
	"googlemaps.github.io/maps"
)

func Matrix(ctx context.Context, c *maps.Client, input distance.MatrixPointsInput) (*distance.MatrixOutput, error) {

	// Batch reverse-geocode all locations
	geocoder := google.NewGeocoder(c)
	parallelizationFactor := 10
	geocodeOut, err := geocode.BatchReverseGeocode(
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

	matrix, err := MatrixFromPlaces(ctx, c, distance.MatrixPlacesInput{
		Origins:      origins,
		Destinations: destinations,
	})
	if err != nil {
		return nil, err
	}

	return matrix, nil
}

func MatrixFromPlaces(ctx context.Context, c *maps.Client, input distance.MatrixPlacesInput) (*distance.MatrixOutput, error) {
	var origins []string
	for _, placeID := range input.Origins {
		origins = append(origins, "place_id:"+placeID)
	}
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

func toMatrixOutput(response *maps.DistanceMatrixResponse) *distance.MatrixOutput {
	var rows []distance.MatrixElementsRow
	for i := range response.Rows {
		row := response.Rows[i]
		var elements []distance.MatrixElement
		for j := range row.Elements {
			elem := row.Elements[j]
			elements = append(elements, toElem(elem))
		}
		rows = append(rows, distance.MatrixElementsRow{Elements: elements})
	}
	return &distance.MatrixOutput{
		OriginAddresses:      response.OriginAddresses,
		DestinationAddresses: response.DestinationAddresses,
		Rows:                 rows,
	}
}

func toElem(element *maps.DistanceMatrixElement) distance.MatrixElement {
	return distance.MatrixElement{
		Status:            element.Status,
		Duration:          element.Duration,
		DurationInTraffic: element.DurationInTraffic,
		Distance:          element.Distance.Meters,
	}
}
