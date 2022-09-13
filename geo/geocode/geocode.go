package geocode

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/codingsince1985/geo-golang"
	geosvc "github.com/openmarketplaceengine/openmarketplaceengine/geo"
	"golang.org/x/sync/errgroup"
)

// ReverseGeocodeOutput is a representation of a reverse-geocode request.
// It should be generic enough to work for various APIs:
// https://github.com/codingsince1985/geo-golang
type ReverseGeocodeOutput struct {
	PlaceID string

	// Embedded Address struct from codingsince1985/geo-golang
	geo.Address
}

type ReverseGeocoder interface {
	ReverseGeocode(ctx context.Context, location geosvc.LatLng) (*ReverseGeocodeOutput, error)
}

// BatchReverseGeocode reverse-geocodes a list of geographic coordinates.
//
// It uses parallelization with the errgroup package, since some APIs (e.g.,
// Google Maps) do not offer a way to reverse geocode in bulk, and since each
// individual request would take roughly 150 ms.
// Inspired by:
// https://www.fullstory.com/blog/why-errgroup-withcontext-in-golang-server-handlers/
func BatchReverseGeocode(
	ctx context.Context,
	reverseGeocoder ReverseGeocoder,
	locations []geosvc.LatLng,
	parallelizationFactor int,
) ([]*ReverseGeocodeOutput, error) {
	g, ctx := errgroup.WithContext(ctx)

	locationsChan := make(chan geosvc.LatLng)

	// Step 1: Produce
	g.Go(func() error {
		defer close(locationsChan)
		for _, location := range locations {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case locationsChan <- location:
			}
		}
		return nil
	})

	type Result struct {
		LatLng           geosvc.LatLng
		GeocodingResults *ReverseGeocodeOutput
	}
	results := make(chan Result)

	// Step 2: Map
	nWorkers := parallelizationFactor
	workers := int32(nWorkers)
	for i := 0; i < nWorkers; i++ {
		g.Go(func() error {
			defer func() {
				// Last one out closes shop
				if atomic.AddInt32(&workers, -1) == 0 {
					close(results)
				}
			}()

			for location := range locationsChan {
				geocodingResults, err := reverseGeocoder.ReverseGeocode(ctx, location)
				if err != nil {
					return fmt.Errorf("failed to reverse geocode location: %w", err)
				} else {
					result := Result{
						LatLng: geosvc.LatLng{
							Lat: location.Lat,
							Lng: location.Lng,
						},
						GeocodingResults: geocodingResults,
					}
					select {
					case <-ctx.Done():
						return ctx.Err()
					case results <- result:
					}
				}
			}
			return nil
		})
	}

	// Step 3: Reduce
	// A normal Go map isn't thread-safe, so we use sync.Map
	ret := new(sync.Map)
	g.Go(func() error {
		for result := range results {
			ret.Store(result.LatLng, result.GeocodingResults)
		}
		return nil
	})

	// Wait blocks until all function calls from the Go method have returned, then
	// returns the first non-nil error (if any) from them.
	err := g.Wait()
	if err != nil {
		return nil, err
	}

	// Step 4: Convert sync.Map into a list
	// The order of outputs has to correspond with the order of inputs
	// e.g., if the input was [point1, point2] then the output should be
	// [place1ID, place2ID]
	var out []*ReverseGeocodeOutput
	for _, location := range locations {
		val, ok := ret.Load(location)
		if !ok {
			return nil, fmt.Errorf("could not find lat/lng for (%s, %s)",
				strconv.FormatFloat(location.Lat, 'f', -1, 64),
				strconv.FormatFloat(location.Lng, 'f', -1, 64),
			)
		}
		geocodingResults, ok := val.(*ReverseGeocodeOutput)
		if !ok {
			return nil, fmt.Errorf("expected sync.Map values to be *ReverseGeocodeOutput, not %T", val)
		}

		out = append(out, geocodingResults)
	}
	return out, nil
}
