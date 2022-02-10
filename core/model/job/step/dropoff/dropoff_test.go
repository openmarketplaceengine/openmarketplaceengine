package dropoff

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDropOff(t *testing.T) {
	sync := make(chan string)
	onExit := func() {
		sync <- "done"
	}
	dropOff := NewDropOff(onExit)

	go dropOff.Complete()
	require.Equal(t, "done", <-sync)
}
