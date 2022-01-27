package loghandler

import (
	"net/http"
	"time"

	"github.com/openmarketplaceengine/openmarketplaceengine/log"
	"go.uber.org/zap"
)

type LogHandler struct {
	handler http.Handler
}

func (l *LogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	l.handler.ServeHTTP(w, r)
	log.GetLogger().Info("served",
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
		zap.Duration("duration", time.Since(start)),
	)
}

func NewLogHandler(handlerToWrap http.Handler) *LogHandler {
	return &LogHandler{handlerToWrap}
}
