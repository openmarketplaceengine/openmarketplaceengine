// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package worker

import (
	"fmt"
	"time"

	"github.com/openmarketplaceengine/openmarketplaceengine/dom"

	"github.com/openmarketplaceengine/openmarketplaceengine/dao"
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

// Persist saves Worker to the database.
func (w *Worker) Persist(ctx dom.Context) error {
	return dao.ExecTX(ctx, w.Insert())
}

//-----------------------------------------------------------------------------

func (w *Worker) Insert() dao.Executable {
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

func GetWorker(ctx dom.Context, workerID dom.SUID) (wrk *Worker, has bool, err error) {
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

//-----------------------------------------------------------------------------

func GetWorkerStatus(ctx dom.Context, workerID dom.SUID) (Status, bool, error) {
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

func SetWorkerStatus(ctx dom.Context, workerID dom.SUID, status Status) error {
	sql := dao.Update(workerTable).
		Set("status", int32(status)).
		Set("updated", time.Now()).
		Where("id = ?", workerID)
	return dao.ExecTX(ctx, sql)
}

//-----------------------------------------------------------------------------
// Worker <-> Vehicle
//-----------------------------------------------------------------------------

// Persist saves WorkerVehicle to the database.
func (w *WorkerVehicle) Persist(ctx dom.Context) error {
	return dao.ExecTX(ctx, w.Insert())
}

func (w *WorkerVehicle) Insert() dao.Executable {
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
