// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package log

const (
	TimeFormat = "2006/01/02 15:04:05.000"
	FileExtn   = ".log"
)

// ConfigHolder interface provides log configuration params.
type ConfigHolder interface {
	LogLevel() string
	LogTrace() bool
	LogDevel() bool
	LogStyle() string // json | plain
	LogTerm() string  // stdout | stderr
	LogFile() string
	LogTrim() bool
	LogCaller() bool
}

// ConfigValues implements ConfigHolder interface for log configuration.
type ConfigValues struct {
	Level  string
	Style  string
	Term   string
	File   string
	Trim   bool
	Devel  bool
	Trace  bool
	Caller bool
}

func (c *ConfigValues) LogLevel() string {
	return c.Level
}

func (c *ConfigValues) LogTrace() bool {
	return c.Trace
}

func (c *ConfigValues) LogDevel() bool {
	return c.Devel
}

func (c *ConfigValues) LogStyle() string {
	return c.Style
}

func (c *ConfigValues) LogTerm() string {
	return c.Term
}

func (c *ConfigValues) LogFile() string {
	return c.File
}

func (c *ConfigValues) LogTrim() bool {
	return c.Trim
}

func (c *ConfigValues) LogCaller() bool {
	return c.Caller
}

//-----------------------------------------------------------------------------

func (c *ConfigValues) WithLevel(level string) *ConfigValues {
	c.Level = level
	return c
}

func (c *ConfigValues) WithStyle(style string) *ConfigValues {
	c.Style = style
	return c
}

func (c *ConfigValues) WithTerm(term string) *ConfigValues {
	c.Term = term
	return c
}

func (c *ConfigValues) WithFile(file string) *ConfigValues {
	c.File = file
	return c
}

func (c *ConfigValues) WithTrim(trim bool) *ConfigValues {
	c.Trim = trim
	return c
}

func (c *ConfigValues) WithDevel(devel bool) *ConfigValues {
	c.Devel = devel
	return c
}

func (c *ConfigValues) WithTrace(trace bool) *ConfigValues {
	c.Trace = trace
	return c
}

func (c *ConfigValues) WithCaller(caller bool) *ConfigValues {
	c.Caller = caller
	return c
}

//-----------------------------------------------------------------------------

func DevelConfig() *ConfigValues {
	return &ConfigValues{
		Level:  "debug",
		Style:  "plain",
		Term:   "stdout",
		File:   "",
		Trim:   false,
		Devel:  true,
		Trace:  true,
		Caller: true,
	}
}
