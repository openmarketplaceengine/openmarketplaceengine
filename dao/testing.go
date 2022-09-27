// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package dao

import (
	"github.com/driverscooperative/geosrv/cfg"
	"github.com/driverscooperative/geosrv/log"
)

//-----------------------------------------------------------------------------
// Testing
//-----------------------------------------------------------------------------

type Tester interface {
	Fatalf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Cleanup(f func())
	SkipNow()
}

//-----------------------------------------------------------------------------

func SkipTest() bool {
	_, ok := cfg.GetEnv(cfg.EnvPgdbAddr)
	return !ok
}

//-----------------------------------------------------------------------------

func WillTest(t Tester, schema string) (ctx Context) {
	defer Pgdb.clearAutos()
	ctx = cfg.Context()
	if SkipTest() {
		t.SkipNow()
		return
	}
	Pgdb.resetForTests(t)
	if len(schema) > 0 {
		err := cfg.SetEnv(cfg.EnvPgdbSchema, schema)
		if err != nil {
			t.Fatalf("setenv %q=%q failed: %s", cfg.EnvPgdbSchema, schema, err)
		}
	}
	err := cfg.Load()
	if err != nil {
		t.Fatalf("config load failed: %s", err)
	}
	err = log.Init(log.DevelConfig().WithTrace(false).WithCaller(false))
	if err != nil {
		t.Fatalf("log init failed: %s", err)
	}
	err = Pgdb.Boot()
	if err != nil {
		log.Fatalf("%s", err)
	}
	t.Cleanup(func() {
		Pgdb.resetForTests(t)
	})
	return
}
