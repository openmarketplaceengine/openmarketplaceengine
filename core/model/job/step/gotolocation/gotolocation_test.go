package gotolocation

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGoToLocation(t *testing.T) {
	sync := make(chan string)
	onExit := func() {
		sync <- "done"
	}
	goToLocation := NewGoToLocation(onExit)

	go func() {
		goToLocation.Traveling()
		goToLocation.NearDestination()
		goToLocation.AtDestination()
		goToLocation.Complete()
	}()
	require.Equal(t, "done", <-sync)
}
