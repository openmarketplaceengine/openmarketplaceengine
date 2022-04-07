package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/openmarketplaceengine/openmarketplaceengine/internal/service/location"
	"github.com/openmarketplaceengine/openmarketplaceengine/internal/service/tollgate/crossing"
	"github.com/openmarketplaceengine/openmarketplaceengine/internal/service/worker"

	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
	"github.com/openmarketplaceengine/openmarketplaceengine/dao"
	"github.com/openmarketplaceengine/openmarketplaceengine/dom"
	"github.com/openmarketplaceengine/openmarketplaceengine/log"
	"github.com/openmarketplaceengine/openmarketplaceengine/srv"
	"github.com/openmarketplaceengine/openmarketplaceengine/stat"
)

var started = time.Now()

var boot cfg.BootList

func init() {
	boot.Add("DATA", dom.Boot, nil)
	boot.Use("PGDB", dao.Pgdb)
	boot.Use("REDS", dao.Reds)
	boot.Add("URLS", route, nil)
	boot.Add("STAT", stat.Boot, nil)
	boot.Use("HTTP", srv.Http)
	boot.Use("GRPC", srv.Grpc)

	location.GrpcRegister()
	crossing.GrpcRegister()
	worker.GrpcRegister()
}

func main() {
	err := cfg.Load()

	if err != nil {
		fatalf("cannot load config: %s\n", err)
	}

	if cfg.Server.MustExit() { // help or env requested
		return
	}

	cfg.Server.ReleaseMemory()

	err = log.Init(cfg.Log())

	if err != nil {
		fatalf("logger init failed: %s\n", err)
	}

	defer log.SafeSync()

	if log.IsDebug() {
		cfg.DebugContext = log.Debugf
	}

	boot.SetLog(log.Infof, log.Errorf)

	err = boot.Boot()

	if err != nil {
		log.Fatalf("%s", err)
	}

	ctx := cfg.Context()

	ctx.WaitDone()

	err = boot.Stop()

	if err != nil {
		log.Fatalf("STOP failed")
	}

	log.Infof("Done")
}

//-----------------------------------------------------------------------------

func route() error {
	if name := cfg.Http.Name; len(name) > 0 {
		srv.Http.SetHeader("Server", name)
	}
	srv.Http.Get("/uptime", uptime)
	return nil
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
	// we need to print to stdout for the cloud provider proper capture
	fmt.Printf(format, args...)
	os.Exit(1)
}
