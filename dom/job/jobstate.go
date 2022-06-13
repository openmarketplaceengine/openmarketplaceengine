// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package job

import "github.com/openmarketplaceengine/openmarketplaceengine/dao"

type State int

const (
	StateUnspecified      State = 0
	StateAvailable        State = 1
	StateGoToPickup       State = 2
	StateNearPickup       State = 3
	StateAwaitingPickup   State = 4
	StateOnTrip           State = 5
	StateNearDestination  State = 6
	StateAtDestination    State = 7
	StateComplete         State = 8
	StateCompleteEarly    State = 9
	StateCanceled         State = 10
	StateCanceledByDriver State = 11
	StateCanceledByRider  State = 12
	StateRiderNoShow      State = 13
)

var stateString = [...]string{
	StateUnspecified:      "UNSPECIFIED",
	StateAvailable:        "AVAILABLE",
	StateGoToPickup:       "GO_TO_PICKUP",
	StateNearPickup:       "NEAR_PICKUP",
	StateAwaitingPickup:   "AWAITING_PICKUP",
	StateOnTrip:           "ON_TRIP",
	StateNearDestination:  "NEAR_DESTINATION",
	StateAtDestination:    "AT_DESTINATION",
	StateComplete:         "COMPLETE",
	StateCompleteEarly:    "COMPLETE_EARLY",
	StateCanceled:         "CANCELED",
	StateCanceledByDriver: "CANCELED_BY_DRIVER",
	StateCanceledByRider:  "CANCELED_BY_RIDER",
	StateRiderNoShow:      "RIDER_NO_SHOW",
}

var stringState = buildStateMap()

//-----------------------------------------------------------------------------

func StateFromString(s string) (state State, found bool) {
	state, found = stringState[s]
	return
}

//-----------------------------------------------------------------------------

func StateFromNumber(n int) (state State, found bool) {
	if n >= 0 && n < len(stateString) {
		state = State(n)
		found = true
	}
	return
}

//-----------------------------------------------------------------------------

func (s State) String() string {
	if s >= 0 && int(s) < len(stateString) {
		return stateString[s]
	}
	return "<invalid>"
}

//-----------------------------------------------------------------------------

func buildStateMap() map[string]State {
	m := make(map[string]State)
	for i := 0; i < len(stateString); i++ {
		m[stateString[i]] = State(i)
	}
	return m
}

//-----------------------------------------------------------------------------
// State DAO
//-----------------------------------------------------------------------------

func SetState(ctx dao.Context, jobID string, state State) (set bool, err error) {
	sql := dao.Update(table).
		Set("state", state.String()).
		Where("id = ?", jobID)
	err = dao.ExecTX(ctx, sql)
	set = (err == nil) && (dao.RowsAffected(sql.Result()) > 0)
	return
}
