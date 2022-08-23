package job

import (
	"github.com/mmcloughlin/geohash"
)

func ToGeoHash(latLon LatLon) string {
	return geohash.EncodeWithPrecision(latLon.Lat, latLon.Lon, 8)
}
