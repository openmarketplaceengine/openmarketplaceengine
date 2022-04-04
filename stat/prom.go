package stat

import (
	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
	"github.com/openmarketplaceengine/openmarketplaceengine/srv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const pfxErr = "stat"

var state cfg.State64

//-----------------------------------------------------------------------------

func Boot() error {
	if !state.TryBoot() {
		return state.StateError(pfxErr)
	}
	srv.Http.Handle("/metrics", promhttp.Handler())
	state.SetRunning()
	return nil
}
