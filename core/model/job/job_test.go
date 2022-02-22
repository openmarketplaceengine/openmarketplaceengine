package job

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJob(t *testing.T) {
	t.Run("testCreateJob", func(t *testing.T) {
		testCreateJob(t)
	})
}

func testCreateJob(t *testing.T) {
	request := RideRequest{
		PickupLocation:     Location{},
		DropOffLocation:    Location{},
		SubjectID:          "",
		RequestedTime:      time.Time{},
		RequestedStartTime: time.Time{},
	}
	job := NewJob(request, time.Now())

	assert.NotNil(t, job)
}
