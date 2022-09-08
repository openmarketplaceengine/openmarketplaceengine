package google

import (
	"context"
	"errors"

	"github.com/codingsince1985/geo-golang"
	"googlemaps.github.io/maps"
)

var errNoResults = errors.New("no results found for coordinates")

type Geocoder struct {
	client *maps.Client
}

func NewGeocoder(client *maps.Client) *Geocoder {
	return &Geocoder{client: client}
}

func (g *Geocoder) ReverseGeocode(ctx context.Context, location LatLon) (*ReverseGeocodeOutput, error) {
	results, err := g.client.ReverseGeocode(ctx, &maps.GeocodingRequest{
		LatLng: &maps.LatLng{
			Lat: location.Lat,
			Lng: location.Lon,
		},
		ResultType:   nil,
		LocationType: nil,
		PlaceID:      "",
		Language:     "",
		Custom:       nil,
	})
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, errNoResults
	}
	return toReverseGeocodeOutput(results), err
}

func toReverseGeocodeOutput(in []maps.GeocodingResult) *ReverseGeocodeOutput {
	// Assuming Google Maps API returns places ordered in such a way that
	// the first element is the most salient/relevant.
	bestResult := in[0]
	return &ReverseGeocodeOutput{
		PlaceID: bestResult.PlaceID,
		Address: geo.Address{
			FormattedAddress: bestResult.FormattedAddress,
			// TODO set other components
			Street:        "",
			HouseNumber:   "",
			Suburb:        "",
			Postcode:      "",
			State:         "",
			StateCode:     "",
			StateDistrict: "",
			County:        "",
			Country:       "",
			CountryCode:   "",
			City:          "",
		},
	}
}
