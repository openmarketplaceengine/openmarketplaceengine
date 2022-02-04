package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/openmarketplaceengine/openmarketplaceengine/api/health"
	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
	"github.com/openmarketplaceengine/openmarketplaceengine/config"
	"github.com/openmarketplaceengine/openmarketplaceengine/middleware/loghandler"
)

func main() {
	err := config.Read()

	if err != nil {
		log.Fatalf("read config err=%s", err)
	}

	port := config.GetString(config.ServicePort)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", health.ServeHTTP)

	loggedMux := loghandler.NewLogHandler(mux)

	log.Printf("server is listening at %s\n", port)

	cfg.DebugContext = log.Printf

	ctx := cfg.Context()

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: loggedMux,
		BaseContext: func(listener net.Listener) context.Context {
			return ctx
		},
	}

	go func() {
		err = srv.ListenAndServe()
		ctx.Stop()
	}()

	ctx.WaitDone()

	if err != nil {
		log.Fatalln(err)
	}

	end, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Second))

	err = srv.Shutdown(end)

	cancel()

	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Done")

}
