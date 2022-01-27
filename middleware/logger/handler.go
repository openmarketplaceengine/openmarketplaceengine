package logger

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type Logger struct {
	handler http.Handler
}

func (l *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	l.handler.ServeHTTP(w, r)
	logger.Info("served",
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
		zap.Duration("duration", time.Since(start)),
	)
}

func NewLogger(handlerToWrap http.Handler) *Logger {
	return &Logger{handlerToWrap}
}
