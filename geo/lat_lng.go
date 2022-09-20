package geo

type LatLng struct {
	Lat float64
	Lng float64
}

func NewLoc(lat, lon float64) (loc LatLng) {
	loc.Lat = lat
	loc.Lng = lon
	return
}
