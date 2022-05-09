// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package app

import "sync"

const (
	Domain    = "openmarketplaceengine"
	SrvAppUID = "omesrv"
	CmdAppUID = "omecmd"
)

type appset struct {
	sync.Mutex
	apps []*Application
}

var apps appset

//-----------------------------------------------------------------------------

func Server() *Application {
	return Get(SrvAppUID)
}

func Client() *Application {
	return Get(CmdAppUID)
}

//-----------------------------------------------------------------------------

func Get(name string) *Application {
	return apps.get(name)
}

//-----------------------------------------------------------------------------

func (m *appset) get(name string) *Application {
	m.Lock()
	defer m.Unlock()
	for i := 0; i < len(m.apps); i++ {
		if app := m.apps[i]; app.Name == name {
			return app
		}
	}
	app := NewApp(name)
	m.apps = append(m.apps, app)
	return app
}
