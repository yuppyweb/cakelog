package decorator_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/yuppyweb/cakelog"
	"github.com/yuppyweb/cakelog/decorator"
)

func newLevelLogger(t *testing.T, log cakelog.Logger) cakelog.Logger {
	t.Helper()

	logger, err := decorator.NewLevel(log, decorator.LevelDebug)
	if err != nil {
		t.Fatalf("failed to create level logger: %v", err)
	}

	return logger
}

// TestNewLevel_NilLogger tests that NewLevel returns a wrapped ErrNilLogger
// when the logger is nil.
func TestNewLevel_NilLogger(t *testing.T) {
	t.Parallel()

	_, err := decorator.NewLevel(nil, decorator.LevelInfo)
	if err == nil {
		t.Fatal("expected error when providing a nil logger, got nil")
	}

	if !errors.Is(err, decorator.ErrNilLogger) {
		t.Errorf(
			"unexpected error when creating level logger with nil logger:\nGot:  %v\nWant: %v",
			err,
			decorator.ErrNilLogger,
		)
	}

	if !strings.Contains(err.Error(), "level logger:") {
		t.Errorf(
			"error message does not contain expected text:\nGot:  %s\nWant to contain: %s",
			err.Error(),
			"level logger:",
		)
	}
}

// TestNewLevel_TypedNilLogger tests that NewLevel returns a wrapped
// ErrNilLogger when a typed nil logger is provided.
func TestNewLevel_TypedNilLogger(t *testing.T) {
	t.Parallel()

	var typedNil *mockLogger

	_, err := decorator.NewLevel(typedNil, decorator.LevelInfo)
	if err == nil {
		t.Fatal("expected error when providing a typed nil logger, got nil")
	}

	if !errors.Is(err, decorator.ErrNilLogger) {
		t.Errorf(
			"unexpected error when creating level logger with typed nil logger:\nGot:  %v\nWant: %v",
			err,
			decorator.ErrNilLogger,
		)
	}

	if !strings.Contains(err.Error(), "level logger:") {
		t.Errorf(
			"error message does not contain expected text:\nGot:  %s\nWant to contain: %s",
			err.Error(),
			"level logger:",
		)
	}
}

// TestNewLevel_InvalidLevel tests that NewLevel returns a wrapped
// ErrInvalidLogLevel when minLevel is outside the valid range.
func TestNewLevel_InvalidLevel(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		newLevel func(cakelog.Logger) (cakelog.Logger, error)
		value    string
	}{
		{
			name: "below debug",
			newLevel: func(log cakelog.Logger) (cakelog.Logger, error) {
				return decorator.NewLevel(log, decorator.LevelDebug-1)
			},
			value: fmt.Sprintf("%d", decorator.LevelDebug-1),
		},
		{
			name: "above error",
			newLevel: func(log cakelog.Logger) (cakelog.Logger, error) {
				return decorator.NewLevel(log, decorator.LevelError+1)
			},
			value: fmt.Sprintf("%d", decorator.LevelError+1),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, err := tc.newLevel(new(mockLogger))
			if err == nil {
				t.Fatal("expected error when providing an invalid level, got nil")
			}

			if !errors.Is(err, decorator.ErrInvalidLogLevel) {
				t.Errorf(
					"unexpected error when creating level logger with invalid level:\nGot:  %v\nWant: %v",
					err,
					decorator.ErrInvalidLogLevel,
				)
			}

			if !strings.Contains(err.Error(), "level logger:") {
				t.Errorf(
					"error message does not contain expected text:\nGot:  %s\nWant to contain: %s",
					err.Error(),
					"level logger:",
				)
			}

			if !strings.Contains(err.Error(), tc.value) {
				t.Errorf(
					"error message does not contain expected text:\nGot:  %s\nWant to contain: %s",
					err.Error(),
					tc.value,
				)
			}
		})
	}
}

