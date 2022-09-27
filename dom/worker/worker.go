// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package worker

import (
	"fmt"
	"time"

	"github.com/driverscooperative/geosrv/dom"

	"github.com/driverscooperative/geosrv/dao"
)

type Status int32

const (
	Offline  = iota // worker is offline
	Ready           // worker is ready for a job
	OnJob           // worker is on the job
	Paused          // worker is on hold
	Disabled        // worker is disabled
)

const (
	workerTable        = "worker"
	workerVehicleTable = "worker_vehicle"
)

// Worker represents information about a driver.
type Worker struct {
	ID        dom.SUID `db:"id"`
	Status    Status   `db:"status"`     // worker's status
	Rating    int32    `db:"rating"`     // worker's rating by customers
	Jobs      int      `db:"jobs"`       // total number of jobs completed
	FirstName string   `db:"first_name"` // first name
	LastName  string   `db:"last_name"`  // last name
	Vehicle   dom.SUID `db:"vehicle"`    // vehicle id foreign key
	Created   dom.Time `db:"created"`    // worker creation time
	Updated   dom.Time `db:"updated"`    // worker last modified time
}

// WorkerVehicle represents relationships among workers and vehicles.
//
// One worker can operate one or more vehicles and one vehicle can be shared among
// several workers.
type WorkerVehicle struct { //nolint:revive
	Worker  dom.SUID `db:"worker"`
	Vehicle dom.SUID `db:"vehicle"`
}

//-----------------------------------------------------------------------------

// Insert inserts Worker to the database.
func (w *Worker) Insert(ctx dom.Context) error {
	return dao.ExecTX(ctx, w.insert())
}

//-----------------------------------------------------------------------------

func (w *Worker) insert() dao.Executable {
	return dao.Insert(workerTable).
		Set("id", w.ID).
		Set("status", w.Status).
		Set("rating", w.Rating).
		Set("jobs", w.Jobs).
		Set("first_name", w.FirstName).
		Set("last_name", w.LastName).
		Set("vehicle", w.Vehicle).
		Set("created", w.Created).
		Set("updated", w.Updated)
}

//-----------------------------------------------------------------------------
// Getters
//-----------------------------------------------------------------------------

func QueryOne(ctx dom.Context, workerID dom.SUID) (wrk *Worker, has bool, err error) {
	wrk = new(Worker)
	has, err = dao.From(workerTable).
		Bind(wrk).
		Where("id = ?", workerID).
		QueryOne(ctx)
	if !has {
		wrk = nil
	}
	return
}

func QueryAll(ctx dom.Context, status *Status, limit int, offset int) ([]*Worker, error) {
	query := dao.From(workerTable).
		Select("*").
		Limit(limit).
		Offset(offset)

	if status != nil {
		query.Where("status = ?", status)
	}

	ary := make([]*Worker, 0, limit)
	err := query.QueryRows(ctx, func(rows *dao.Rows) error {
		for rows.Next() {
			var w Worker
			if err := rows.Scan(
				&w.ID,
				&w.Status,
				&w.Rating,
				&w.Jobs,
				&w.FirstName,
				&w.LastName,
				&w.Vehicle,
				&w.Created,
				&w.Updated,
			); err != nil {
				return err
			}
			ary = append(ary, &w)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return ary, nil
}

//-----------------------------------------------------------------------------

func QueryWorkerStatus(ctx dom.Context, workerID dom.SUID) (Status, bool, error) {
	var status int32
	has, err := dao.From(workerTable).
		Select("status").To(&status).
		Where("id = ?", workerID).
		QueryOne(ctx)
	return Status(status), has, err
}

//-----------------------------------------------------------------------------
// Setters
//-----------------------------------------------------------------------------

func UpdateWorkerStatus(ctx dom.Context, workerID dom.SUID, status Status) error {
	sql := dao.Update(workerTable).
		Set("status", int32(status)).
		Set("updated", time.Now()).
		Where("id = ?", workerID)
	return dao.ExecTX(ctx, sql)
}

//-----------------------------------------------------------------------------
// Worker <-> Vehicle
//-----------------------------------------------------------------------------

// Insert inserts WorkerVehicle to the database.
func (w *WorkerVehicle) Insert(ctx dom.Context) error {
	return dao.ExecTX(ctx, w.insert())
}

func (w *WorkerVehicle) insert() dao.Executable {
	return dao.Insert(workerVehicleTable).
		Set("worker", w.Worker).
		Set("vehicle", w.Vehicle)
}

//-----------------------------------------------------------------------------

// String representation of WorkerStatus.
func (s Status) String() string {
	switch s {
	case Offline:
		return "offline"
	case Ready:
		return "ready"
	case OnJob:
		return "onjob"
	case Paused:
		return "paused"
	case Disabled:
		return "disabled"
	}
	return fmt.Sprintf("WorkerStatus<%d>", s)
}
