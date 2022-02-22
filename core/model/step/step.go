package step

import "github.com/openmarketplaceengine/openmarketplaceengine/core/model/job"

type ID string

// Action defines primitive action constituting a job.Job.
type Action int

const (
	GoToLocation Action = iota
	Pickup
	DropOff
	CollectCache
	CollectVoucher
	CallPhone
)

// Step is a part of Job execution
// JobID refers to job.Job step belongs to.
type Step struct {
	ID     ID
	JobID  job.ID
	Action Action
}
