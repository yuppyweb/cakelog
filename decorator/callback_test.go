package decorator_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/yuppyweb/cakelog"
	"github.com/yuppyweb/cakelog/decorator"
)

func newCallbackLogger(
	t *testing.T,
	log cakelog.Logger,
	opts ...decorator.CallbackOption,
) cakelog.Logger {
	t.Helper()

	logger, err := decorator.NewCallback(log, opts...)
	if err != nil {
		t.Fatalf("failed to create callback logger: %v", err)
	}

	return logger
}

func assertCallbackError(t *testing.T, err error, want error) {
	t.Helper()

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, want) {
		t.Errorf(
			"unexpected error when creating callback logger:\nGot:  %v\nWant: %v",
			err,
			want,
		)
	}

	if !strings.Contains(err.Error(), "callback logger:") {
		t.Errorf(
			"error message does not contain expected text:\nGot:  %s\nWant to contain: %s",
			err.Error(),
			"callback logger:",
		)
	}
}

// TestNewCallback_NilLogger tests that NewCallback returns a wrapped
// ErrNilLogger when the logger is nil.
func TestNewCallback_NilLogger(t *testing.T) {
	t.Parallel()

	_, err := decorator.NewCallback(nil)

	assertCallbackError(t, err, decorator.ErrNilLogger)
}

// TestNewCallback_TypedNilLogger tests that NewCallback returns a wrapped
// ErrNilLogger when a typed nil logger is provided.
func TestNewCallback_TypedNilLogger(t *testing.T) {
	t.Parallel()

	var typedNil *mockLogger

	_, err := decorator.NewCallback(typedNil)

	assertCallbackError(t, err, decorator.ErrNilLogger)
}

// TestNewCallback_NilOption tests that NewCallback returns a wrapped
// ErrNilCallbackOption when an option is nil.
func TestNewCallback_NilOption(t *testing.T) {
	t.Parallel()

	_, err := decorator.NewCallback(new(mockLogger), nil)

	assertCallbackError(t, err, decorator.ErrNilCallbackOption)
}

// TestNewCallback_NilCallbackFunc tests that NewCallback returns a wrapped
// sentinel error when a nil callback is configured for a level.
func TestNewCallback_NilCallbackFunc(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		opt  decorator.CallbackOption
		want error
	}{
		{
			name: "debug",
			opt:  decorator.WithDebugCallback(nil),
			want: decorator.ErrNilCallbackDebugFunc,
		},
		{
			name: "info",
			opt:  decorator.WithInfoCallback(nil),
			want: decorator.ErrNilCallbackInfoFunc,
		},
		{
			name: "warn",
			opt:  decorator.WithWarnCallback(nil),
			want: decorator.ErrNilCallbackWarnFunc,
		},
		{
			name: "error",
			opt:  decorator.WithErrorCallback(nil),
			want: decorator.ErrNilCallbackErrorFunc,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, err := decorator.NewCallback(new(mockLogger), tc.opt)

			assertCallbackError(t, err, tc.want)
		})
	}
}

// TestNewCallback_DefaultNoop tests that NewCallback with no options
// forwards log calls and does not panic.
func TestNewCallback_DefaultNoop(t *testing.T) {
	t.Parallel()

	base := new(mockLogger)
	logger := newCallbackLogger(t, base)
	ctx := context.Background()
	expectedErr := errors.New("test error")

	logger.Debug(ctx, "debug message")
	logger.Info(ctx, "info message")
	logger.Warn(ctx, "warn message")
	logger.Error(ctx, expectedErr)

	assertLoggedMsg(t, base.debugIn, ctx, "debug message")
	assertLoggedMsg(t, base.infoIn, ctx, "info message")
	assertLoggedMsg(t, base.warnIn, ctx, "warn message")
	assertLoggedErr(t, base.errorIn, ctx, expectedErr)
}

