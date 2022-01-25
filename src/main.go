package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/openmarketplaceengine/openmarketplaceengine/src/api/health"
	"github.com/openmarketplaceengine/openmarketplaceengine/src/config"
	"github.com/openmarketplaceengine/openmarketplaceengine/src/middleware/logger"
)

func main() {
	err := config.Read()
	if err != nil {
		log.Fatalf("read config err=%s", err)
	}
	port := config.GetString(config.ServicePort)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", health.ServeHTTP)

	loggedMux := logger.NewLogger(mux)
	log.Printf("server is listening at %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), loggedMux))
}
