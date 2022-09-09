package service

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
	"github.com/openmarketplaceengine/openmarketplaceengine/dao"
	"github.com/openmarketplaceengine/openmarketplaceengine/dispatch/jobstore"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestController(t *testing.T) {

	err := cfg.Load()
	require.NoError(t, err)

	if !dao.Reds.State.Running() {
		require.NoError(t, dao.Reds.Boot())
	}
	client := dao.Reds.StoreClient

	c := Controller{
		JobStore: jobstore.NewJobStore(client),
	}

	t.Run("testGetJobs200", func(t *testing.T) {
		testGetJobs200(t, c)
	})

	t.Run("testPostJobs200", func(t *testing.T) {
		testPostJobs200(t, c)
	})
}

func testGetJobs200(t *testing.T, c Controller) {
	rr := makeGetJobs(t, c, http.StatusOK)

	var jobsResponse JobsResponse
	err := jobsResponse.Decode(rr.Body)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(jobsResponse.Jobs), 0)
}

func testPostJobs200(t *testing.T, c Controller) {

	jobsPayload := JobsPayload{
		Jobs: []*jobstore.Job{
			{
				ID: "job1",
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
	}

	rr0 := makePostJobs(t, c, jobsPayload, http.StatusOK)
	var jobsResponse JobsResponse
	err := jobsResponse.Decode(rr0.Body)
	require.NoError(t, err)
	require.Equal(t, jobsPayload.Jobs, jobsResponse.Jobs)
}

func makeGetJobs(t *testing.T, c Controller, expectStatus int) *httptest.ResponseRecorder {
	r, err := http.NewRequest("GET", "/jobs", nil)
	require.NoError(t, err)
	routeParams := chi.RouteParams{
		Keys:   []string{},
		Values: []string{},
	}
	return makeRequest(t, c.GetJobs, r, routeParams, expectStatus)
}

func makePostJobs(t *testing.T, c Controller, payload JobsPayload, expectStatus int) *httptest.ResponseRecorder {
	body, _ := json.Marshal(payload)

	r, err := http.NewRequest("PUT", "/jobs", bytes.NewReader(body))
	require.NoError(t, err)
	routeParams := chi.RouteParams{
		Keys:   []string{},
		Values: []string{},
	}
	return makeRequest(t, c.PostJobs, r, routeParams, expectStatus)
}

func makeRequest(t *testing.T, handlerFunc func(http.ResponseWriter, *http.Request), r *http.Request, params chi.RouteParams, expectStatus int) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	rctx := chi.NewRouteContext()
	rctx.URLParams = params

	ctx := context.WithValue(r.Context(), chi.RouteCtxKey, rctx)
	handler := http.HandlerFunc(handlerFunc)

	handler.ServeHTTP(rr, r.WithContext(ctx))

	if status := rr.Code; status != expectStatus {
		response, err := ioutil.ReadAll(bytes.NewReader(rr.Body.Bytes()))
		require.NoError(t, err)
		t.Errorf("%s params=%v status expected(%d) != received(%d) response: %s.", r.URL.String(), params, expectStatus, rr.Code, string(response))
	}
	return rr
}
