// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package dom

import (
	"fmt"

	"github.com/openmarketplaceengine/openmarketplaceengine/dao"
)

type WorkerStatus int32

const (
	WorkerOffline  = iota // worker is offline
	WorkerReady           // worker is ready for a job
	WorkerOnjob           // worker is on the job
	WorkerPaused          // worker is on hold
	WorkerDisabled        // worker is disabled
)

const (
	workerTable        = "worker"
	workerVehicleTable = "worker_vehicle"
)

// Worker represents information about a driver.
type Worker struct {
	ID        SUID         `db:"id"`
	Status    WorkerStatus `db:"status"`     // worker's status
	Rating    int32        `db:"rating"`     // worker's rating by customers
	Jobs      int          `db:"jobs"`       // total number of jobs completed
	FirstName string       `db:"first_name"` // first name
	LastName  string       `db:"last_name"`  // last name
	Vehicle   SUID         `db:"vehicle"`    // vehicle id foreign key
	Created   Time         `db:"created"`    // worker creation time
	Updated   Time         `db:"updated"`    // worker last modified time
}

// WorkerVehicle represents relationships among workers and vehicles.
//
// One worker can operate one or more vehicles and one vehicle can be shared among
// several workers.
type WorkerVehicle struct {
	Worker  SUID `db:"worker"`
	Vehicle SUID `db:"vehicle"`
}

//-----------------------------------------------------------------------------

// Persist saves Worker to the database.
func (w *Worker) Persist(ctx Context) error {
	return dao.ExecTX(ctx, w.Insert())
}

//-----------------------------------------------------------------------------

func (w *Worker) Insert() dao.Executable {
	sql := dao.Insert(workerTable)
	sql.Set("id", w.ID)
	sql.Set("status", w.Status)
	sql.Set("rating", w.Rating)
	sql.Set("jobs", w.Jobs)
	sql.Set("first_name", w.FirstName)
	sql.Set("last_name", w.LastName)
	sql.Set("vehicle", w.Vehicle)
	sql.Set("created", w.Created)
	sql.Set("updated", w.Updated)
	return sql
}

//-----------------------------------------------------------------------------
// Worker <-> Vehicle
//-----------------------------------------------------------------------------

// Persist saves WorkerVehicle to the database.
func (w *WorkerVehicle) Persist(ctx Context) error {
	return dao.ExecTX(ctx, w.Insert())
}

func (w *WorkerVehicle) Insert() dao.Executable {
	sql := dao.Insert(workerVehicleTable)
	sql.Set("worker", w.Worker)
	sql.Set("vehicle", w.Vehicle)
	return sql
}

//-----------------------------------------------------------------------------

// String representation of WorkerStatus.
func (s WorkerStatus) String() string {
	switch s {
	case WorkerOffline:
		return "offline"
	case WorkerReady:
		return "ready"
	case WorkerOnjob:
		return "onjob"
	case WorkerPaused:
		return "paused"
	case WorkerDisabled:
		return "disabled"
	}
	return fmt.Sprintf("WorkerStatus<%d>", s)
}
