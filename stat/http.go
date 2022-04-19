package stat

import (
	"net/http"

	"github.com/openmarketplaceengine/openmarketplaceengine/srv"
	"github.com/openmarketplaceengine/openmarketplaceengine/srv/htp"
)

var _http List

func bootHTTP() {
	const path = "/status"
	slog.Infof("HTTP stats endpoint: %s", path)
	srv.Http.Get(path, httpStat)
}

func httpStat(w http.ResponseWriter, r *http.Request) {
	res := htp.GetRes(w, r)
	defer res.Release()
	buf := AcquireJSONBuffer(4)
	defer buf.Release()
	err := listJSON(res.Ctx, &_http, buf)
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

func AddStat(name string, help string, stat Func) {
	_http.Add(name, help, stat)
}

func Group(name string, help string) *List {
	return _http.Group(name, help)
}
