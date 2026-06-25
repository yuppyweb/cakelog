package cakelog_test

import (
	"context"
	"errors"
	"testing"

	"github.com/yuppyweb/cakelog"
)

// TestNopLogger tests the NopLogger implementation to ensure all logging methods execute without errors.
func TestNopLogger(t *testing.T) {
	t.Parallel()

	logger := cakelog.NewNopLogger()

	logger.Debug(context.Background(), "debug message", "debug", 42)
	logger.Info(context.Background(), "info message", "info", 75)
	logger.Warn(context.Background(), "warn message", "warn", 88)
	logger.Error(context.Background(), errors.New("error message"), "error", 90)
}
