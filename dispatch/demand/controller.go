package demand

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/openmarketplaceengine/openmarketplaceengine/dispatch/estimate"
	"github.com/openmarketplaceengine/openmarketplaceengine/dispatch/geoqueue"
	"github.com/openmarketplaceengine/openmarketplaceengine/dispatch/htp"
	"github.com/openmarketplaceengine/openmarketplaceengine/dispatch/validate"
	"net/http"
	"net/url"
	"strconv"
)

const areaKey = "ny"

type Demand struct {
	ID      string `json:"id" validate:"required"`
	PickUp  LatLon `json:"pickUp" validate:"required"`
	DropOff LatLon `json:"dropOff" validate:"required"`
}

type Status struct {
	Status string  `json:"status"`
	Demand *Demand `json:"demand"`
}

func (p *Status) Decode(b *bytes.Buffer) error {
	if err := json.NewDecoder(b).Decode(&p); err != nil {
		return fmt.Errorf("decoding bytes error: %w", err)
	}

	return nil
}

type LatLon struct {
	Lat float64 `json:"lat" validate:"required,latitude"`
	Lon float64 `json:"lon" validate:"required,longitude"`
}

type params struct {
	id           string
	radiusMeters int
	lat          float64
	lon          float64
}

type Demands struct {
	ID     string               `json:"id"`
	Demand []*estimate.Estimate `json:"demand"`
}

func (p *Demands) Decode(b *bytes.Buffer) error {
	if err := json.NewDecoder(b).Decode(&p); err != nil {
		return fmt.Errorf("decoding bytes error: %w", err)
	}

	return nil
}

type Controller struct {
	demandService *Service
}

func NewController(demandService *Service) *Controller {
	return &Controller{demandService: demandService}
}

func (c *Controller) GetDemands(w http.ResponseWriter, r *http.Request) {

	p := requireParams(w, r)
	if p == nil {
		return
	}

	res, err := c.demandService.GetEstimates(r.Context(), areaKey, geoqueue.LatLon{
		Lat: p.lat,
		Lon: p.lon,
	}, p.radiusMeters)
	if err != nil {
		htp.Render500(w, r, err)

		return
	}
	render.JSON(w, r, Demands{
		ID:     p.id,
		Demand: res,
	})
}

func (c *Controller) GetDemand(w http.ResponseWriter, r *http.Request) {

	id := requireId(w, r)
	if id == "" {
		return
	}

	res, err := c.demandService.GetDemand(r.Context(), areaKey, id)
	if err != nil {
		htp.Render500(w, r, err)

		return
	}
	render.JSON(w, r, res)
}

func (c *Controller) DeleteDemand(w http.ResponseWriter, r *http.Request) {

	id := requireId(w, r)
	if id == "" {
		return
	}

	err := c.demandService.DeleteDemand(r.Context(), areaKey, id)
	if err != nil {
		htp.Render500(w, r, err)

		return
	}
	render.JSON(w, r, Status{
		Status: "deleted",
		Demand: nil,
	})
}

func requireId(w http.ResponseWriter, r *http.Request) string {
	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		htp.Render400(w, r, []validate.Error{
			{
				Field:   "id",
				Value:   idParam,
				Message: "wrong id",
				Details: fmt.Sprintf("path parameter required, received: %s", idParam),
			},
		}, nil)

		return ""
	}
	return idParam
}

func (c *Controller) PostDemand(w http.ResponseWriter, r *http.Request) {
	payload := requireValidDemandPayload(w, r)
	if payload == nil {
		return
	}

	err := c.demandService.AddDemand(r.Context(), areaKey, payload)
	if err != nil {
		htp.Render500(w, r, err)

		return
	}
	render.JSON(w, r, Status{
		Status: "created",
		Demand: payload,
	})
}

func requireValidDemandPayload(w http.ResponseWriter, r *http.Request) *Demand {
	var payload Demand
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		htp.Render400(w, r,
			[]validate.Error{{
				Field:   "",
				Value:   "",
				Message: fmt.Sprintf("bad payload %T", Demand{}),
				Details: fmt.Sprintf("decode payload error: %s", err),
			}},
			Demand{
				ID: "",
				PickUp: LatLon{
					Lat: 37.656177,
					Lon: -122.473048,
				},
				DropOff: LatLon{
					Lat: 37.656177,
					Lon: -122.473048,
				},
			})

		return nil
	}

	errors := validate.Struct(payload, "Demand.")
	if errors != nil {
		htp.Render400(w, r, errors, Demand{
			ID: "",
			PickUp: LatLon{
				Lat: 37.656177,
				Lon: -122.473048,
			},
			DropOff: LatLon{
				Lat: 37.656177,
				Lon: -122.473048,
			},
		})

		return nil
	}
	return &payload
}

func requireParams(w http.ResponseWriter, r *http.Request) *params {
	idParam := r.URL.Query().Get("id")
	lonParam := r.URL.Query().Get("lon")
	latParam := r.URL.Query().Get("lat")
	radiusMetersParam := r.URL.Query().Get("radiusMeters")
	errors := validate.Vars(validate.Var{
		Field: idParam,
		Name:  "id",
		Tag:   "required",
	}, validate.Var{
		Field: lonParam,
		Name:  "lon",
		Tag:   "required,longitude",
	}, validate.Var{
		Field: latParam,
		Name:  "lat",
		Tag:   "required,latitude",
	}, validate.Var{
		Field: radiusMetersParam,
		Name:  "radiusMeters",
		Tag:   "required,gt=0",
	})
	if errors != nil {
		v := url.Values{}
		v.Set("id", "driver-123")
		v.Set("lon", "-122.473048")
		v.Set("lat", "37.656177")
		v.Set("radiusMeters", "3000")
		htp.Render400(w, r, errors, v.Encode())

		return nil
	}
	lat, err := strconv.ParseFloat(latParam, 64)
	if err != nil {
		panic(fmt.Errorf("parse lat error: %w", err))
	}
	lon, err := strconv.ParseFloat(lonParam, 64)
	if err != nil {
		panic(fmt.Errorf("parse lon error: %w", err))
	}
	radiusMeters, err := strconv.ParseInt(radiusMetersParam, 10, 0)
	if err != nil {
		panic(fmt.Errorf("parse radiusMeters error: %w", err))
	}

	return &params{
		id:           idParam,
		radiusMeters: int(radiusMeters),
		lat:          lat,
		lon:          lon,
	}
}
