package ranking

import (
	"github.com/openmarketplaceengine/openmarketplaceengine/geo"
	"github.com/uber/h3-go"
)

// getCell indexes the location at the specified resolution, returning the index
// of the cell containing the location. This buckets the geographic point into
// the H3 grid. See the algorithm description for more information:
// https://h3geo.org/docs/core-library/geoToH3desc
//
// https://h3geo.org/docs/api/indexing#geotoh3
func getCell(l geo.LatLng, res int) string {
	i := h3.FromGeo(h3.GeoCoord{
		Latitude:  l.Lat,
		Longitude: l.Lng,
	}, res)
	return h3.ToString(i)
}

// cellNeighbors produces indices within k distance of the origin index.
// k-ring 0 is defined as the origin index.
// k-ring 1 is defined as k-ring 0 and all neighboring indices, and so on.
//
// https://h3geo.org/docs/api/traversal#kring
func cellNeighbors(l geo.LatLng, res int, k int) []string {
	i := h3.FromGeo(h3.GeoCoord{
		Latitude:  l.Lat,
		Longitude: l.Lng,
	}, res)
	indexes := h3.KRing(i, k)
	var out []string
	for _, idx := range indexes {
		out = append(out, h3.ToString(idx))
	}
	return out
}

func hasNeighbor(neighbors []string, desiredNeighbor string) bool {
	for _, e := range neighbors {
		if e == desiredNeighbor {
			return true
		}
	}
	return false
}
