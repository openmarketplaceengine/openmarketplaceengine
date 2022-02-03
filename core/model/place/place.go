package place

import (
	"github.com/openmarketplaceengine/openmarketplaceengine/core/model/location"
)

// Place represents geographic locations, points, venues, etc.
type Place struct {
	Name     string
	Location location.Location
}
