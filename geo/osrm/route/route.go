package route

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/openmarketplaceengine/openmarketplaceengine/geo/osrm"
	"github.com/openmarketplaceengine/openmarketplaceengine/geo/osrm/options"
)

// Annotations true , false (default), nodes , distance , duration , datasources , weight , speed
type Annotations string

const (
	True        Annotations = "true"
	False       Annotations = "false"
	Nodes       Annotations = "nodes"
	Distance    Annotations = "distance"
	Duration    Annotations = "duration"
	Datasources Annotations = "datasources"
	Weight      Annotations = "weight"
	Speed       Annotations = "speed"
)

// Geometries polyline (default), polyline6 , geojson
type Geometries string

const (
	Polyline  Geometries = "polyline"
	Polyline6 Geometries = "polyline6"
	Geojson   Geometries = "geojson"
)

// Overview simplified (default), full , false
type Overview string

const (
	Simplified Overview = "simplified"
	Full       Overview = "full"
)

type Request struct {
	Coordinates      []osrm.LngLat
	Alternatives     bool
	Steps            bool
	Annotations      Annotations
	Geometries       Geometries
	Overview         Overview
	ContinueStraight string
	Waypoints        []int
}

type Route struct {
	Legs []struct {
		Steps    []Step  `json:"steps"`
		Summary  string  `json:"summary"`
		Weight   float64 `json:"weight"`
		Duration float64 `json:"duration"`
		Distance float64 `json:"distance"`
	} `json:"legs"`
	WeightName string  `json:"weight_name"`
	Weight     float64 `json:"weight"`
	Duration   float64 `json:"duration"`
	Distance   float64 `json:"distance"`
}

type Step struct {
	Geometry string `json:"geometry"`
	Maneuver struct {
		BearingAfter  int       `json:"bearing_after"`
		BearingBefore int       `json:"bearing_before"`
		Location      []float64 `json:"location"`
		Type          string    `json:"type"`
	} `json:"maneuver"`
	Mode          string `json:"mode"`
	DrivingSide   string `json:"driving_side"`
	Name          string `json:"name"`
	Intersections []struct {
		In       int       `json:"in"`
		Entry    []bool    `json:"entry"`
		Bearings []int     `json:"bearings"`
		Location []float64 `json:"location"`
	} `json:"intersections"`
	Weight   float64 `json:"weight"`
	Duration float64 `json:"duration"`
	Distance float64 `json:"distance"`
}

type Response struct {
	Code      string          `json:"code"`
	Waypoints []osrm.Waypoint `json:"waypoints"`
	Routes    []Route         `json:"routes"`
}

// Routes will return the fastest routes for RouteRequest.
// http request goes to http://project-osrm.org/docs/v5.23.0/api/#route-service endpoint.
func Routes(c *http.Client, request Request) (*Response, error) {

	coords := "polyline(" + url.PathEscape(osrm.ToPolyline(request.Coordinates)) + ")"

	opts := options.UrlEncode(map[string]interface{}{
		"steps":       request.Steps,
		"overview":    request.Overview,
		"annotations": request.Annotations,
	})

	uri := fmt.Sprintf("https://router.project-osrm.org/route/v1/driving/%s?%s", coords, opts)

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
