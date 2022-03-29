// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package log

import (
	"fmt"
	stdlog "log"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

const (
	compactEncoderName = "compact"
)

type zapLog struct {
	s *zap.SugaredLogger
	c zapcore.Core
	z *zap.Logger
}

var (
	zlog = noopZap()
)

var sink []string
var stdlogReset func()

//-----------------------------------------------------------------------------

func init() {
	err := zap.RegisterEncoder(compactEncoderName, createCompactEncoder)
	if err != nil {
		panic(err)
	}
}

//-----------------------------------------------------------------------------

func noopZap() (z zapLog) {
	z.set(zap.NewNop())
	return
}

func (z *zapLog) set(log *zap.Logger) {
	z.s = log.Sugar()
	z.c = log.Core()
	z.z = log
}

//-----------------------------------------------------------------------------

func NewStdLog(level Level) *stdlog.Logger {
	slog, _ := zap.NewStdLogAt(zlog.z, level)
	return slog
}

//-----------------------------------------------------------------------------

func Init(c ConfigHolder) error {
	reset()
	log, err := newZap(c)
	if err != nil {
		return err
	}
	zlog.set(log)
	stdlogReset = zap.RedirectStdLog(log)
	return nil
}

//-----------------------------------------------------------------------------

func reset() {
	if len(sink) > 0 {
		_ = zlog.Sync()
		sink = nil
	}
	if stdlogReset != nil {
		stdlogReset()
		stdlogReset = nil
	}
}

//-----------------------------------------------------------------------------

func newZap(c ConfigHolder) (*zap.Logger, error) {
	sink = filterEmpty(c.LogTerm(), c.LogFile())
	if len(sink) == 0 {
		return zap.NewNop(), nil
	}
	lev, err := zapLev(c.LogLevel())
	if err != nil {
		return nil, err
	}
	if c.LogTrim() {
		trimFile(sink)
	}
	encoder := compactEncoderName
	if c.LogStyle() == "json" {
		encoder = "json"
	}
	cfg := zap.Config{
		Level:             zap.NewAtomicLevelAt(LevelDebug),
		Development:       c.LogDevel(),
		DisableCaller:     !c.LogCaller(),
		DisableStacktrace: !c.LogTrace(),
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:      encoder,
		EncoderConfig: encoderConfig(c),
		OutputPaths:   sink,
	}
	return cfg.Build(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return &levelCore{int32(lev), core}
	}))
}

//-----------------------------------------------------------------------------

func zapLev(name string) (lev Level, err error) {
	switch name {
	case "debug":
		lev = LevelDebug
	case "info":
		lev = LevelInfo
	case "warn":
		lev = LevelWarn
	case "error":
		lev = LevelError
	case "panic":
		lev = LevelPanic
	case "fatal":
		lev = LevelFatal
	default:
		err = fmt.Errorf("invalid log level: %q", name)
	}
	return
}

//-----------------------------------------------------------------------------

func encoderConfig(c ConfigHolder) zapcore.EncoderConfig {
	enc := zapcore.EncoderConfig{
		TimeKey:          "T",
		LevelKey:         "L",
		NameKey:          "N",
		CallerKey:        "C",
		FunctionKey:      zapcore.OmitKey,
		MessageKey:       "M",
		StacktraceKey:    "S",
		LineEnding:       zapcore.DefaultLineEnding,
		EncodeLevel:      zapcore.CapitalLevelEncoder,
		EncodeTime:       zapcore.ISO8601TimeEncoder,
		EncodeDuration:   zapcore.StringDurationEncoder,
		EncodeCaller:     zapcore.ShortCallerEncoder,
		ConsoleSeparator: " ",
	}
	if !c.LogTrace() {
		enc.StacktraceKey = zapcore.OmitKey
	}
	if !c.LogCaller() {
		enc.CallerKey = zapcore.OmitKey
	}
	return enc
}

//-----------------------------------------------------------------------------
// compactEncoder
//-----------------------------------------------------------------------------

type compactEncoder struct {
	*zapcore.EncoderConfig
}

//-----------------------------------------------------------------------------

func createCompactEncoder(c zapcore.EncoderConfig) (zapcore.Encoder, error) {
	return compactEncoder{&c}, nil
}

//-----------------------------------------------------------------------------

func (c compactEncoder) AddArray(_ string, _ zapcore.ArrayMarshaler) error   { return nil }
func (c compactEncoder) AddObject(_ string, _ zapcore.ObjectMarshaler) error { return nil }
func (c compactEncoder) AddReflected(_ string, _ interface{}) error          { return nil }

