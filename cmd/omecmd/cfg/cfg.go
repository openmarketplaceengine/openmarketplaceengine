// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package cfg

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/openmarketplaceengine/openmarketplaceengine/app"
	"github.com/openmarketplaceengine/openmarketplaceengine/app/enc/uri"
	"gopkg.in/yaml.v2"
)

const (
	CfgFile = "omecmd.yaml"
	defAddr = "localhost:8090"
)

type Cmdcfg struct {
	Server string
	debug  bool
}

var _cfg = Cmdcfg{
	Server: defAddr,
}

var file = &app.Client().Dirs.ConfFile

//-----------------------------------------------------------------------------

func init() {
	args := &app.Client().Args
	args.String("s", "OME gRPC server `address`", srvArg).Option(true).DefVal(defAddr)
	args.BoolVar(&_cfg.debug, "d", false, "Enable debug output")
}

//-----------------------------------------------------------------------------

func srvArg(_ context.Context, arg string) error {
	if _cfg.Server == arg {
		return nil
	}
	host, err := uri.CheckHostPort(arg)
	if err != nil {
		return err
	}
	_cfg.Server = host
	log.Println("using gRPC server:", _cfg.Server)
	return file.Encode(&_cfg, 0, yaml.Marshal)
}

//-----------------------------------------------------------------------------

func MustInit() {
	if file.Valid() {
		err := file.Decode(&_cfg, yaml.Unmarshal)
		if err == nil {
			return
		}
		log.Printf("error loading config: %s\n", err)
	}
	err := file.Encode(&_cfg, 0, yaml.Marshal)
	if err != nil {
		app.Exit(fmt.Errorf("error saving config: %w", err))
	}
}

//-----------------------------------------------------------------------------

func Server(pathElem ...string) string {
	return uri.Join(_cfg.Server, pathElem...)
}

//-----------------------------------------------------------------------------

func Debugf(format string, args ...interface{}) {
	if _cfg.debug {
		log.Printf(format, args...)
	}
}

//-----------------------------------------------------------------------------

func SafeClose(c io.Closer) {
	if c != nil {
		_ = c.Close()
	}
}
