package stat

import (
	"github.com/driverscooperative/geosrv/cfg"
	"github.com/driverscooperative/geosrv/log"
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
	bootHTTP()
	state.SetRunning()
	return nil
}
