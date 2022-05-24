// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package job

import (
	"time"

	"github.com/openmarketplaceengine/openmarketplaceengine/dom"

	"github.com/openmarketplaceengine/openmarketplaceengine/dao"
)

const table = "job"

type Job struct {
	ID          dom.SUID  `db:"id"`
	WorkerID    dom.SUID  `db:"worker_id"`
	Created     time.Time `db:"created"`
	Updated     time.Time `db:"updated"`
	State       string    `db:"state"`
	PickupDate  time.Time `db:"pickup_date"`
	PickupAddr  string    `db:"pickup_addr"`
	PickupLat   float64   `db:"pickup_lat"`
	PickupLon   float64   `db:"pickup_lon"`
	DropoffAddr string    `db:"dropoff_addr"`
	DropoffLat  float64   `db:"dropoff_lat"`
	DropoffLon  float64   `db:"dropoff_lon"`
	TripType    string    `db:"trip_type"`
	Category    string    `db:"category"`
}

func (j *Job) Upsert(ctx dom.Context) (dao.Result, dao.UpsertStatus, error) {
	return dao.Upsert(ctx, j.insert, j.update)
}

func (j *Job) insert() dao.Executable {
	sql := dao.Insert(table).
		Set("id", j.ID).
		Set("worker_id", j.WorkerID)
	j.setSQL(sql)
	return sql
}

func (j *Job) update() dao.Executable {
	sql := dao.Update(table)
	j.setSQL(sql)
	sql.Where("id = ?", j.ID)
	return sql
}

func (j *Job) setSQL(sql *dao.SQL) {
	sql.
		Set("created", j.Created).
		Set("updated", j.Updated).
		Set("state", j.State).
		Set("pickup_date", j.PickupDate).
		Set("pickup_addr", j.PickupAddr).
		Set("pickup_lat", j.PickupLat).
		Set("pickup_lon", j.PickupLon).
		Set("dropoff_addr", j.DropoffAddr).
		Set("dropoff_lat", j.DropoffLat).
		Set("dropoff_lon", j.DropoffLon).
		Set("trip_type", j.TripType).
		Set("category", j.Category)
}
