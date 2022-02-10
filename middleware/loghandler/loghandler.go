package loghandler

import (
	"net/http"
	"time"

	"github.com/openmarketplaceengine/openmarketplaceengine/log"
)

type LogHandler struct {
	handler http.Handler
}

func (l *LogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	l.handler.ServeHTTP(w, r)
	log.Debugf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
}

func NewLogHandler(handlerToWrap http.Handler) *LogHandler {
	return &LogHandler{handlerToWrap}
}
