// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package job

import (
	"fmt"
	"testing"
	"time"

	"github.com/driverscooperative/geosrv/cfg"
	"github.com/driverscooperative/geosrv/dao"
	"github.com/driverscooperative/geosrv/dom"
	"github.com/stretchr/testify/require"
)

func TestStateFromString(t *testing.T) {
	for i := range stateString {
		_, found := StateFromString(stateString[i])
		require.True(t, found)
	}
	_, found := StateFromString("dummy")
	require.False(t, found)
}

func TestState_String(t *testing.T) {
	for k, v := range stringState {
		require.Equal(t, k, v.String())
	}
}

//-----------------------------------------------------------------------------
// SetState
//-----------------------------------------------------------------------------

func TestSetState(t *testing.T) {
	ctx := testJobState(t)
	for i := 1; i < len(stateString); i++ {
		testSetState(t, ctx, fmt.Sprint(i), State(i), true)
	}
	testSetState(t, ctx, dao.NewXid(), StateAvailable, false)
}

func testSetState(t *testing.T, ctx dom.Context, jobID string, state State, mustSet bool) {
	set, err := SetState(ctx, jobID, state)
	require.NoError(t, err)
	if mustSet {
		require.True(t, set)
	} else {
		require.False(t, set)
	}
}

//-----------------------------------------------------------------------------
// GetState
//-----------------------------------------------------------------------------

func TestGetState(t *testing.T) {
	ctx := testJobState(t)
	for i := 1; i < len(stateString); i++ {
		testGetState(t, ctx, fmt.Sprint(i), State(i), true)
	}
	testGetState(t, ctx, dao.NewXid(), StateUnspecified, false)
}

func testGetState(t *testing.T, ctx dom.Context, jobID string, state State, mustGet bool) {
	s, found, err := GetState(ctx, jobID)
	require.NoError(t, err)
	if !mustGet {
		require.False(t, found)
		return
	}
	require.True(t, found)
	require.Equal(t, state, s)
}

//-----------------------------------------------------------------------------
// Helpers
//-----------------------------------------------------------------------------

func testJobState(t testing.TB) dom.Context {
	dom.WillTest(t, "test", false)
	ctx := cfg.Context()
	for i := 1; i < len(stateString); i++ {
		j := genJob()
		j.ID = fmt.Sprintf("%d", i)
		j.WorkerID = j.ID
		j.State = State(i).String()
		_, _, err := j.Upsert(ctx)
		require.NoError(t, err)
	}
	return ctx
}

func genJob() *Job {
	return &Job{
		ID:          dao.NewXid(),
		WorkerID:    dao.NewXid(),
		Created:     time.Now(),
		Updated:     time.Now(),
		State:       StateAvailable.String(),
		PickupDate:  time.Now(),
		PickupAddr:  "QU: 49th Pl & Maspeth Ave",
		PickupLat:   40.72193,
		PickupLon:   -73.919973,
		DropoffAddr: "BK: Rutland Rd & E 91st St",
		DropoffLat:  40.6614552,
		DropoffLon:  -73.9282859,
		TripType:    "Passenger",
		Category:    "AAR",
	}
}