// TestNewCallback_LastOptionWins tests that the last option for a level
// replaces earlier callbacks for that level.
func TestNewCallback_LastOptionWins(t *testing.T) {
	t.Parallel()

	base := new(mockLogger)
	first := false
	second := false

	logger := newCallbackLogger(
		t,
		base,
		decorator.WithDebugCallback(func(context.Context, string, ...any) {
			first = true
		}),
		decorator.WithDebugCallback(func(context.Context, string, ...any) {
			second = true
		}),
	)

	logger.Debug(context.Background(), "debug message")

	if first {
		t.Errorf("expected the first debug callback not to run")
	}

	if !second {
		t.Errorf("expected the last debug callback to run")
	}
}

// TestCallbackLogger_Debug tests that Debug forwards to the underlying
// logger, then invokes the debug callback with the same arguments.
func TestCallbackLogger_Debug(t *testing.T) {
	t.Parallel()

	type ctxKey struct{}

	base := new(mockLogger)
	expectedCtx := context.WithValue(context.Background(), ctxKey{}, "test debug")
	execed := false

	logger := newCallbackLogger(
		t,
		base,
		decorator.WithDebugCallback(func(ctx context.Context, msg string, args ...any) {
			if ctx != expectedCtx {
				t.Errorf("expected callback context to be %v, got %v", expectedCtx, ctx)
			}

			if msg != "debug message" {
				t.Errorf("expected callback message to be %q, got %q", "debug message", msg)
			}

			if len(args) != 2 || args[0] != "debug" || args[1] != 42 {
				t.Errorf("expected callback args [debug 42], got %v", args)
			}

			execed = true
		}),
	)

	logger.Debug(expectedCtx, "debug message", "debug", 42)

	assertLoggedMsg(t, base.debugIn, expectedCtx, "debug message", "debug", 42)

	if !execed {
		t.Errorf("expected debug callback to be executed")
	}
}

// TestCallbackLogger_Info tests that Info forwards to the underlying
// logger, then invokes the info callback with the same arguments.
func TestCallbackLogger_Info(t *testing.T) {
	t.Parallel()

	type ctxKey struct{}

	base := new(mockLogger)
	expectedCtx := context.WithValue(context.Background(), ctxKey{}, "test info")
	execed := false

	logger := newCallbackLogger(
		t,
		base,
		decorator.WithInfoCallback(func(ctx context.Context, msg string, args ...any) {
			if ctx != expectedCtx {
				t.Errorf("expected callback context to be %v, got %v", expectedCtx, ctx)
			}

			if msg != "info message" {
				t.Errorf("expected callback message to be %q, got %q", "info message", msg)
			}

			if len(args) != 2 || args[0] != "info" || args[1] != 75 {
				t.Errorf("expected callback args [info 75], got %v", args)
			}

			execed = true
		}),
	)

	logger.Info(expectedCtx, "info message", "info", 75)

	assertLoggedMsg(t, base.infoIn, expectedCtx, "info message", "info", 75)

	if !execed {
		t.Errorf("expected info callback to be executed")
	}
}

// TestCallbackLogger_Warn tests that Warn forwards to the underlying
// logger, then invokes the warn callback with the same arguments.
func TestCallbackLogger_Warn(t *testing.T) {
	t.Parallel()

	type ctxKey struct{}

	base := new(mockLogger)
	expectedCtx := context.WithValue(context.Background(), ctxKey{}, "test warn")
	execed := false

	logger := newCallbackLogger(
		t,
		base,
		decorator.WithWarnCallback(func(ctx context.Context, msg string, args ...any) {
			if ctx != expectedCtx {
				t.Errorf("expected callback context to be %v, got %v", expectedCtx, ctx)
			}

			if msg != "warn message" {
				t.Errorf("expected callback message to be %q, got %q", "warn message", msg)
			}

			if len(args) != 2 || args[0] != "warn" || args[1] != 80 {
				t.Errorf("expected callback args [warn 80], got %v", args)
			}

			execed = true
		}),
	)

	logger.Warn(expectedCtx, "warn message", "warn", 80)

	assertLoggedMsg(t, base.warnIn, expectedCtx, "warn message", "warn", 80)

	if !execed {
		t.Errorf("expected warn callback to be executed")
	}
}

