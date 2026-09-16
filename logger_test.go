package cakelog_test

import (
	"context"
	"errors"
	"testing"

	"github.com/yuppyweb/cakelog"
)

// TestNopLogger tests that NopLogger returns a non-nil logger and that
// Debug, Info, Warn, and Error accept typical arguments without panicking.
func TestNopLogger(t *testing.T) {
	t.Parallel()

	logger := cakelog.NopLogger()
	if logger == nil {
		t.Fatal("expected NopLogger to return a non-nil logger")
	}

	type ctxKey struct{}

	ctx := context.WithValue(context.Background(), ctxKey{}, "nop")
	err := errors.New("error message")

	logger.Debug(ctx, "debug message", "debug", 42)
	logger.Info(ctx, "info message", "info", 75)
	logger.Warn(ctx, "warn message", "warn", 88)
	logger.Error(ctx, err, "error", 90)
}

// TestNopLogger_EmptyArgs tests that all methods accept calls with no
// extra arguments, including a nil error.
func TestNopLogger_EmptyArgs(t *testing.T) {
	t.Parallel()

	logger := cakelog.NopLogger()
	ctx := context.Background()

	logger.Debug(ctx, "debug message")
	logger.Info(ctx, "info message")
	logger.Warn(ctx, "warn message")
	logger.Error(ctx, nil)
}

// TestNopLogger_MixedArgs tests that all methods accept mixed key-value
// pairs and maps, including an empty message.
func TestNopLogger_MixedArgs(t *testing.T) {
	t.Parallel()

	logger := cakelog.NopLogger()
	ctx := context.Background()
	fields := map[string]any{"status": 200, "path": "/health"}

	logger.Debug(ctx, "", "method", "GET", fields)
	logger.Info(ctx, "request finished", "method", "GET", fields)
	logger.Warn(ctx, "slow request", fields, "duration", 150)
	logger.Error(ctx, errors.New("request failed"), fields, "retry", true)
}
