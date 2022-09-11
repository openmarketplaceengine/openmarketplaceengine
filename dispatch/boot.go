package dispatch

import (
	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
	"github.com/openmarketplaceengine/openmarketplaceengine/dao"
	"github.com/openmarketplaceengine/openmarketplaceengine/dispatch/demand"
	"github.com/openmarketplaceengine/openmarketplaceengine/dispatch/jobstore"
	"github.com/openmarketplaceengine/openmarketplaceengine/dispatch/service"
	"github.com/openmarketplaceengine/openmarketplaceengine/log"
	"github.com/openmarketplaceengine/openmarketplaceengine/srv"
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

	state.SetRunning()
	return nil
}
