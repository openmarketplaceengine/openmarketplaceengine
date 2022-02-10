package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
	"github.com/openmarketplaceengine/openmarketplaceengine/log"
)

var started = time.Now()

func main() {
	err := cfg.Load()

	if err != nil {
		fatalf("cannot load config: %s", err)
	}

	if cfg.Server.MustExit() { // help or env requested
		return
	}

	err = log.Init(&cfg.Server.Log)

	if err != nil {
		fatalf("logger init failed: %s", err)
	}

	defer log.SafeSync()

	if log.IsDebug() {
		cfg.DebugContext = log.Debugf
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/uptime", uptime)

	ctx := cfg.Context()

	srv := cfg.Http.CreateServer(ctx, mux)

	srv.ErrorLog = log.NewStdLog(log.LevelError)

	srv.RegisterOnShutdown(func() {
		log.Infof("HTTP server shutdown")
	})

	go func() {
		log.Infof("HTTP starting at %s\n", cfg.Http.Addr)
		err = srv.ListenAndServe()
		ctx.Stop()
	}()

	ctx.WaitDone()

	if err != nil {
		log.Fatalf("%s", err)
	}

	err = cfg.Http.Shutdown(srv)

	if err != nil {
		log.Fatalf("%s", err)
	}

	log.Infof("Done")
}

//-----------------------------------------------------------------------------

func uptime(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		code := http.StatusMethodNotAllowed
		http.Error(w, http.StatusText(code), code)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = fmt.Fprintf(w, "{\"uptime\": %d}", time.Since(started)/time.Second)
}

//-----------------------------------------------------------------------------

func fatalf(format string, args ...interface{}) {
	println(fmt.Sprintf(format, args...))
	os.Exit(1)
}
