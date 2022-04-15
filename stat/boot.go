package stat

import "github.com/openmarketplaceengine/openmarketplaceengine/cfg"

const pfxErr = "stat"

var state cfg.State64

//-----------------------------------------------------------------------------

func Boot() error {
	if !state.TryBoot() {
		return state.StateError(pfxErr)
	}
	bootProm()
	state.SetRunning()
	return nil
}
