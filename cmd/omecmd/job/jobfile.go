// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package job

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
	"unsafe"

	"github.com/openmarketplaceengine/openmarketplaceengine/app/enc/geo"
	"gopkg.in/yaml.v2"
)

type jobfile struct {
	ID            string          `json:"id" yaml:"id"`
	WorkerID      string          `json:"worker_id" yaml:"worker_id"`
	Created       stamp           `json:"created_at" yaml:"created_at"`
	Updated       stamp           `json:"updated_at" yaml:"updated_at"`
	State         string          `json:"state" yaml:"state"`
	PickupDate    datetime        `json:"pickup_datetime" yaml:"pickup_datetime"`
	PickupAddr    string          `json:"pickup_location_readable" yaml:"pickup_location_readable"`
	PickupGeo     geo.LocationWKB `json:"pickup_location_geo" yaml:"pickup_location_geo"`
	DropoffAddr   string          `json:"dropoff_location_readable" yaml:"dropoff_location_readable"`
	DropoffGeo    geo.LocationWKB `json:"dropoff_location_geo" yaml:"dropoff_location_geo"`
	CustomerName  string          `json:"customer_name" yaml:"customer_name"`
	CustomerPhone string          `json:"customer_phone" yaml:"customer_phone"`
	TripType      string          `json:"trip_type" yaml:"trip_type"`
	Category      string          `json:"category" yaml:"category"`
}

//-----------------------------------------------------------------------------

func (j *jobfile) reset() {
	*j = jobfile{}
}

//-----------------------------------------------------------------------------

func (j *jobfile) dumpJSON() {
	buf, err := json.MarshalIndent(j, "", "  ")
	if err != nil {
		println(err.Error())
		return
	}
	buf = append(buf, '\n')
	_, _ = os.Stdout.Write(buf)
}

//-----------------------------------------------------------------------------

func (j *jobfile) dumpYAML() {
	buf, err := yaml.Marshal(j)
	if err != nil {
		println(err.Error())
		return
	}
	// buf = append(buf, '\n')
	_, _ = os.Stdout.Write(buf)
}

//-----------------------------------------------------------------------------

func parseTime(dest *time.Time, data []byte, layout ...string) (err error) {
	value := *(*string)(unsafe.Pointer(&data))
	for i := 0; i < len(layout); i++ {
		*dest, err = time.Parse(layout[i], value)
		if err == nil {
			return
		}
	}
	return fmt.Errorf("invalid time format: %q", value)
}

//-----------------------------------------------------------------------------
// Time Types
//-----------------------------------------------------------------------------

type stamp struct {
	time.Time
}

func (t *stamp) UnmarshalText(data []byte) error {
	return parseTime(&t.Time, data, time.RFC3339Nano)
}

type datetime struct {
	time.Time
}

func (t *datetime) UnmarshalText(data []byte) error {
	return parseTime(&t.Time, data, time.RFC3339)
}
