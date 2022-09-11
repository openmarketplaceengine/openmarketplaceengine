/**
package with geohash related functions

approximate precision references:
6 chars correspond to < 800 meters
7 chars correspond to < 100 meters
8 chars correspond to < 20 meters

https://www.movable-type.co.uk/scripts/geohash.html
*/

package geohash

import (
	"github.com/mmcloughlin/geohash"
)

const (
	Precision20  = 8
	Precision100 = 7
	Precision800 = 6
)

func ToGeoHash(lat float64, lon float64, precision uint) string {
	return geohash.EncodeWithPrecision(lat, lon, precision)
}
