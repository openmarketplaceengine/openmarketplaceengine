package stat

import (
	"net/http"

	"github.com/openmarketplaceengine/openmarketplaceengine/srv"
	"github.com/openmarketplaceengine/openmarketplaceengine/srv/htp"
)

var top List
var (
	app = top.Group("ome", "Application Metrics")
	sys = top.Group("sys", "System Metrics")
)

func App() *List {
	return app
}

func Sys() *List {
	return sys
}

//-----------------------------------------------------------------------------

func bootHTTP() {
	const path = "/status"
	slog.Infof("HTTP stats endpoint: %s", path)
	srv.Http.Get(path, httpStat)
}

//-----------------------------------------------------------------------------

func httpStat(w http.ResponseWriter, r *http.Request) {
	res := htp.GetRes(w, r)
	defer res.Release()
	buf := AcquireJSONBuffer(4)
	defer buf.Release()
	err := listJSON(res.Ctx, &top, buf)
	if err != nil {
		slog.Errorf("%s", err)
		res.ServerError()
		return
	}
	res.SetJSON()
	err = res.SendData(buf.End())
	if err != nil {
		slog.Errorf("%s", err)
	}
}
