package osrm

import (
	"context"
	"net/http"

	"github.com/codingsince1985/geo-golang/openstreetmap"
	geosvc "github.com/openmarketplaceengine/openmarketplaceengine/geo"
	"github.com/openmarketplaceengine/openmarketplaceengine/geo/geocode"
)

type Geocoder struct {
	client *http.Client
}

func NewGeocoder(client *http.Client) *Geocoder {
	return &Geocoder{client: client}
}

func (g *Geocoder) ReverseGeocode(ctx context.Context, location geosvc.LatLng) (*geocode.ReverseGeocodeOutput, error) {
	geocoder := openstreetmap.Geocoder()

	address, err := geocoder.ReverseGeocode(location.Lat, location.Lng)
	if err != nil {
		return nil, err
	}

	return &geocode.ReverseGeocodeOutput{
		Address: *address,
	}, nil
}
