// Package cakelog provides a simple logging interface with multiple adapter implementations.
package cakelog

import (
	"context"
)

// Logger defines the logging interface with methods for different log levels.
type Logger interface {
	// Debug logs a debug-level message.
	Debug(ctx context.Context, msg string, args ...any)
	// Info logs an info-level message.
	Info(ctx context.Context, msg string, args ...any)
	// Warn logs a warning-level message.
	Warn(ctx context.Context, msg string, args ...any)
	// Error logs an error-level message with an error value.
	Error(ctx context.Context, err error, args ...any)
}

// NopLogger is a no-operation logger that discards all log messages.
type NopLogger struct{}

// NewNopLogger returns a new instance of NopLogger.
func NewNopLogger() *NopLogger {
	return &NopLogger{}
}

// Debug implements the Logger interface (no-op).
func (*NopLogger) Debug(context.Context, string, ...any) {}

// Info implements the Logger interface (no-op).
func (*NopLogger) Info(context.Context, string, ...any) {}

// Warn implements the Logger interface (no-op).
func (*NopLogger) Warn(context.Context, string, ...any) {}

// Error implements the Logger interface (no-op).
func (*NopLogger) Error(context.Context, error, ...any) {}

// Ensure NopLogger implements the Logger interface.
var _ Logger = (*NopLogger)(nil)
