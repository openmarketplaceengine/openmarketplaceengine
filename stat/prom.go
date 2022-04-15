package stat

import (
	"github.com/openmarketplaceengine/openmarketplaceengine/srv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

//-----------------------------------------------------------------------------

func bootProm() {
	srv.Http.Handle("/metrics", promhttp.Handler())
}
