package stat

import (
	"github.com/openmarketplaceengine/openmarketplaceengine/cfg"
	"github.com/openmarketplaceengine/openmarketplaceengine/log"
)

const pfxErr = "stat"

var state cfg.State64

var slog = log.Log()

//-----------------------------------------------------------------------------

func Boot() error {
	if !state.TryBoot() {
		return state.StateError(pfxErr)
	}
	slog = log.Named("STAT")
	_ = slog
	bootProm()
	bootHttp()
	state.SetRunning()
	return nil
}
