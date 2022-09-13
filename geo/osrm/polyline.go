package osrm

import "github.com/twpayne/go-polyline"

// ToPolyline will return coordinates google polyline encoded
// see https://developers.google.com/maps/documentation/utilities/polylineutility
// Output example: _p~iF~ps|U_ulLnnqC_mqNvxq`@
func ToPolyline(lngLats []LngLat) string {
	var coords = make([][]float64, len(lngLats))
	for i, c := range lngLats {
		coords[i] = []float64{c[1], c[0]}
	}
	return string(polyline.EncodeCoords(coords))
}
