package timestampuid

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetTimestampUID(t *testing.T) {
	unique := make(map[string]bool)
	count := 10
	for i := 0; i < count; i++ {
		unique[NewTimestampUid()] = true
		time.Sleep(1 * time.Millisecond)
	}

	assert.Len(t, unique, count)
}
