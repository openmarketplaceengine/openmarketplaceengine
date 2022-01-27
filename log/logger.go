package log

import (
	"fmt"

	"go.uber.org/zap"
)

var logger *zap.Logger

func init() {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		panic(fmt.Sprintf("create logger error: %s", err))
	}
}

func GetLogger() *zap.Logger {
	return logger
}
