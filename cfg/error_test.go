package cfg

import (
	"testing"
)

func TestTracePanic(t *testing.T) {
	defer func() {
		if rec := recover(); rec != nil {
			t.Logf("\npanic: %v\n%s", rec, TracePanic())
		}
	}()
	raisePanic()
}

func raisePanic() {
	panic("This is a test panic message")
}