// TestLevelLogger_FiltersByMinLevel tests that each minLevel forwards
// calls at that level and above and drops the rest.
func TestLevelLogger_FiltersByMinLevel(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name      string
		newLogger func(cakelog.Logger) (cakelog.Logger, error)
		wantDebug int
		wantInfo  int
		wantWarn  int
		wantError int
	}{
		{
			name: "debug",
			newLogger: func(log cakelog.Logger) (cakelog.Logger, error) {
				return decorator.NewLevel(log, decorator.LevelDebug)
			},
			wantDebug: 1,
			wantInfo:  1,
			wantWarn:  1,
			wantError: 1,
		},
		{
			name: "info",
			newLogger: func(log cakelog.Logger) (cakelog.Logger, error) {
				return decorator.NewLevel(log, decorator.LevelInfo)
			},
			wantInfo:  1,
			wantWarn:  1,
			wantError: 1,
		},
		{
			name: "warn",
			newLogger: func(log cakelog.Logger) (cakelog.Logger, error) {
				return decorator.NewLevel(log, decorator.LevelWarn)
			},
			wantWarn:  1,
			wantError: 1,
		},
		{
			name: "error",
			newLogger: func(log cakelog.Logger) (cakelog.Logger, error) {
				return decorator.NewLevel(log, decorator.LevelError)
			},
			wantError: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			base := new(mockLogger)

			logger, err := tc.newLogger(base)
			if err != nil {
				t.Fatalf("failed to create level logger: %v", err)
			}

			type ctxKey struct{}

			expectedCtx := context.WithValue(context.Background(), ctxKey{}, tc.name)
			expectedErr := errors.New("error message")

			logger.Debug(expectedCtx, "debug message", "debug", 42)
			logger.Info(expectedCtx, "info message", "info", 75)
			logger.Warn(expectedCtx, "warn message", "warn", 80)
			logger.Error(expectedCtx, expectedErr, "error", 99)

			if len(base.debugIn) != tc.wantDebug {
				t.Errorf("expected %d Debug calls, got %d", tc.wantDebug, len(base.debugIn))
			}

			if len(base.infoIn) != tc.wantInfo {
				t.Errorf("expected %d Info calls, got %d", tc.wantInfo, len(base.infoIn))
			}

			if len(base.warnIn) != tc.wantWarn {
				t.Errorf("expected %d Warn calls, got %d", tc.wantWarn, len(base.warnIn))
			}

			if len(base.errorIn) != tc.wantError {
				t.Errorf("expected %d Error calls, got %d", tc.wantError, len(base.errorIn))
			}

			if tc.wantDebug != 0 {
				assertLoggedMsg(t, base.debugIn, expectedCtx, "debug message", "debug", 42)
			}

			if tc.wantInfo != 0 {
				assertLoggedMsg(t, base.infoIn, expectedCtx, "info message", "info", 75)
			}

			if tc.wantWarn != 0 {
				assertLoggedMsg(t, base.warnIn, expectedCtx, "warn message", "warn", 80)
			}

			if tc.wantError != 0 {
				assertLoggedErr(t, base.errorIn, expectedCtx, expectedErr, "error", 99)
			}
		})
	}
}

// TestLevelLogger_EmptyArgs tests that all methods forward calls with no
// extra arguments, including a nil error, when minLevel is LevelDebug.
func TestLevelLogger_EmptyArgs(t *testing.T) {
	t.Parallel()

	base := new(mockLogger)
	logger := newLevelLogger(t, base)

	ctx := context.Background()

	logger.Debug(ctx, "debug message")
	logger.Info(ctx, "info message")
	logger.Warn(ctx, "warn message")
	logger.Error(ctx, nil)

	assertLoggedMsg(t, base.debugIn, ctx, "debug message")
	assertLoggedMsg(t, base.infoIn, ctx, "info message")
	assertLoggedMsg(t, base.warnIn, ctx, "warn message")
	assertLoggedErr(t, base.errorIn, ctx, nil)
}

// TestLevelLogger_ArgsNotCopied tests that args are forwarded without
// copying, so later mutations of the caller's slice are visible to the
// underlying logger.
func TestLevelLogger_ArgsNotCopied(t *testing.T) {
	t.Parallel()

	base := new(mockLogger)
	logger := newLevelLogger(t, base)

	args := []any{"key", "value"}
	logger.Info(context.Background(), "message", args...)

	args[1] = "mutated"

	if len(base.infoIn) != 1 {
		t.Fatalf("expected Info to be called once, got %d", len(base.infoIn))
	}

	if base.infoIn[0].args[1] != "mutated" {
		t.Errorf(
			"expected underlying logger to see mutated args, got %v",
			base.infoIn[0].args[1],
		)
	}
}

// TestLevelLogger_PanicNotRecovered tests that a panic in the underlying
// logger is not recovered.
func TestLevelLogger_PanicNotRecovered(t *testing.T) {
	t.Parallel()

	base := &hookLogger{
		infoFn: func() {
			panic("level logger panic")
		},
	}

	logger := newLevelLogger(t, base)

	defer func() {
		got := recover()
		if got == nil {
			return
		}

		if got != "level logger panic" {
			t.Errorf("expected panic %q, got %v", "level logger panic", got)
		}
	}()

	logger.Info(context.Background(), "message")

	t.Fatal("expected panic from underlying logger")
}
