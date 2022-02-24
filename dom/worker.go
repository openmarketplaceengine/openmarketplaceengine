// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package dom

import "fmt"

type WorkerStatus int32

const (
	WorkerOffline  = iota // worker is offline
	WorkerReady           // worker is ready for a job
	WorkerOnjob           // worker is on the job
	WorkerPaused          // worker is on hold
	WorkerDisabled        // worker is disabled
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
