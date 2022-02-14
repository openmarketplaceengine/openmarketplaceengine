// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package cfg

import (
	"fmt"
)

type PrintFunc = func(format string, args ...interface{})

// Logger is a basic subset of logging functions used in this package.
type Logger interface {
	IsDebug() bool
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

// LogConfig contains logging configuration params.
type LogConfig struct {
	Level  string `default:"info" usage:"log #level#: [debug|info|warn|error|panic|fatal]"`
	Style  string `default:"plain" usage:"log format #style#: [plain|json]"`
	Term   string `default:"stdout" usage:"log standard output name: #[stdout|stderr|null]#"`
	File   string `usage:"log file #path#"`
	Trim   bool   `default:"false" usage:"#enable# trimming of the log file on start"`
	Trace  bool   `default:"false" usage:"#enable# logging error stack traces"`
	Devel  bool   `default:"false" usage:"#enable# developer's mode"`
	Caller bool   `default:"false" usage:"#enable# logging the caller function name"`
}

// Check validates LogConfig params.
func (c *LogConfig) Check(name ...string) error {
	const nullSink = "null"
	if len(c.Level) == 0 {
		return fmt.Errorf("%s: empty level", field(name))
	}
	if !c.levelValid(c.Level) {
		return fmt.Errorf("%s: invalid level name: %q", field(name), c.Level)
	}
	if !matchString(c.Style, "plain", "json") {
		return fmt.Errorf("%s: invalid style name: %q", field(name), c.Style)
	}
	if !matchString(c.Term, "stdout", "stderr", nullSink, "") {
		return fmt.Errorf("%s: invalid standard output: %q", field(name), c.Term)
	}
	if c.Term == nullSink {
		c.Term = ""
	}
	if c.File == nullSink {
		c.File = ""
	}
	return nil
}

//-----------------------------------------------------------------------------

func (c *LogConfig) levelValid(level string) bool {
	switch level {
	case "debug", "info", "warn", "error", "panic", "fatal":
		return true
	}
	return false
}

//-----------------------------------------------------------------------------
// log.Config Implementation
//-----------------------------------------------------------------------------

func (c *LogConfig) LogLevel() string {
	return c.Level
}

func (c *LogConfig) LogTrace() bool {
	return c.Trace
}

func (c *LogConfig) LogDevel() bool {
	return c.Devel
}

func (c *LogConfig) LogStyle() string {
	return c.Style
}

func (c *LogConfig) LogTerm() string {
	return c.Term
}

func (c *LogConfig) LogFile() string {
	return c.File
}

func (c *LogConfig) LogTrim() bool {
	return c.Trim
}

func (c *LogConfig) LogCaller() bool {
	return c.Caller
}
