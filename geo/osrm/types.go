package osrm

import (
	"fmt"

	"github.com/openmarketplaceengine/openmarketplaceengine/geo"
)

// LngLat is Longitude[0] Latitude[1] pair used on OSRM api
// Note here reversed sequence comparing to geoservices.LatLng
type LngLat [2]float64

// Waypoint describes source or destination
type Waypoint struct {
	Hint     string  `json:"hint"`
	Distance float64 `json:"distance"`
	Name     string  `json:"name"`
	Location LngLat  `json:"location"`
}

func (c LngLat) Textual() string {
	return fmt.Sprintf("%v,%v", c[1], c[0])
}

func (c LngLat) ToLatLng() geo.LatLng {
	return geo.LatLng{
		Lat: c[1],
		Lng: c[0],
	}
}
