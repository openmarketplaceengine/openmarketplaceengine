package log

//-----------------------------------------------------------------------------

func IsDebug() bool {
	return z.Core().Enabled(LevelDebug)
}

func IsInfo() bool {
	return z.Core().Enabled(LevelInfo)
}

func IsWarn() bool {
	return z.Core().Enabled(LevelWarn)
}

func IsError() bool {
	return z.Core().Enabled(LevelError)
}

//-----------------------------------------------------------------------------

func Debugf(format string, args ...interface{}) {
	s.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	s.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	s.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	s.Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	s.Fatalf(format, args...)
}

//-----------------------------------------------------------------------------

func Sync() error {
	return z.Sync()
}

func SafeSync() {
	_ = z.Sync()
}
