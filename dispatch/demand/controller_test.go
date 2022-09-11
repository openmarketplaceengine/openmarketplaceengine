package demand

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
	"github.com/openmarketplaceengine/openmarketplaceengine/dao"
	"github.com/openmarketplaceengine/openmarketplaceengine/dispatch/htp"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestController(t *testing.T) {

	err := cfg.Load()
	require.NoError(t, err)

	if !dao.Reds.State.Running() {
		require.NoError(t, dao.Reds.Boot())
	}

	c := Controller{
		service: NewService(),
	}

	t.Run("testGetDemand200", func(t *testing.T) {
		testGetDemand200(t, c)
	})

	t.Run("testGetDemand400", func(t *testing.T) {
		testGetDemand400(t, c)
	})

	t.Run("testPostDemand200", func(t *testing.T) {
		testPostDemand200(t, c)
	})

	t.Run("testPostDemand400", func(t *testing.T) {
		testPostDemand400(t, c)
	})
}

func testGetDemand200(t *testing.T, c Controller) {

	p := params{
		id:           uuid.NewString(),
		radiusMeters: 3000,
		lat:          37.656177,
		lon:          -122.473048,
	}

	rr0 := makeGetDemand(t, c, p, http.StatusOK)
	var response Estimates
	err := response.Decode(rr0.Body)
	require.NoError(t, err)
	require.Equal(t, p.id, response.ID)
	require.GreaterOrEqual(t, len(response.Estimates), 0)
}

func testGetDemand400(t *testing.T, c Controller) {

	p := params{
		id:           "",
		radiusMeters: 0,
		lat:          0,
		lon:          0,
	}

	rr := makeGetDemand(t, c, p, http.StatusBadRequest)
	var errs htp.ValidationErrors
	err := errs.Decode(rr.Body)
	require.NoError(t, err)

	require.Len(t, errs.Errors, 1)
}

func testPostDemand200(t *testing.T, c Controller) {

	jobs := Jobs{
		Jobs: []Job{
			{
				ID: "job-123",
				PickUp: LatLon{
					Lat: 37.656177,
					Lon: -122.473048,
				},
				DropOff: LatLon{
					Lat: 37.656177,
					Lon: -122.473048,
				},
			},
		},
	}

	rr0 := makePostJobs(t, c, jobs, http.StatusOK)
	var response CrudStatus
	err := response.Decode(rr0.Body)
	require.NoError(t, err)
	require.Len(t, response.Jobs, 1)
	require.Equal(t, jobs.Jobs[0].ID, response.Jobs[0].ID)
}

func testPostDemand400(t *testing.T, c Controller) {

	p := Jobs{
		Jobs: []Job{
			{
				ID: "",
				PickUp: LatLon{
					Lat: 0,
					Lon: 0,
				},
				DropOff: LatLon{
					Lat: 0,
					Lon: 0,
				},
			},
		},
	}

	rr := makePostJobs(t, c, p, http.StatusBadRequest)
	var errs htp.ValidationErrors
	err := errs.Decode(rr.Body)
	require.NoError(t, err)

	require.Len(t, errs.Errors, 5)
}

func makeGetDemand(t *testing.T, c Controller, p params, expectStatus int) *httptest.ResponseRecorder {
	u := url.URL{Path: "/demand"}
	v := url.Values{}
	v.Set("id", p.id)
	v.Set("radiusMeters", fmt.Sprintf("%v", p.radiusMeters))
	v.Add("lon", fmt.Sprintf("%v", p.lon))
	v.Add("lat", fmt.Sprintf("%v", p.lat))
	u.RawQuery = v.Encode()

	r, err := http.NewRequest("GET", u.String(), nil)

	require.NoError(t, err)
	routeParams := chi.RouteParams{
		Keys:   []string{},
		Values: []string{},
	}
	return makeRequest(t, c.GetEstimates, r, routeParams, expectStatus)
}

func makePostJobs(t *testing.T, c Controller, jobs Jobs, expectStatus int) *httptest.ResponseRecorder {
	body, _ := json.Marshal(jobs)

	r, err := http.NewRequest("POST", "/demand", bytes.NewReader(body))
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
