package stat

import (
	"net/http"

	"github.com/openmarketplaceengine/openmarketplaceengine/srv"
)

var _http List

func bootHttp() {
	const path = "/status"
	slog.Infof("HTTP stats endpoint: %s", path)
	srv.Http.Get(path, httpStat)
}

func httpStat(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	buf := AcquireJSONBuffer(2)
	defer buf.Release()
	listJSON(r.Context(), &_http, buf)
	_, _ = buf.WriteTo(w)
}

func AddStat(name string, help string, stat Func) {
	_http.Add(name, help, stat)
}

func Group(name string, help string) *List {
	return _http.Group(name, help)
}
