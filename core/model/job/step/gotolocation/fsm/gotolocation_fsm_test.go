package fsm

import (
	_fsm "github.com/cocoonspace/fsm"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGoToLocationFSM(t *testing.T) {
	sync := make(chan string)

	hooks := make(map[_fsm.State]func())
	hooks[MovingState] = func() {
		sync <- "moving"
	}
	hooks[NearState] = func() {
		sync <- "near"
	}
	hooks[ArrivedState] = func() {
		sync <- "arrived"
	}
	hooks[CanceledState] = func() {
		sync <- "cancelled"
	}

	fsm := NewGoToLocationFSM(hooks)

	go func() {
		fsm.f.Event(Near)
		fsm.f.Event(Arrived)
		fsm.f.Event(Canceled)
	}()
	require.Equal(t, "near", <-sync)
	require.Equal(t, "arrived", <-sync)
	require.Equal(t, "cancelled", <-sync)
	require.Equal(t, CanceledState, fsm.f.Current())
}
