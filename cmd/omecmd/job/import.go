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

	"github.com/openmarketplaceengine/openmarketplaceengine/app"
	"github.com/openmarketplaceengine/openmarketplaceengine/app/arg"
	"github.com/openmarketplaceengine/openmarketplaceengine/app/dir"
	"github.com/openmarketplaceengine/openmarketplaceengine/app/enc/geo"
	"github.com/openmarketplaceengine/openmarketplaceengine/cmd/omecmd/cfg"
	"github.com/openmarketplaceengine/openmarketplaceengine/internal/api/job/v1beta1"
	typeV1beta1 "github.com/openmarketplaceengine/openmarketplaceengine/internal/api/type/v1beta1"
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
	svc := v1beta1.NewJobServiceClient(con)
	var req v1beta1.ImportJobRequest
	var res *v1beta1.ImportJobResponse
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
	pickupLoc  typeV1beta1.Location
	dropoffLoc typeV1beta1.Location
	pickupDate timestamppb.Timestamp
	created    timestamppb.Timestamp
	updated    timestamppb.Timestamp
)

func updateRequest(req *v1beta1.ImportJobRequest, j *jobfile) {
	req.Reset()
	req.Job = new(v1beta1.JobInfo)
	updateTimestamp(&created, j.Created.Time)
	updateTimestamp(&updated, j.Updated.Time)
	updateTimestamp(&pickupDate, j.PickupDate.Time)
	updateLocation(&pickupLoc, j.PickupGeo)
	updateLocation(&dropoffLoc, j.DropoffGeo)
	r := req.Job
	r.Id = j.ID
	r.WorkerId = j.WorkerID
	r.Created = &created
	r.Updated = &updated
	r.State = j.State
	r.PickupDate = &pickupDate
	r.PickupAddr = j.PickupAddr
	r.PickupLoc = &pickupLoc
	r.DropoffAddr = j.DropoffAddr
	r.DropoffLoc = &dropoffLoc
	r.TripType = j.TripType
	r.Category = j.Category
}

// &Timestamp{Seconds: int64(t.Unix()), Nanos: int32(t.Nanosecond())}.
func updateTimestamp(dst *timestamppb.Timestamp, src time.Time) {
	dst.Reset()
	dst.Seconds = src.Unix()
	dst.Nanos = int32(src.Nanosecond())
}

func updateLocation(dst *typeV1beta1.Location, src geo.LocationWKB) {
	dst.Reset()
	dst.Latitude = src.Latitude
	dst.Longitude = src.Longitude
}
