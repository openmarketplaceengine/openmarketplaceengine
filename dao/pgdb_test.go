package dao

import "testing"

func TestWillTest(t *testing.T) {
	WillTest(t, "test")
}

func TestWillTestReentrant(t *testing.T) {
	WillTest(t, "test")
	WillTest(t, "test")
	WillTest(t, "test")
}
