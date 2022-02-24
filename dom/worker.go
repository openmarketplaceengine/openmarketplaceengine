// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package dom

type WorkerStatus int32

const (
	WorkerCreated   = iota // a worker registered in the system
	WorkerApproval         // worker is pending regulatory/system approval
	WorkerOffline          // worker is offline
	WorkerReady            // worker is ready for a job
	WorkerOnjob            // worker is on the job
	WorkerSuspended        // worker is suspended for some reason
	WorkerRetired          // worker is retired and will no longer be available
)

// Worker represents information about a driver.
type Worker struct {
	ID        UUID
	Status    WorkerStatus // worker's status
	Rank      int32        // worker's rank by customers
	Jobs      int          // total number of jobs completed
	FirstName string       // first name
	LastName  string       // last name
	Vehicle   UUID         // vehicle id foreign key
	Address   UUID         // address id foreign key
	Email     string
	Phone     string
	Comment   string
	Created   Time
	Approved  Time
	Suspended Time
	Retired   Time
	FirstJob  Time
	LastJob   Time
}

// WorkerVehicle represents relationships among workers and vehicles.
//
// One worker can operate one or more vehicles and one vehicle can be shared among
// several workers.
type WorkerVehicle struct {
	Worker  UUID
	Vehicle UUID
}
