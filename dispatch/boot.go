package dispatch

import (
	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
	"github.com/openmarketplaceengine/openmarketplaceengine/dao"
	"github.com/openmarketplaceengine/openmarketplaceengine/dispatch/demand"
	"github.com/openmarketplaceengine/openmarketplaceengine/dispatch/jobstore"
	"github.com/openmarketplaceengine/openmarketplaceengine/dispatch/metrics"
	"github.com/openmarketplaceengine/openmarketplaceengine/dispatch/service"
	"github.com/openmarketplaceengine/openmarketplaceengine/log"
	"github.com/openmarketplaceengine/openmarketplaceengine/srv"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const pfxErr = "disp"

var state cfg.State64

var slog = log.Log()

//-----------------------------------------------------------------------------

func Boot() error {
	if !state.TryBoot() {
		return state.StateError(pfxErr)
	}
	slog = log.Named("STAT")
	_ = slog
	c := service.Controller{JobStore: jobstore.NewJobStore(dao.Reds.StoreClient)}
	srv.Http.Get("/jobs", c.GetJobs)
	srv.Http.Post("/jobs", c.PostJobs)

	loc := demand.NewController(demand.NewService())
	srv.Http.Get("/demand", loc.GetEstimates)
	srv.Http.Get("/demand/{id}", loc.GetJob)
	srv.Http.Delete("/demand/{id}", loc.DeleteOne)
	srv.Http.Delete("/demand", loc.DeleteMany)
	srv.Http.Post("/demand", loc.PostJobs)

	r := prometheus.NewRegistry()
	r.MustRegister(metrics.MatrixApiHits.(prometheus.Counter))
	r.MustRegister(metrics.MatrixApiCallDuration.(prometheus.Histogram))
	r.MustRegister(metrics.EstimatesApiCallDuration.(prometheus.Histogram))
	handler := promhttp.HandlerFor(r, promhttp.HandlerOpts{})
	srv.Http.Handle("/demand/metrics", handler)

	state.SetRunning()
	return nil
}
