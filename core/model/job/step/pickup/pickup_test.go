package pickup

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPickUp(t *testing.T) {
	sync := make(chan string)
	onExit := func() {
		sync <- "done"
	}
	pickup := NewPickup(onExit)

	go pickup.Complete()
	require.Equal(t, "done", <-sync)
}
