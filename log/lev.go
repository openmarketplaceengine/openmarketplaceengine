// Copyright 2022 The Drivers Cooperative. All rights reserved.
// Use of this source code is governed by a dual
// license that can be found in the LICENSE file.

package log

import "go.uber.org/zap/zapcore"

type Level = zapcore.Level

const (
	// LevelDebug logs are typically voluminous, and are usually disabled in
	// production.
	LevelDebug = zapcore.DebugLevel
	// LevelInfo is the default logging priority.
	LevelInfo = zapcore.InfoLevel
	// LevelWarn logs are more important than Info, but don't need individual
	// human review.
	LevelWarn = zapcore.WarnLevel
	// LevelError logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	LevelError = zapcore.ErrorLevel
	// LevelPanic logs a message, then panics.
	LevelPanic = zapcore.PanicLevel
	// LevelFatal logs a message, then calls os.Exit(1).
	LevelFatal = zapcore.FatalLevel
	// LevelNone disables logging
	LevelNone = LevelFatal + 1
)
