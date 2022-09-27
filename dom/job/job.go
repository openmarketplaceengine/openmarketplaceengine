// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package job

import (
	"fmt"
	"time"

	"github.com/driverscooperative/geosrv/dao"
)

const table = "job"

type Job struct {
	ID          dao.SUID  `db:"id"`
	WorkerID    dao.SUID  `db:"worker_id"`
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

type AvailableJob struct {
	Job
	DistanceMeters int32
}

func (j *Job) Upsert(ctx dao.Context) (dao.Result, dao.UpsertStatus, error) {
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

func QueryOne(ctx dao.Context, jobID dao.SUID) (job *Job, has bool, err error) {
	job = new(Job)
	has, err = dao.From(table).
		Bind(job).
		Where("id = ?", jobID).
		QueryOne(ctx)
	if !has {
		job = nil
	}
	return
}

// QueryByPickupDistance is used to select jobs nearest to specified lat/lon within precision in meters.
// lat/lon - point to count distance from
// radiusMeters - radius in meters
// limit - rows limit
// returns AvailableJob array.
func QueryByPickupDistance(ctx dao.Context, lon float64, lat float64, state string, radiusMeters int32, limit int32) ([]*AvailableJob, error) {
	distance := fmt.Sprintf("st_distance(st_point(%v, %v, 4326)::geography, st_point(pickup_lon, pickup_lat, 4326)::geography) as distance", lon, lat)
	within := fmt.Sprintf("st_dwithin(st_point(%v, %v, 4326)::geography, st_point(pickup_lon, pickup_lat, 4326)::geography, %v)", lon, lat, radiusMeters)
	sql := dao.From(table).
		Select("*").
		Select(distance).
		Where(within).
		OrderBy("distance").
		Limit(limit)
	if state != "" {
		sql.Where("state = ?", state)
	}
	ary := make([]*AvailableJob, 0, limit)
	err := sql.QueryRows(ctx, func(rows *dao.Rows) error {
		for rows.Next() {
			var j AvailableJob
			var d float64
			if err := rows.Scan(
				&j.ID,
				&j.WorkerID,
				&j.Created,
				&j.Updated,
				&j.State,
				&j.PickupDate,
				&j.PickupAddr,
				&j.PickupLat,
				&j.PickupLon,
				&j.DropoffAddr,
				&j.DropoffLat,
				&j.DropoffLon,
				&j.TripType,
				&j.Category,
				&d,
			); err != nil {
				return err
			}
			j.DistanceMeters = int32(d)
			ary = append(ary, &j)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return ary, nil
}

func DeleteAll(ctx dao.Context) error {
	del := dao.Delete(table)
	return dao.ExecTX(ctx, del)
}