// TestCallbackLogger_Error tests that Error forwards to the underlying
// logger, then invokes the error callback with the same arguments.
func TestCallbackLogger_Error(t *testing.T) {
	t.Parallel()

	type ctxKey struct{}

	base := new(mockLogger)
	expectedCtx := context.WithValue(context.Background(), ctxKey{}, "test error")
	expectedErr := errors.New("test error")
	execed := false

	logger := newCallbackLogger(
		t,
		base,
		decorator.WithErrorCallback(func(ctx context.Context, err error, args ...any) {
			if ctx != expectedCtx {
				t.Errorf("expected callback context to be %v, got %v", expectedCtx, ctx)
			}

			if !errors.Is(err, expectedErr) {
				t.Errorf("expected callback error to be %v, got %v", expectedErr, err)
			}

			if len(args) != 2 || args[0] != "error" || args[1] != 99 {
				t.Errorf("expected callback args [error 99], got %v", args)
			}

			execed = true
		}),
	)

	logger.Error(expectedCtx, expectedErr, "error", 99)

	assertLoggedErr(t, base.errorIn, expectedCtx, expectedErr, "error", 99)

	if !execed {
		t.Errorf("expected error callback to be executed")
	}
}

// TestCallbackLogger_NilError tests that an Error callback may receive a
// nil error.
func TestCallbackLogger_NilError(t *testing.T) {
	t.Parallel()

	base := new(mockLogger)
	execed := false

	logger := newCallbackLogger(
		t,
		base,
		decorator.WithErrorCallback(func(_ context.Context, err error, _ ...any) {
			if err != nil {
				t.Errorf("expected nil error, got %v", err)
			}

			execed = true
		}),
	)

	logger.Error(context.Background(), nil)

	assertLoggedErr(t, base.errorIn, context.Background(), nil)

	if !execed {
		t.Errorf("expected error callback to be executed")
	}
}

// TestCallbackLogger_AllCallbacks tests that only the callback for the
// called level runs.
func TestCallbackLogger_AllCallbacks(t *testing.T) {
	t.Parallel()

	base := new(mockLogger)
	ctx := context.Background()
	execedDebug := false
	execedInfo := false
	execedWarn := false
	execedError := false

	logger := newCallbackLogger(
		t,
		base,
		decorator.WithDebugCallback(func(context.Context, string, ...any) {
			execedDebug = true
		}),
		decorator.WithInfoCallback(func(context.Context, string, ...any) {
			execedInfo = true
		}),
		decorator.WithWarnCallback(func(context.Context, string, ...any) {
			execedWarn = true
		}),
		decorator.WithErrorCallback(func(context.Context, error, ...any) {
			execedError = true
		}),
	)

	logger.Debug(ctx, "debug message")

	if !execedDebug {
		t.Errorf("expected debug callback to be executed")
	}

	if execedInfo || execedWarn || execedError {
		t.Errorf("expected only debug callback to be executed")
	}

	logger.Info(ctx, "info message")

	if !execedInfo {
		t.Errorf("expected info callback to be executed")
	}

	if execedWarn || execedError {
		t.Errorf("expected warn and error callbacks not to be executed yet")
	}

	logger.Warn(ctx, "warn message")

	if !execedWarn {
		t.Errorf("expected warn callback to be executed")
	}

	if execedError {
		t.Errorf("expected error callback not to be executed yet")
	}

	logger.Error(ctx, errors.New("test error"))

	if !execedError {
		t.Errorf("expected error callback to be executed")
	}
}

// TestCallbackLogger_MultipleCallbacks tests that each callback runs once
// per log call.
func TestCallbackLogger_MultipleCallbacks(t *testing.T) {
	t.Parallel()

	base := new(mockLogger)
	ctx := context.Background()
	countDebug := 0
	countInfo := 0
	countWarn := 0
	countError := 0

	logger := newCallbackLogger(
		t,
		base,
		decorator.WithDebugCallback(func(context.Context, string, ...any) {
			countDebug++
		}),
		decorator.WithInfoCallback(func(context.Context, string, ...any) {
			countInfo++
		}),
		decorator.WithWarnCallback(func(context.Context, string, ...any) {
			countWarn++
		}),
		decorator.WithErrorCallback(func(context.Context, error, ...any) {
			countError++
		}),
	)

	for range 5 {
		logger.Debug(ctx, "debug message")
		logger.Info(ctx, "info message")
		logger.Warn(ctx, "warn message")
		logger.Error(ctx, errors.New("test error"))
	}

	if countDebug != 5 {
		t.Errorf("expected debug callback to be executed 5 times, got %d", countDebug)
	}

	if countInfo != 5 {
		t.Errorf("expected info callback to be executed 5 times, got %d", countInfo)
	}

	if countWarn != 5 {
		t.Errorf("expected warn callback to be executed 5 times, got %d", countWarn)
	}

	if countError != 5 {
		t.Errorf("expected error callback to be executed 5 times, got %d", countError)
	}
}