func (c compactEncoder) AddBinary(_ string, _ []byte)          {}
func (c compactEncoder) AddByteString(_ string, _ []byte)      {}
func (c compactEncoder) AddBool(_ string, _ bool)              {}
func (c compactEncoder) AddComplex128(_ string, _ complex128)  {}
func (c compactEncoder) AddComplex64(_ string, _ complex64)    {}
func (c compactEncoder) AddDuration(_ string, _ time.Duration) {}
func (c compactEncoder) AddFloat64(_ string, _ float64)        {}
func (c compactEncoder) AddFloat32(_ string, _ float32)        {}
func (c compactEncoder) AddInt(_ string, _ int)                {}
func (c compactEncoder) AddInt64(_ string, _ int64)            {}
func (c compactEncoder) AddInt32(_ string, _ int32)            {}
func (c compactEncoder) AddInt16(_ string, _ int16)            {}
func (c compactEncoder) AddInt8(_ string, _ int8)              {}
func (c compactEncoder) AddString(_, _ string)                 {}
func (c compactEncoder) AddTime(_ string, _ time.Time)         {}
func (c compactEncoder) AddUint(_ string, _ uint)              {}
func (c compactEncoder) AddUint64(_ string, _ uint64)          {}
func (c compactEncoder) AddUint32(_ string, _ uint32)          {}
func (c compactEncoder) AddUint16(_ string, _ uint16)          {}
func (c compactEncoder) AddUint8(_ string, _ uint8)            {}
func (c compactEncoder) AddUintptr(_ string, _ uintptr)        {}
func (c compactEncoder) OpenNamespace(_ string)                {}

//-----------------------------------------------------------------------------

func (c compactEncoder) Clone() zapcore.Encoder {
	return compactEncoder{c.EncoderConfig}
}

//-----------------------------------------------------------------------------

var bpool = buffer.NewPool()

func (c compactEncoder) EncodeEntry(e zapcore.Entry, _ []zapcore.Field) (*buffer.Buffer, error) {
	b := bpool.Get()
	b.AppendTime(e.Time, TimeFormat)
	c.sep(b)
	b.AppendString(zapLevStr(e.Level))
	if len(e.LoggerName) > 0 {
		c.sep(b)
		b.AppendByte('[')
		b.AppendString(e.LoggerName)
		b.AppendByte(']')
	}
	hasMsg := false
	if m, n := trimEOL(e.Message); n > 0 {
		c.sep(b)
		b.AppendString(m)
		hasMsg = true
	}
	if e.Caller.Defined && len(c.CallerKey) > 0 {
		if hasMsg {
			b.AppendByte(' ')
		}
		b.AppendByte('(')
		b.AppendString(e.Caller.TrimmedPath())
		b.AppendByte(')')
		hasMsg = true
	}
	if len(e.Stack) > 0 && len(c.StacktraceKey) > 0 {
		if hasMsg {
			b.AppendString(c.LineEnding)
		}
		indent(b, e.Stack, "    ")
		b.TrimNewline()
	}
	b.AppendString(c.LineEnding)
	return b, nil
}

func (c compactEncoder) sep(b *buffer.Buffer) {
	b.AppendString(c.ConsoleSeparator)
}

//-----------------------------------------------------------------------------

func zapLevStr(lev zapcore.Level) string {
	switch lev {
	case zapcore.DebugLevel:
		return "[DBG]"
	case zapcore.InfoLevel:
		return "[INF]"
	case zapcore.WarnLevel:
		return "[WRN]"
	case zapcore.ErrorLevel:
		return "[ERR]"
	case zapcore.PanicLevel:
	case zapcore.DPanicLevel:
		return "[PANIC]"
	case zapcore.FatalLevel:
		return "[FATAL]"
	}
	return "[INVAL]"
}

//-----------------------------------------------------------------------------

func filterEmpty(arg ...string) []string {
	ret := make([]string, 0, len(arg))
	for _, s := range arg {
		if len(s) > 0 {
			ret = append(ret, s)
		}
	}
	return ret
}

//-----------------------------------------------------------------------------

func trimEOL(s string) (string, int) {
	n := len(s)
	if n > 0 && s[n-1] == '\n' {
		s = s[:n-1]
		n--
	}
	return s, n
}

//-----------------------------------------------------------------------------

func indent(b *buffer.Buffer, s string, ind string) {
	for len(s) > 0 {
		b.AppendString(ind)
		x := strings.IndexByte(s, '\n')
		if x == -1 {
			b.AppendString(s)
			return
		}
		if x == 0 {
			b.AppendByte('\n')
			s = s[1:]
			continue
		}
		b.AppendString(s[:x+1])
		s = s[x+1:]
	}
}

//-----------------------------------------------------------------------------

func trimFile(files []string) {
	for _, name := range files {
		if strings.HasSuffix(name, FileExtn) {
			_ = os.Remove(name)
		}
	}
}

//-----------------------------------------------------------------------------

type levelCore struct {
	level int32
	core  zapcore.Core
}

func (c *levelCore) Level() Level {
	return Level(int8(atomic.LoadInt32(&c.level)))
}

func (c *levelCore) SetLevel(level Level) {
	atomic.StoreInt32(&c.level, int32(level))
}

func (c *levelCore) Enabled(level Level) bool {
	return level >= c.Level()
}

func (c *levelCore) With(fields []zapcore.Field) zapcore.Core {
	level := atomic.LoadInt32(&c.level)
	return &levelCore{level, c.core.With(fields)}
}

func (c *levelCore) Check(entry zapcore.Entry, checked *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if !c.Enabled(entry.Level) {
		return checked
	}
	return c.core.Check(entry, checked)
}

func (c *levelCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	return c.core.Write(entry, fields)
}

func (c *levelCore) Sync() error {
	return c.core.Sync()
}
