package log

// Logger interface describes logging functions.
type Logger interface {
	IsDebug() bool
	IsInfo() bool
	IsWarn() bool
	IsError() bool
	IsLevel(lev Level) bool
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Panicf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Sync() error
}

// Log returns default Logger implementation.
func Log() Logger {
	return &z
}

//-----------------------------------------------------------------------------

func IsLevel(level Level) bool {
	return z.IsLevel(level)
}

func IsDebug() bool {
	return z.IsDebug()
}

func IsInfo() bool {
	return z.IsInfo()
}

func IsWarn() bool {
	return z.IsWarn()
}

func IsError() bool {
	return z.IsError()
}

//-----------------------------------------------------------------------------

func Debugf(format string, args ...interface{}) {
	z.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	z.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	z.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	z.Errorf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	z.Panicf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	z.Fatalf(format, args...)
}

//-----------------------------------------------------------------------------

func Sync() error {
	return z.Sync()
}

func SafeSync() {
	_ = z.Sync()
}

//-----------------------------------------------------------------------------
// Zap Log
//-----------------------------------------------------------------------------

func (z *zapLog) IsLevel(level Level) bool {
	return z.c.Enabled(level)
}

func (z *zapLog) IsDebug() bool {
	return z.c.Enabled(LevelDebug)
}

func (z *zapLog) IsInfo() bool {
	return z.c.Enabled(LevelInfo)
}

func (z *zapLog) IsWarn() bool {
	return z.c.Enabled(LevelWarn)
}

func (z *zapLog) IsError() bool {
	return z.c.Enabled(LevelError)
}

//-----------------------------------------------------------------------------

func (z *zapLog) Debugf(format string, args ...interface{}) {
	z.s.Debugf(format, args...)
}

func (z *zapLog) Infof(format string, args ...interface{}) {
	z.s.Infof(format, args...)
}

func (z *zapLog) Warnf(format string, args ...interface{}) {
	z.s.Warnf(format, args...)
}

func (z *zapLog) Errorf(format string, args ...interface{}) {
	z.s.Errorf(format, args...)
}

func (z *zapLog) Panicf(format string, args ...interface{}) {
	z.s.Panicf(format, args...)
}

func (z *zapLog) Fatalf(format string, args ...interface{}) {
	z.s.Fatalf(format, args...)
}

//-----------------------------------------------------------------------------

func (z *zapLog) Sync() error {
	return z.z.Sync()
}