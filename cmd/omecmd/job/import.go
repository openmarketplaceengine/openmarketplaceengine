// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package job

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	jobimpv1 "github.com/openmarketplaceengine/openmarketplaceengine/api/gen/jobimp/v1"
	"github.com/openmarketplaceengine/openmarketplaceengine/app"
	"github.com/openmarketplaceengine/openmarketplaceengine/app/arg"
	"github.com/openmarketplaceengine/openmarketplaceengine/app/dir"
	"github.com/openmarketplaceengine/openmarketplaceengine/app/enc/geo"
	"github.com/openmarketplaceengine/openmarketplaceengine/cmd/omecmd/cfg"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gopkg.in/yaml.v2"
)

const (
	fyaml = ".yaml"
	fjson = ".json"
)

//-----------------------------------------------------------------------------

func init() {
	const flags = arg.FileMustExist | arg.FileSkipInvalid | arg.PathPrintError
	v := arg.FileValidator(flags, fyaml, fjson)
	app.Client().Args.Files("jobimp", "Import job description from YAML or JSON `file(s)`", flags, Jobimp).Validator(v)
}

//-----------------------------------------------------------------------------

func Jobimp(ctx context.Context, files []string) error {
	con, err := cfg.Dial(ctx)
	if err != nil {
		return err
	}
	defer cfg.SafeClose(con)
	svc := jobimpv1.NewJobimpServiceClient(con)
	var req jobimpv1.JobimpRequest
	var res *jobimpv1.JobimpResponse
	cfg.Debugf("importing %d job file(s)", len(files))
	var job jobfile
	var dec dir.DecodeFunc
	for i := range files {
		if app.Done(ctx) {
			return ctx.Err()
		}
		file := files[i]
		cfg.Debugf("parsing: %q", file)
		job.reset()
		switch {
		case strings.HasSuffix(file, fyaml):
			dec = yaml.Unmarshal
		case strings.HasSuffix(file, fjson):
			dec = json.Unmarshal
		default:
			return fmt.Errorf("invalid file extension: %q", file)
		}
		err := dir.DecodeFile(file, &job, dec)
		if err != nil {
			return fmt.Errorf("error parsing: %q: %v", file, err)
		}
		updateRequest(&req, &job)
		cfg.Debugf("sending: %q", file)
		res, err = svc.ImportJob(ctx, &req)
		if err != nil {
			return err
		}
		cfg.Debugf("result : %s", res.Action.String())
	}
	cfg.Debugf("done")
	return nil
}

//-----------------------------------------------------------------------------
// Request Fulfillment
//-----------------------------------------------------------------------------

var (
	pickup_loc  jobimpv1.Location
	dropoff_loc jobimpv1.Location
	pickup_date timestamppb.Timestamp
	created     timestamppb.Timestamp
	updated     timestamppb.Timestamp
)

func updateRequest(r *jobimpv1.JobimpRequest, j *jobfile) {
	r.Reset()
	updateTimestamp(&created, j.Created.Time)
	updateTimestamp(&updated, j.Updated.Time)
	updateTimestamp(&pickup_date, j.PickupDate.Time)
	updateLocation(&pickup_loc, j.PickupGeo)
	updateLocation(&dropoff_loc, j.DropoffGeo)
	r.Id = j.ID
	r.WorkerId = j.WorkerID
	r.Created = &created
	r.Updated = &updated
	r.State = j.State
	r.PickupDate = &pickup_date
	r.PickupAddr = j.PickupAddr
	r.PickupLoc = &pickup_loc
	r.DropoffAddr = j.DropoffAddr
	r.DropoffLoc = &dropoff_loc
	r.TripType = j.TripType
	r.Category = j.Category
}

// &Timestamp{Seconds: int64(t.Unix()), Nanos: int32(t.Nanosecond())}
func updateTimestamp(dst *timestamppb.Timestamp, src time.Time) {
	dst.Reset()
	dst.Seconds = src.Unix()
	dst.Nanos = int32(src.Nanosecond())
}

func updateLocation(dst *jobimpv1.Location, src geo.LocationWKB) {
	dst.Reset()
	dst.Lat = src.Lat
	dst.Lon = src.Lon
}
