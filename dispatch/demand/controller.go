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
	"github.com/openmarketplaceengine/openmarketplaceengine/dispatch/metrics"
	"github.com/openmarketplaceengine/openmarketplaceengine/dispatch/validate"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const areaKey = "ny"

type Job struct {
	ID      string `json:"id" validate:"required"`
	PickUp  LatLon `json:"pickUp" validate:"required"`
	DropOff LatLon `json:"dropOff" validate:"required"`
}

type Jobs struct {
	Jobs []Job `json:"jobs" validate:"required,dive"`
}

type JobIds struct {
	Ids []string `json:"jobIds" validate:"required"`
}

type CrudStatus struct {
	Status string `json:"status"`
	Jobs   []Job  `json:"jobs"`
}

func (p *CrudStatus) Decode(b *bytes.Buffer) error {
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

type Estimates struct {
	ID        string               `json:"id"`
	Count     int                  `json:"count"`
	Estimates []*estimate.Estimate `json:"estimates"`
}

func (p *Estimates) Decode(b *bytes.Buffer) error {
	if err := json.NewDecoder(b).Decode(&p); err != nil {
		return fmt.Errorf("decoding bytes error: %w", err)
	}

	return nil
}

type Controller struct {
	service *Service
}

func NewController(service *Service) *Controller {
	return &Controller{service: service}
}

func (c *Controller) GetEstimates(w http.ResponseWriter, r *http.Request) {
	defer func(begin time.Time) {
		milliseconds := time.Since(begin).Milliseconds()
		value := float64(milliseconds)
		metrics.EstimatesApiCallDuration.Observe(value)
	}(time.Now())

	p := requireParams(w, r)
	if p == nil {
		return
	}

	res, err := c.service.GetEstimates(r.Context(), areaKey, geoqueue.LatLon{
		Lat: p.lat,
		Lon: p.lon,
	}, p.radiusMeters)
	if err != nil {
		htp.Render500(w, r, err)

		return
	}
	render.JSON(w, r, Estimates{
		ID:        p.id,
		Count:     len(res),
		Estimates: res,
	})
}

func (c *Controller) GetJob(w http.ResponseWriter, r *http.Request) {

	id := requireId(w, r)
	if id == "" {
		return
	}

	res, err := c.service.GetJob(r.Context(), areaKey, id)
	if err != nil {
		htp.Render500(w, r, err)

		return
	}
	render.JSON(w, r, res)
}

func (c *Controller) DeleteOne(w http.ResponseWriter, r *http.Request) {

	id := requireId(w, r)
	if id == "" {
		return
	}

	err := c.service.DeleteJobs(r.Context(), areaKey, id)
	if err != nil {
		htp.Render500(w, r, err)

		return
	}
	render.JSON(w, r, CrudStatus{
		Status: "deleted",
		Jobs:   nil,
	})
}

func (c *Controller) DeleteMany(w http.ResponseWriter, r *http.Request) {

	ids := requireValidJobIds(w, r)
	if ids == nil {
		return
	}

	err := c.service.DeleteJobs(r.Context(), areaKey, ids.Ids...)
	if err != nil {
		htp.Render500(w, r, err)

		return
	}
	render.JSON(w, r, CrudStatus{
		Status: "deleted",
		Jobs:   nil,
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

func (c *Controller) PostJobs(w http.ResponseWriter, r *http.Request) {
	payload := requireValidJobs(w, r)
	if payload == nil {
		return
	}

	for _, job := range payload.Jobs {
		err := c.service.AddJob(r.Context(), areaKey, &job)
		if err != nil {
			htp.Render500(w, r, err)

			return
		}
	}
	render.JSON(w, r, CrudStatus{
		Status: "created",
		Jobs:   payload.Jobs,
	})
}

func requireValidJobs(w http.ResponseWriter, r *http.Request) *Jobs {
	var jobs Jobs
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&jobs)
	if err != nil {
		htp.Render400(w, r,
			[]validate.Error{{
				Field:   "",
				Value:   "",
				Message: fmt.Sprintf("bad payload %T", Jobs{}),
				Details: fmt.Sprintf("decode payload error: %s", err),
			}},
			Job{
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

	errors := validate.Struct(jobs, "Jobs.")
	if errors != nil {
		htp.Render400(w, r, errors, Job{
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
	return &jobs
}

func requireValidJobIds(w http.ResponseWriter, r *http.Request) *JobIds {
	var ids JobIds
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&ids)
	if err != nil {
		htp.Render400(w, r,
			[]validate.Error{{
				Field:   "",
				Value:   "",
				Message: fmt.Sprintf("bad payload %T", JobIds{}),
				Details: fmt.Sprintf("decode payload error: %s", err),
			}},
			JobIds{
				Ids: []string{"job1"},
			},
		)

		return nil
	}

	errors := validate.Struct(ids, "JobIds.")
	if errors != nil {
		htp.Render400(w, r, errors,
			JobIds{
				Ids: []string{"job1"},
			},
		)

		return nil
	}
	return &ids
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
