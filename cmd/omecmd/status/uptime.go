// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package status

import (
	"context"

	"github.com/openmarketplaceengine/openmarketplaceengine/app"
	"github.com/openmarketplaceengine/openmarketplaceengine/cmd/omecmd/cfg"
)

func init() {
	args := &app.Client().Args
	args.Void("uptime", "Server uptime", uptime)
}

func uptime(_ context.Context) error {
	cfg.Debugf("Requesting server uptime...")
	return nil
}
