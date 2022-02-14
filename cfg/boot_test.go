// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package cfg

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestBootList_Boot(t *testing.T) {
	b := testBoot(t)
	addBoot(b, "DBMS", false, false)
	addBoot(b, "GRPC", false, false)
	addBoot(b, "HTTP", false, false)
	require.NoError(t, b.Boot())
	require.NoError(t, b.Stop())
}

func TestBootList_Error(t *testing.T) {
	b := testBoot(t)
	addBoot(b, "DBMS", false, false)
	addBoot(b, "GRPC", false, false)
	addBoot(b, "HTTP", true, false)
	require.Error(t, b.Boot())
}

//-----------------------------------------------------------------------------
// Helpers
//-----------------------------------------------------------------------------

func addBoot(b *BootList, name string, failBoot bool, failStop bool) { //nolint
	b.Add(name, testBootFunc(failBoot), testBootFunc(failStop))
}

//-----------------------------------------------------------------------------

func testBootFunc(fail bool) BootFunc {
	return func() error {
		if fail {
			return fmt.Errorf("requested test error: %d", byte(rand.Uint32()))
		}
		return nil
	}
}

//-----------------------------------------------------------------------------

func testBoot(t testing.TB) *BootList {
	b := new(BootList)
	if testing.Verbose() {
		b.Debugf = testPrintFunc(t, "DEBUG")
		b.Errorf = testPrintFunc(t, "ERROR")
	}
	return b
}

//-----------------------------------------------------------------------------

func testPrintFunc(t testing.TB, level string) PrintFunc {
	const timeFormat = "15:04:05.000"
	return func(format string, args ...interface{}) {
		t.Helper()
		t.Logf("%s [%s] %s", time.Now().Format(timeFormat), level, fmt.Sprintf(format, args...))
	}
}
