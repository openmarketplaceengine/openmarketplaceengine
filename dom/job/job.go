// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package job

import (
	"fmt"
	"time"

	"github.com/openmarketplaceengine/openmarketplaceengine/dao"
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
	Distance Distance
}

type Distance struct {
	FromLat float64
	FromLon float64
	Unit    RangeUnit
	Range   float64
}

type RangeUnit string

const M RangeUnit = "m"
const Km RangeUnit = "km"
const Mile RangeUnit = "mile"

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

//QueryByPickupDistance is used to select jobs nearest to specified location within specified range limit.
//NB! This is MVP version, not suitable for production use.
func QueryByPickupDistance(ctx dao.Context, fromLat float64, fromLon float64, state string, rangeLimit float64, rangeUnit RangeUnit, limit int) ([]*AvailableJob, error) {
	stmt := jobsInRangeStmt(fromLat, fromLon, state, rangeLimit, rangeUnit, limit)
	s := dao.NewSQL(stmt)
	jobs := make([]*AvailableJob, 0, limit)
	err := s.QueryRows(ctx, func(rows *dao.Rows) error {
		for rows.Next() {
			var j AvailableJob
			var r float64
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
				&r,
			); err != nil {
				return err
			}
			j.Distance = Distance{
				FromLat: fromLat,
				FromLon: fromLon,
				Unit:    rangeUnit,
				Range:   r,
			}
			jobs = append(jobs, &j)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return jobs, nil
}

//jobsInRangeStmt returns plain sql of the
//query taken from https://gis.stackexchange.com/questions/31628/find-features-within-given-coordinates-and-distance-using-mysql
//that follows Math from https://www.movable-type.co.uk/scripts/latlong.html
//
//select id,
//		(
//		3959 * acos(
//			   cos(radians(78.3232))
//			   * cos(radians(pickup_lat))
//			   * cos(radians(pickup_lon) - radians(65.3234))
//		   + sin(radians(78.3232)) * sin(radians(pickup_lat))
//		)
//		) as distance
//from job
//order by distance
//having distance < 30
//limit 20;
//To search by kilometers instead of miles, replace 3959 with 6371.
func jobsInRangeStmt(latitude float64, longitude float64, state string, rangeLimit float64, rangeUnit RangeUnit, limit int) string {
	var u int
	if rangeUnit == Km {
		u = 6371
	} else if rangeUnit == M {
		u = 6371 * 1000
	} else {
		u = 3959
	}

	return fmt.Sprintf(`
	select t.*
	from (select *,
				 %v * acos(
				 	cos(radians(%v))
					 * cos(radians(pickup_lat))
					 * cos(radians(pickup_lon) - radians(%v))
					 + sin(radians(%v)) * sin(radians(pickup_lat))
				) as range
		  from job) as t
	where t.state = '%s'
	  and t.range < %v
	order by t.range
	limit %v
	`, u, latitude, longitude, latitude, state, rangeLimit, limit)
}

func DeleteAll(ctx dao.Context) error {
	del := dao.Delete(table)
	return dao.ExecTX(ctx, del)
}
