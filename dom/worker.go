// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package dom

type WorkerStatus int32

const (
	WorkerCreated  = iota // a worker registered in the system
	WorkerOffline         // worker is offline
	WorkerReady           // worker is ready for a job
	WorkerOnjob           // worker is on the job
	WorkerDisabled        // worker is disabled and will no longer be available
)

// Worker represents information about a driver.
type Worker struct {
	ID        UUID
	Status    WorkerStatus // worker's status
	Rating    int32        // worker's rating by customers
	Jobs      int          // total number of jobs completed
	FirstName string       // first name
	LastName  string       // last name
	Vehicle   UUID         // vehicle id foreign key
	Created   Time         // worker creation time
	Updated   Time         // worker last modified time
}

// WorkerVehicle represents relationships among workers and vehicles.
//
// One worker can operate one or more vehicles and one vehicle can be shared among
// several workers.
type WorkerVehicle struct {
	Worker  UUID
	Vehicle UUID
}
