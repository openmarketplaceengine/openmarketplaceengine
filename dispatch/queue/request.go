package queue

import (
	"fmt"
	"github.com/openmarketplaceengine/openmarketplaceengine/pkg/validate"
)

type Request struct {
	ID      string
	PickUp  LatLon
	DropOff LatLon
}

type LatLon struct {
	Lat float64
	Lon float64
}

func ValidateLatLon(requests []*Request) error {
	v := validate.Validator{}
	for _, d := range requests {
		v.ValidateFloat64(fmt.Sprintf("request(%s).pickUp_lat", d.ID), d.PickUp.Lat).Latitude()
		v.ValidateFloat64(fmt.Sprintf("request(%s).pickUp_lon", d.ID), d.PickUp.Lon).Longitude()
		v.ValidateFloat64(fmt.Sprintf("request(%s).dropOff_lat", d.ID), d.DropOff.Lat).Latitude()
		v.ValidateFloat64(fmt.Sprintf("request(%s).dropOff_lon", d.ID), d.DropOff.Lon).Longitude()
	}
	return v.Error()
}
