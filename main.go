package main

import (
	"fmt"
	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
	"github.com/openmarketplaceengine/openmarketplaceengine/dao"
	"github.com/openmarketplaceengine/openmarketplaceengine/dispatch"
	"github.com/openmarketplaceengine/openmarketplaceengine/dom"
	_ "github.com/openmarketplaceengine/openmarketplaceengine/internal/controller" // gRPC services registration
	"github.com/openmarketplaceengine/openmarketplaceengine/log"
	"github.com/openmarketplaceengine/openmarketplaceengine/srv"
	"github.com/openmarketplaceengine/openmarketplaceengine/stat"
	"os"
)

var boot cfg.BootList

func init() {
	boot.Add("DATA", dom.Boot, nil)
	boot.Use("PGDB", dao.Pgdb)
	boot.Use("REDS", dao.Reds)
	boot.Add("STAT", stat.Boot, nil)
	boot.Add("DISP", dispatch.Boot, nil)
	boot.Use("GRPC", srv.Grpc)
	boot.Use("HTTP", srv.Http)
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

func fatalf(format string, args ...interface{}) {
	// we need to print to stdout for the cloud provider proper capture
	fmt.Printf(format, args...)
	os.Exit(1)
}
