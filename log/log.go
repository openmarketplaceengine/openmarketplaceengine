package log

type PrintFunc = func(format string, args ...interface{})

// Logger interface describes logging functions.
type Logger interface {
	IsDebug() bool
	IsInfo() bool
	IsWarn() bool
	IsError() bool
	IsLevel(lev Level) bool
	LevelFunc(lev Level) PrintFunc
	Levelf(level Level, format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Panicf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Named(s string) Logger
	Sync() error
}

// Log returns default Logger implementation.
func Log() Logger {
	return &zlog
}

//-----------------------------------------------------------------------------

func IsLevel(level Level) bool {
	return zlog.IsLevel(level)
}

func IsDebug() bool {
	return zlog.IsDebug()
}

func IsInfo() bool {
	return zlog.IsInfo()
}

func IsWarn() bool {
	return zlog.IsWarn()
}

func IsError() bool {
	return zlog.IsError()
}

func LevelFunc(lev Level) PrintFunc {
	return zlog.LevelFunc(lev)
}

func Levelf(level Level, format string, args ...interface{}) {
	zlog.Levelf(level, format, args...)
}

func Named(s string) Logger {
	return zlog.Named(s)
}

//-----------------------------------------------------------------------------

func Debugf(format string, args ...interface{}) {
	zlog.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	zlog.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	zlog.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	zlog.Errorf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	zlog.Panicf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	zlog.Fatalf(format, args...)
}

//-----------------------------------------------------------------------------

func Sync() error {
	return zlog.Sync()
}

func SafeSync() {
	_ = zlog.Sync()
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

func (z *zapLog) LevelFunc(lev Level) PrintFunc {
	switch lev {
	case LevelDebug:
		return z.Debugf
	case LevelInfo:
		return z.Infof
	case LevelWarn:
		return z.Warnf
	case LevelError:
		return z.Errorf
	case LevelPanic:
		return z.Panicf
	case LevelFatal:
		return Fatalf
	default:
		return nil
	}
}

func (z *zapLog) Levelf(level Level, format string, args ...interface{}) {
	switch level {
	case LevelDebug:
		z.Debugf(format, args...)
	case LevelInfo:
		z.Infof(format, args...)
	case LevelWarn:
		z.Warnf(format, args...)
	case LevelError:
		z.Errorf(format, args...)
	case LevelPanic:
		z.Panicf(format, args...)
	case LevelFatal:
		z.Fatalf(format, args...)
	}
}

//-----------------------------------------------------------------------------

func (z *zapLog) Named(s string) Logger {
	named := new(zapLog)
	named.set(z.z.Named(s))
	return named
}

//-----------------------------------------------------------------------------

func (z *zapLog) Sync() error {
	return z.z.Sync()
}
