package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
)

var started = time.Now()

func main() {
	err := cfg.Load()

	if err != nil {
		log.Fatalf("cannot load config: %s\n", err)
	}

	if cfg.Server.MustExit() { // help or env requested
		return
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/uptime", uptime)

	cfg.DebugContext = log.Printf

	ctx := cfg.Context()

	srv := cfg.Http.CreateServer(ctx, mux)

	srv.RegisterOnShutdown(func() {
		log.Println("HTTP server shutdown")
	})

	go func() {
		log.Printf("HTTP starting at %s\n", cfg.Http.Addr)
		err = srv.ListenAndServe()
		ctx.Stop()
	}()

	ctx.WaitDone()

	if err != nil {
		log.Fatalln(err)
	}

	err = cfg.Http.Shutdown(srv)

	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Done")
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
