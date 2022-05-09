// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package app

import (
	"encoding/json"
	"os"

	"github.com/openmarketplaceengine/openmarketplaceengine/app/arg"
)

type Application struct {
	Name string
	Dirs Dir
	Args arg.ArgSet
}

//-----------------------------------------------------------------------------

func NewApp(name string) *Application {
	a := &Application{Name: name}
	a.Args.Init()
	return a
}

//-----------------------------------------------------------------------------

func (a *Application) Init(domain string, cfgFileName string) error {
	return a.Dirs.Init(domain, a.Name, cfgFileName)
}

func (a *Application) MustInit(domain string, cfgFileName string) {
	if err := a.Init(domain, cfgFileName); err != nil {
		Exit(err)
	}
}

//-----------------------------------------------------------------------------

func (a *Application) DumpJSON() bool {
	buf, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		println("json.MarshalIndent failed:", err.Error())
		return false
	}
	buf = append(buf, '\n')
	_, err = os.Stdout.Write(buf)
	return err == nil
}