// TestCallbackLogger_CallbackAfterLog tests that the callback runs after
// the underlying logger has received the call.
func TestCallbackLogger_CallbackAfterLog(t *testing.T) {
	t.Parallel()

	base := new(mockLogger)
	ctx := context.Background()
	execedDebug := false
	execedInfo := false
	execedWarn := false
	execedError := false

	logger := newCallbackLogger(
		t,
		base,
		decorator.WithDebugCallback(func(context.Context, string, ...any) {
			if len(base.debugIn) != 1 {
				t.Errorf("expected Debug to be called once, got %d calls", len(base.debugIn))
			}

			execedDebug = true
		}),
		decorator.WithInfoCallback(func(context.Context, string, ...any) {
			if len(base.infoIn) != 1 {
				t.Errorf("expected Info to be called once, got %d calls", len(base.infoIn))
			}

			execedInfo = true
		}),
		decorator.WithWarnCallback(func(context.Context, string, ...any) {
			if len(base.warnIn) != 1 {
				t.Errorf("expected Warn to be called once, got %d calls", len(base.warnIn))
			}

			execedWarn = true
		}),
		decorator.WithErrorCallback(func(context.Context, error, ...any) {
			if len(base.errorIn) != 1 {
				t.Errorf("expected Error to be called once, got %d calls", len(base.errorIn))
			}

			execedError = true
		}),
	)

	logger.Debug(ctx, "debug message")
	logger.Info(ctx, "info message")
	logger.Warn(ctx, "warn message")
	logger.Error(ctx, errors.New("test error"))

	if !execedDebug {
		t.Errorf("expected debug callback to be executed")
	}

	if !execedInfo {
		t.Errorf("expected info callback to be executed")
	}

	if !execedWarn {
		t.Errorf("expected warn callback to be executed")
	}

	if !execedError {
		t.Errorf("expected error callback to be executed")
	}
}

// TestCallbackLogger_CopiesArgs tests that the callback receives a copy of
// the args slice, so mutations in the callback do not affect the logger.
func TestCallbackLogger_CopiesArgs(t *testing.T) {
	t.Parallel()

	base := new(mockLogger)

	logger := newCallbackLogger(
		t,
		base,
		decorator.WithInfoCallback(func(_ context.Context, _ string, args ...any) {
			args[0] = "from-callback"
		}),
	)

	logger.Info(context.Background(), "info message", "original")

	if len(base.infoIn) != 1 || len(base.infoIn[0].args) != 1 {
		t.Fatalf("expected 1 info call with 1 argument, got %+v", base.infoIn)
	}

	if base.infoIn[0].args[0] != "original" {
		t.Errorf("expected logger args to stay original, got %v", base.infoIn[0].args[0])
	}
}

// TestCallbackLogger_PanicInLoggerSkipsCallback tests that a panic in the
// underlying logger prevents the callback from running.
func TestCallbackLogger_PanicInLoggerSkipsCallback(t *testing.T) {
	t.Parallel()

	called := false
	inner := &hookLogger{
		infoFn: func() {
			panic("underlying logger panic")
		},
	}

	logger := newCallbackLogger(
		t,
		inner,
		decorator.WithInfoCallback(func(context.Context, string, ...any) {
			called = true
		}),
	)

	defer func() {
		recovered := recover()
		if recovered == nil {
			t.Fatal("expected panic from underlying logger")
		}

		if called {
			t.Errorf("expected callback not to run after logger panic")
		}
	}()

	logger.Info(context.Background(), "info message")
}
