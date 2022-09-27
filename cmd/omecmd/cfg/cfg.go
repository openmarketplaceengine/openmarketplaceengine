// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package cfg

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/driverscooperative/geosrv/app"
	"github.com/driverscooperative/geosrv/app/enc/uri"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gopkg.in/yaml.v2"
)

const (
	CfgFile = "omecmd.yaml"
	defAddr = "localhost:8090"
	defDir  = "."
)

const (
	defTimeout = 5
	minTimeout = 0
	maxTimeout = 60
)

type Cmdcfg struct {
	Server  string
	timeout int64 // connection timeout in seconds
	dstDir  string
	debug   bool
}

var _cfg = Cmdcfg{
	Server:  defAddr,
	timeout: defTimeout,
	dstDir:  defDir,
}

var file = &app.Client().Dirs.ConfFile

//-----------------------------------------------------------------------------

func init() {
	args := &app.Client().Args
	args.String("s", "OME gRPC server `address`", srvArg).Option(true).DefVal(defAddr)
	args.BoolVar(&_cfg.debug, "d", false, "Enable debug output")
	args.Dir("dir", "Destination directory `path` for the current command", 0, func(_ context.Context, path string) error {
		_cfg.dstDir = path
		return nil
	}).DefVal(defDir).Option(true)
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
		format = "[DBG] " + format
		log.Printf(format, args...)
	}
}

//-----------------------------------------------------------------------------

func Errorf(format string, args ...interface{}) {
	format = "[ERR] " + format
	log.Printf(format, args...)
}

//-----------------------------------------------------------------------------

func SafeClose(c io.Closer) {
	if c != nil {
		_ = c.Close()
	}
}

//-----------------------------------------------------------------------------

func timeout() time.Duration {
	t := _cfg.timeout
	if t < minTimeout || t > maxTimeout {
		t = defTimeout
	}
	return time.Duration(t) * time.Second
}

//-----------------------------------------------------------------------------

func Dial(ctx context.Context) (*grpc.ClientConn, error) {
	srv := _cfg.Server
	Debugf("connecting to %s", srv)
	timectx, cancel := context.WithTimeout(ctx, timeout())
	con, err := grpc.DialContext(timectx, srv, grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	cancel()
	if err != nil {
		if err == context.DeadlineExceeded {
			err = fmt.Errorf("connection timeout %s", srv)
		}
		return nil, err
	}
	Debugf("connection established")
	return con, nil
}

//-----------------------------------------------------------------------------

func DstDir() (string, error) {
	var err error
	if _cfg.dstDir == defDir {
		_cfg.dstDir, err = os.Getwd()
	}
	if _cfg.debug && err == nil {
		Debugf("target directory: %q", _cfg.dstDir)
	}
	return _cfg.dstDir, err
}
