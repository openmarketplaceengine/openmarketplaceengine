package stat

import (
	"github.com/driverscooperative/geosrv/srv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

//-----------------------------------------------------------------------------

func bootProm() {
	const path = "/metrics"
	slog.Infof("Prometheus endpoint: %s", path)
	srv.Http.Handle(path, promhttp.Handler())
}
