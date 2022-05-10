// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"os"

	"github.com/openmarketplaceengine/openmarketplaceengine/app"
	"github.com/openmarketplaceengine/openmarketplaceengine/cmd/omecmd/cfg"
	_ "github.com/openmarketplaceengine/openmarketplaceengine/cmd/omecmd/status"
)

func init() {
	setlog()
}

func main() {
	if len(os.Args) == 1 {
		println("Missing command line arguments. Use '-h' for help.")
		os.Exit(2)
	}
	cli := app.Client()
	cli.MustInit(app.Domain, cfg.CfgFile)
	cfg.MustInit()
	err := cli.Args.Parse(app.Context(), os.Args[1:])
	app.Exit(err)
}

func setlog() {
	// log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime)
	log.SetPrefix("")
}
