// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package status

import (
	"context"
	"fmt"
	"time"

	"github.com/openmarketplaceengine/openmarketplaceengine/app"
	"github.com/openmarketplaceengine/openmarketplaceengine/cmd/omecmd/cfg"
	"github.com/openmarketplaceengine/openmarketplaceengine/internal/idl/api/status/v1beta1"
)

func init() {
	args := &app.Client().Args
	args.Void("uptime", "Server uptime", uptime)
}

func uptime(ctx context.Context) error {
	con, err := cfg.Dial(ctx)
	if err != nil {
		return err
	}
	defer cfg.SafeClose(con)
	svc := v1beta1.NewUptimeServiceClient(con)
	res, err2 := svc.GetUptime(ctx, &v1beta1.GetUptimeRequest{})
	if err2 != nil {
		return err2
	}
	fmt.Println("Started:", res.Started.AsTime().Local())
	fmt.Println(" Uptime:", time.Duration(res.Uptime))
	return nil
}
