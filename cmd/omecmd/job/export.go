// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package job

import (
	"context"

	"github.com/openmarketplaceengine/openmarketplaceengine/app"
	"github.com/openmarketplaceengine/openmarketplaceengine/app/dir"
	"github.com/openmarketplaceengine/openmarketplaceengine/cmd/omecmd/cfg"
	rpc "github.com/openmarketplaceengine/openmarketplaceengine/internal/idl/api/job/v1beta1"
	"gopkg.in/yaml.v2"
)

func init() {
	app.Client().Args.Rest("jobexp", "Export job description with `id(s)` to YAML", JobExp)
}

func JobExp(ctx context.Context, ids []string) error { //nolint
	dst, derr := cfg.DstDir()
	if derr != nil {
		return derr
	}
	con, cerr := cfg.Dial(ctx)
	if cerr != nil {
		return cerr
	}
	defer cfg.SafeClose(con)
	svc := rpc.NewJobServiceClient(con)
	req := rpc.ExportJobRequest{Ids: ids}
	res, err := svc.ExportJob(ctx, &req)
	if err != nil {
		return err
	}
	var file jobfile
	for _, item := range res.Jobs {
		if item.Job == nil {
			cfg.Errorf("job not found for id: %q", item.Id)
			continue
		}
		path := dir.FastJoin(dst, item.Id) + fyaml
		cfg.Debugf("writing job: %q", item.Id)
		file.FromInfo(item.Job)
		err = dir.EncodeFile(path, &file, 0, yaml.Marshal)
		if err != nil {
			return err
		}
	}
	return nil
}

func (j *jobfile) FromInfo(inf *rpc.JobInfo) {
	j.ID = inf.Id
	j.WorkerID = inf.WorkerId
	j.Created.Time = inf.Created.AsTime()
	j.Updated.Time = inf.Updated.AsTime()
	j.State = inf.State
	j.PickupDate.Time = inf.PickupDate.AsTime()
	j.PickupAddr = inf.PickupAddr
	j.PickupGeo.Latitude = inf.PickupLoc.GetLatitude()
	j.PickupGeo.Longitude = inf.PickupLoc.GetLongitude()
	j.DropoffAddr = inf.DropoffAddr
	j.DropoffGeo.Latitude = inf.DropoffLoc.GetLatitude()
	j.DropoffGeo.Longitude = inf.DropoffLoc.GetLongitude()
	j.TripType = inf.TripType
	j.Category = inf.Category
}