package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-chi/render"
	"github.com/openmarketplaceengine/openmarketplaceengine/dispatch/htp"
	"github.com/openmarketplaceengine/openmarketplaceengine/dispatch/jobstore"
	"github.com/openmarketplaceengine/openmarketplaceengine/dispatch/validate"
	"net/http"
)

const areaKey = "ny"

type JobsPayload struct {
	Jobs []*jobstore.Job `json:"jobs"`
}

type JobsResponse struct {
	Count int             `json:"count"`
	Jobs  []*jobstore.Job `json:"jobs"`
}

func (p *JobsResponse) Decode(b *bytes.Buffer) error {
	if err := json.NewDecoder(b).Decode(&p); err != nil {
		return fmt.Errorf("decoding bytes error: %w", err)
	}

	return nil
}

type Controller struct {
	JobStore *jobstore.JobStore
}

func (c *Controller) GetJobs(w http.ResponseWriter, r *http.Request) {

	jobs, err := c.JobStore.GetAll(r.Context(), areaKey)
	if err != nil {
		htp.Render500(w, r, err)

		return
	}
	render.JSON(w, r, JobsResponse{
		Count: len(jobs),
		Jobs:  jobs,
	})
}

func (c *Controller) PostJobs(w http.ResponseWriter, r *http.Request) {

	payload := requireValidAddJobsPayload(w, r)
	if payload == nil {
		return
	}

	err := c.JobStore.StoreMany(r.Context(), areaKey, payload.Jobs)
	if err != nil {
		htp.Render500(w, r, err)

		return
	}
	render.JSON(w, r, JobsResponse{
		Count: len(payload.Jobs),
		Jobs:  payload.Jobs,
	})
}

func requireValidAddJobsPayload(w http.ResponseWriter, r *http.Request) *JobsPayload {
	var payload JobsPayload
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		htp.Render400(w, r,
			[]validate.Error{{
				Field:   "",
				Value:   "",
				Message: fmt.Sprintf("bad payload %T", JobsPayload{}),
				Details: fmt.Sprintf("decode payload error: %s", err),
			}},
			JobsPayload{
				Jobs: []*jobstore.Job{
					{
						ID: "",
						PickUp: jobstore.LatLon{
							Lat: 0,
							Lon: 0,
						},
						DropOff: jobstore.LatLon{
							Lat: 0,
							Lon: 0,
						},
					},
				},
			})

		return nil
	}

	return &payload
}
