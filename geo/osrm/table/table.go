package table

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/openmarketplaceengine/openmarketplaceengine/geo/osrm"
	"github.com/openmarketplaceengine/openmarketplaceengine/geo/osrm/options"
)

type Annotations string

const (
	Duration         Annotations = "duration"
	Distance         Annotations = "distance"
	DurationDistance Annotations = "duration,distance"
)

type Request struct {
	Coordinates  []osrm.LngLat
	Origins      []int
	Destinations []int
	Annotations  Annotations
	// FallbackSpeed      *float64
	// FallbackCoordinate *osrm.LngLat
	// ScaleFactor        *float64
}

// Response object from OSRM table request
// Durations - array of arrays that stores the matrix in row-major order.
// durations[i][j] gives the travel time from the i-th waypoint to the j-th waypoint.
// Values are given in seconds. Can be null if no route between i and j can be found.
//
// Distances - array of arrays that stores the matrix in row-major order.
// distances[i][j] gives the travel distance from the i-th waypoint to the j-th waypoint.
// Values are given in meters. Can be null if no route between i and j can be found.
//
// Sources - array of Waypoint objects describing all sources in order
// Destinations -  array of Waypoint objects describing all destinations in order
type Response struct {
	Code         string          `json:"code"`
	Distances    [][]float64     `json:"distances"`
	Durations    [][]float64     `json:"durations"`
	Destinations []osrm.Waypoint `json:"destinations"`
	Sources      []osrm.Waypoint `json:"sources"`
}

// Table will return the fastest route between TableRequest Origins and each Destination.
// i.e. For worker to get ranked routes to available jobs.
// http request goes to http://project-osrm.org/docs/v5.23.0/api/#table-service endpoint.
func Table(c *http.Client, request Request) (*Response, error) {

	coords := "polyline(" + url.PathEscape(osrm.ToPolyline(request.Coordinates)) + ")"

	opts := options.UrlEncode(map[string]interface{}{
		"sources":      request.Origins,
		"destinations": request.Destinations,
		"annotations":  request.Annotations,
	})
	uri := fmt.Sprintf("https://router.project-osrm.org/table/v1/driving/%s?%s", coords, opts)

	res, err := c.Get(uri)
	defer res.Body.Close()
	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusBadRequest {
		return nil, fmt.Errorf("unexpected response: %q status: %d", bytes, res.StatusCode)
	}

	var response Response
	err = json.Unmarshal(bytes, &response)
	if err != nil {
		return nil, fmt.Errorf("unmarshal response: %q error: %w", bytes, err)
	}

	return &response, nil
}
