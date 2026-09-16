// Package cakelog provides a unified Logger interface for application code.
//
// Swap the backend through adapters and add behavior with stackable
// decorators. Application code depends on Logger, not on a specific
// logging library.
package cakelog

import (
	"context"
)

// Logger is the unified logging interface used by application code,
// adapters, and decorators.
//
// Args may mix alternating key-value pairs and maps. Built-in adapters
// parse them the same way. A backend may then keep duplicate keys or
// collapse them, depending on its own field model.
//
// Error has no separate message argument. Built-in adapters log
// err.Error(), or an empty message when err is nil.
//
// ctx is passed on every call. An adapter forwards it when the backend
// accepts context, and otherwise ignores it.
type Logger interface {
	// Debug logs a debug-level message.
	Debug(ctx context.Context, msg string, args ...any)
	// Info logs an info-level message.
	Info(ctx context.Context, msg string, args ...any)
	// Warn logs a warning-level message.
	Warn(ctx context.Context, msg string, args ...any)
	// Error logs err at error level. There is no separate message argument.
	// err may be nil.
	Error(ctx context.Context, err error, args ...any)
}

// nopLogger discards every log call. Methods never panic.
type nopLogger struct{}

// NopLogger returns a logger that discards every call. It is safe to use
// in tests and as a disabled backend. Methods never panic, including when
// err is nil and when args mix key-value pairs and maps.
func NopLogger() Logger {
	return nopLogger{}
}

// Debug discards the debug-level message.
func (nopLogger) Debug(context.Context, string, ...any) {}

// Info discards the info-level message.
func (nopLogger) Info(context.Context, string, ...any) {}

// Warn discards the warning-level message.
func (nopLogger) Warn(context.Context, string, ...any) {}

// Error discards the error-level call. err may be nil.
func (nopLogger) Error(context.Context, error, ...any) {}

// Ensure nopLogger implements the Logger interface.
var _ Logger = nopLogger{}
