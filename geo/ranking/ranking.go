package ranking

import (
	"context"
	"log"
	"time"

	"github.com/driverscooperative/geosrv/geo"
)

// Rank returns an ordered list of origin points which are nearest to a
// destination point.
//
// Elements are ordered by descending distance (i.e., nearest points first).
// Traffic or travel duration are not accounted for.
//
// Implementation relies on H3, a system for geospatial indexing.
func Rank(ctx context.Context, origins []geo.LatLng, destination geo.LatLng) []geo.LatLng {
	start := time.Now()

	// TODO configure acceptable resolutions?
	// Geo-hash/index the destination
	r7Cell := getCell(destination, 7)
	r8Cell := getCell(destination, 8)
	r9Cell := getCell(destination, 9)
	r10Cell := getCell(destination, 10)

	var neighborsR10K1, neighborsR10K2 []geo.LatLng
	var neighborsR9K1, neighborsR9K2 []geo.LatLng
	var neighborsR8K1, neighborsR8K2 []geo.LatLng
	var neighborsR7K1 []geo.LatLng

	for _, e := range origins {
		r10K1Neighbors := cellNeighbors(e, 10, 1)
		if hasNeighbor(r10K1Neighbors, r10Cell) {
			neighborsR10K1 = append(neighborsR10K1, e)
			continue
		}
		r10K2Neighbors := cellNeighbors(e, 10, 2)
		if hasNeighbor(r10K2Neighbors, r10Cell) {
			neighborsR10K2 = append(neighborsR10K2, e)
			continue
		}
		r9K1Neighbors := cellNeighbors(e, 9, 1)
		if hasNeighbor(r9K1Neighbors, r9Cell) {
			neighborsR9K1 = append(neighborsR9K1, e)
			continue
		}
		r9K2Neighbors := cellNeighbors(e, 9, 2)
		if hasNeighbor(r9K2Neighbors, r9Cell) {
			neighborsR9K2 = append(neighborsR9K2, e)
			continue
		}
		r8K1Neighbors := cellNeighbors(e, 8, 1)
		if hasNeighbor(r8K1Neighbors, r8Cell) {
			neighborsR8K1 = append(neighborsR8K1, e)
			continue
		}
		r8K2Neighbors := cellNeighbors(e, 8, 2)
		if hasNeighbor(r8K2Neighbors, r8Cell) {
			neighborsR8K2 = append(neighborsR8K2, e)
			continue
		}
		r7K1Neighbors := cellNeighbors(e, 7, 1)
		if hasNeighbor(r7K1Neighbors, r7Cell) {
			neighborsR7K1 = append(neighborsR7K1, e)
			continue
		}
	}

	log.Println(time.Since(start))

	var neighbors []geo.LatLng
	neighbors = append(neighbors, neighborsR10K1...)
	neighbors = append(neighbors, neighborsR10K2...)
	neighbors = append(neighbors, neighborsR9K1...)
	neighbors = append(neighbors, neighborsR9K2...)
	neighbors = append(neighbors, neighborsR8K1...)
	neighbors = append(neighbors, neighborsR8K2...)
	neighbors = append(neighbors, neighborsR7K1...)

	return neighbors
}
