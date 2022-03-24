// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package dom

import (
	"fmt"
	"time"

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

func GetWorker(ctx Context, workerID SUID) (*Worker, error) {
	var wrk Worker
	err := dao.From(workerTable).
		Bind(&wrk).
		Where("id = ?", workerID).
		QueryOne(ctx)
	if err != nil {
		return nil, err
	}
	return &wrk, nil
}

//-----------------------------------------------------------------------------

func GetWorkerStatus(ctx Context, workerID SUID) (WorkerStatus, error) {
	var status int32
	err := dao.From(workerTable).
		Select("status").To(&status).
		Where("id = ?", workerID).
		QueryOne(ctx)
	if err != nil {
		return 0, err
	}
	return WorkerStatus(status), nil
}

//-----------------------------------------------------------------------------
// Setters
//-----------------------------------------------------------------------------

func SetWorkerStatus(ctx Context, workerID SUID, status WorkerStatus) error {
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
func (w *WorkerVehicle) Persist(ctx Context) error {
	return dao.ExecTX(ctx, w.Insert())
}

func (w *WorkerVehicle) Insert() dao.Executable {
	return dao.Insert(workerVehicleTable).
		Set("worker", w.Worker).
		Set("vehicle", w.Vehicle)
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
