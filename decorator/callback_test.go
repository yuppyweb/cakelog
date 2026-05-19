package decorator_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/yuppyweb/cakelog/decorator"
)

// TestWithDebugCallback_NilOptions checks that WithDebugCallback returns
// ErrNilCallbackOpts error when passed nil options.
func TestWithDebugCallback_NilOptions(t *testing.T) {
	t.Parallel()

	opt := decorator.WithDebugCallback(func(ctx context.Context) {})

	err := opt(nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, decorator.ErrNilCallbackOpts) {
		t.Errorf("expected error to be ErrNilCallbackOpts, got %v", err)
	}
}

// TestWithDebugCallback_NilFunc checks that WithDebugCallback returns
// ErrNilDebugFunc error when passed nil function.
func TestWithDebugCallback_NilFunc(t *testing.T) {
	t.Parallel()

	opt := decorator.WithDebugCallback(nil)

	err := opt(decorator.DefaultCallbackOptions())
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, decorator.ErrNilDebugFunc) {
		t.Errorf("expected error to be ErrNilDebugFunc, got %v", err)
	}
}

// TestWithDebugCallback_Success checks successful creation of WithDebugCallback
// with valid parameters.
func TestWithDebugCallback_Success(t *testing.T) {
	t.Parallel()

	opt := decorator.WithDebugCallback(func(ctx context.Context) {})

	err := opt(decorator.DefaultCallbackOptions())
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

// TestWithInfoCallback_NilOptions checks that WithInfoCallback returns
// ErrNilCallbackOpts error when passed nil options.
func TestWithInfoCallback_NilOptions(t *testing.T) {
	t.Parallel()

	opt := decorator.WithInfoCallback(func(ctx context.Context) {})

	err := opt(nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, decorator.ErrNilCallbackOpts) {
		t.Errorf("expected error to be ErrNilCallbackOpts, got %v", err)
	}
}

// TestWithInfoCallback_NilFunc checks that WithInfoCallback returns
// ErrNilInfoFunc error when passed nil function.
func TestWithInfoCallback_NilFunc(t *testing.T) {
	t.Parallel()

	opt := decorator.WithInfoCallback(nil)

	err := opt(decorator.DefaultCallbackOptions())
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, decorator.ErrNilInfoFunc) {
		t.Errorf("expected error to be ErrNilInfoFunc, got %v", err)
	}
}

// TestWithInfoCallback_Success checks successful creation of WithInfoCallback
// with valid parameters.
func TestWithInfoCallback_Success(t *testing.T) {
	t.Parallel()

	opt := decorator.WithInfoCallback(func(ctx context.Context) {})

	err := opt(decorator.DefaultCallbackOptions())
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

// TestWithWarnCallback_NilOptions checks that WithWarnCallback returns
// ErrNilCallbackOpts error when passed nil options.
func TestWithWarnCallback_NilOptions(t *testing.T) {
	t.Parallel()

	opt := decorator.WithWarnCallback(func(ctx context.Context) {})

	err := opt(nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, decorator.ErrNilCallbackOpts) {
		t.Errorf("expected error to be ErrNilCallbackOpts, got %v", err)
	}
}

// TestWithWarnCallback_NilFunc checks that WithWarnCallback returns
// ErrNilWarnFunc error when passed nil function.
func TestWithWarnCallback_NilFunc(t *testing.T) {
	t.Parallel()

	opt := decorator.WithWarnCallback(nil)

	err := opt(decorator.DefaultCallbackOptions())
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, decorator.ErrNilWarnFunc) {
		t.Errorf("expected error to be ErrNilWarnFunc, got %v", err)
	}
}

// TestWithWarnCallback_Success checks successful creation of WithWarnCallback
// with valid parameters.
func TestWithWarnCallback_Success(t *testing.T) {
	t.Parallel()

	opt := decorator.WithWarnCallback(func(ctx context.Context) {})

	err := opt(decorator.DefaultCallbackOptions())
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

// TestWithErrorCallback_NilOptions checks that WithErrorCallback returns
// ErrNilCallbackOpts error when passed nil options.
func TestWithErrorCallback_NilOptions(t *testing.T) {
	t.Parallel()

	opt := decorator.WithErrorCallback(func(ctx context.Context) {})

	err := opt(nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, decorator.ErrNilCallbackOpts) {
		t.Errorf("expected error to be ErrNilCallbackOpts, got %v", err)
	}
}

// TestWithErrorCallback_NilFunc checks that WithErrorCallback returns
// ErrNilErrorFunc error when passed nil function.
func TestWithErrorCallback_NilFunc(t *testing.T) {
	t.Parallel()

	opt := decorator.WithErrorCallback(nil)

	err := opt(decorator.DefaultCallbackOptions())
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, decorator.ErrNilErrorFunc) {
		t.Errorf("expected error to be ErrNilErrorFunc, got %v", err)
	}
}

// TestWithErrorCallback_Success checks successful creation of WithErrorCallback
// with valid parameters.
func TestWithErrorCallback_Success(t *testing.T) {
	t.Parallel()

	opt := decorator.WithErrorCallback(func(ctx context.Context) {})

	err := opt(decorator.DefaultCallbackOptions())
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

// TestNewCallbackLogger_NilLogger checks that NewCallbackLogger returns
// ErrNilCallbackLogger error when passed nil logger.
func TestNewCallbackLogger_NilLogger(t *testing.T) {
	t.Parallel()

	_, err := decorator.NewCallbackLogger(nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, decorator.ErrNilCallbackLogger) {
		t.Errorf("expected error to be ErrNilCallbackLogger, got %v", err)
	}
}

// TestNewCallbackLogger_NilOption checks that NewCallbackLogger returns
// ErrNilCallbackOpt error when passed nil option.
func TestNewCallbackLogger_NilOption(t *testing.T) {
	t.Parallel()

	_, err := decorator.NewCallbackLogger(new(mockLogger), nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, decorator.ErrNilCallbackOpt) {
		t.Errorf("expected error to be ErrNilCallbackOpt, got %v", err)
	}
}

// TestNewCallbackLogger_NilDebugCallback checks that NewCallbackLogger returns
// error with correct message when passed nil debug callback.
func TestNewCallbackLogger_NilDebugCallback(t *testing.T) {
	t.Parallel()

	opt := decorator.WithDebugCallback(nil)

	_, err := decorator.NewCallbackLogger(new(mockLogger), opt)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, decorator.ErrNilDebugFunc) {
		t.Errorf("expected error to be ErrNilDebugFunc, got %v", err)
	}

	if !strings.Contains(err.Error(), "failed to apply option:") {
		t.Errorf(
			"error message does not contain expected text:\nGot:  %s\nWant to contain: %s",
			err.Error(), "failed to apply option:")
	}
}

// TestNewCallbackLogger_NilInfoCallback checks that NewCallbackLogger returns
// error with correct message when passed nil info callback.
func TestNewCallbackLogger_NilInfoCallback(t *testing.T) {
	t.Parallel()

	opt := decorator.WithInfoCallback(nil)

	_, err := decorator.NewCallbackLogger(new(mockLogger), opt)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, decorator.ErrNilInfoFunc) {
		t.Errorf("expected error to be ErrNilInfoFunc, got %v", err)
	}

	if !strings.Contains(err.Error(), "failed to apply option:") {
		t.Errorf(
			"error message does not contain expected text:\nGot:  %s\nWant to contain: %s",
			err.Error(), "failed to apply option:")
	}
}

// TestNewCallbackLogger_NilWarnCallback checks that NewCallbackLogger returns
// error with correct message when passed nil warn callback.
func TestNewCallbackLogger_NilWarnCallback(t *testing.T) {
	t.Parallel()

	opt := decorator.WithWarnCallback(nil)

	_, err := decorator.NewCallbackLogger(new(mockLogger), opt)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, decorator.ErrNilWarnFunc) {
		t.Errorf("expected error to be ErrNilWarnFunc, got %v", err)
	}

	if !strings.Contains(err.Error(), "failed to apply option:") {
		t.Errorf(
			"error message does not contain expected text:\nGot:  %s\nWant to contain: %s",
			err.Error(), "failed to apply option:")
	}
}

// TestNewCallbackLogger_NilErrorCallback checks that NewCallbackLogger returns
// error with correct message when passed nil error callback.
func TestNewCallbackLogger_NilErrorCallback(t *testing.T) {
	t.Parallel()

	opt := decorator.WithErrorCallback(nil)

	_, err := decorator.NewCallbackLogger(new(mockLogger), opt)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, decorator.ErrNilErrorFunc) {
		t.Errorf("expected error to be ErrNilErrorFunc, got %v", err)
	}

	if !strings.Contains(err.Error(), "failed to apply option:") {
		t.Errorf(
			"error message does not contain expected text:\nGot:  %s\nWant to contain: %s",
			err.Error(), "failed to apply option:")
	}
}

// TestCallbackLogger_Debug checks that CallbackLogger correctly invokes
// debug callback and passes parameters to the underlying logger.
func TestCallbackLogger_Debug(t *testing.T) {
	t.Parallel()

	type ctxKey struct{}

	mockLogger := new(mockLogger)
	expectedCtx := context.WithValue(context.Background(), ctxKey{}, "test debug")
	execed := false

	opt := decorator.WithDebugCallback(func(ctx context.Context) {
		if ctx != expectedCtx {
			t.Errorf("expected context to be %v, got %v", expectedCtx, ctx)
		}

		execed = true
	})

	logger, err := decorator.NewCallbackLogger(mockLogger, opt)
	if err != nil {
		t.Fatalf("unexpected error when creating CallbackLogger: %v", err)
	}

	logger.Debug(expectedCtx, "debug message", "debug", 42)

	if len(mockLogger.debugIn) != 1 {
		t.Fatalf("Expected Debug to be called once, got %d calls", len(mockLogger.debugIn))
	}

	if mockLogger.debugIn[0].ctx != expectedCtx {
		t.Errorf(
			"Expected Debug context to be %v, got %v",
			expectedCtx,
			mockLogger.debugIn[0].ctx,
		)
	}

	if mockLogger.debugIn[0].msg != "debug message" {
		t.Errorf(
			"Expected Debug message to be 'debug message', got '%s'",
			mockLogger.debugIn[0].msg,
		)
	}

	if len(mockLogger.debugIn[0].args) != 2 {
		t.Fatalf(
			"Expected Debug to be called with 2 arguments, got %d",
			len(mockLogger.debugIn[0].args),
		)
	}

	if mockLogger.debugIn[0].args[0] != "debug" {
		t.Errorf(
			"Expected first argument to be 'debug', got '%v'",
			mockLogger.debugIn[0].args[0],
		)
	}

	if mockLogger.debugIn[0].args[1] != 42 {
		t.Errorf(
			"Expected second argument to be 42, got '%v'",
			mockLogger.debugIn[0].args[1],
		)
	}

	if !execed {
		t.Errorf("expected debug callback to be executed")
	}
}

// TestCallbackLogger_Info checks that CallbackLogger correctly invokes
// info callback and passes parameters to the underlying logger.
func TestCallbackLogger_Info(t *testing.T) {
	t.Parallel()

	type ctxKey struct{}

	mockLogger := new(mockLogger)
	expectedCtx := context.WithValue(context.Background(), ctxKey{}, "test info")
	execed := false

	opt := decorator.WithInfoCallback(func(ctx context.Context) {
		if ctx != expectedCtx {
			t.Errorf("expected context to be %v, got %v", expectedCtx, ctx)
		}

		execed = true
	})

	logger, err := decorator.NewCallbackLogger(mockLogger, opt)
	if err != nil {
		t.Fatalf("unexpected error when creating CallbackLogger: %v", err)
	}

	logger.Info(expectedCtx, "info message", "info", 75)

	if len(mockLogger.infoIn) != 1 {
		t.Fatalf("Expected Info to be called once, got %d calls", len(mockLogger.infoIn))
	}

	if mockLogger.infoIn[0].ctx != expectedCtx {
		t.Errorf(
			"Expected Info context to be %v, got %v",
			expectedCtx,
			mockLogger.infoIn[0].ctx,
		)
	}

	if mockLogger.infoIn[0].msg != "info message" {
		t.Errorf(
			"Expected Info message to be 'info message', got '%s'",
			mockLogger.infoIn[0].msg,
		)
	}

	if len(mockLogger.infoIn[0].args) != 2 {
		t.Fatalf(
			"Expected Info to be called with 2 arguments, got %d",
			len(mockLogger.infoIn[0].args),
		)
	}

	if mockLogger.infoIn[0].args[0] != "info" {
		t.Errorf(
			"Expected first argument to be 'info', got '%v'",
			mockLogger.infoIn[0].args[0],
		)
	}

	if mockLogger.infoIn[0].args[1] != 75 {
		t.Errorf(
			"Expected second argument to be 75, got '%v'",
			mockLogger.infoIn[0].args[1],
		)
	}

	if !execed {
		t.Errorf("expected info callback to be executed")
	}
}

// TestCallbackLogger_Warn checks that CallbackLogger correctly invokes
// warn callback and passes parameters to the underlying logger.
func TestCallbackLogger_Warn(t *testing.T) {
	t.Parallel()

	type ctxKey struct{}

	mockLogger := new(mockLogger)
	expectedCtx := context.WithValue(context.Background(), ctxKey{}, "test warn")
	execed := false

	opt := decorator.WithWarnCallback(func(ctx context.Context) {
		if ctx != expectedCtx {
			t.Errorf("expected context to be %v, got %v", expectedCtx, ctx)
		}

		execed = true
	})

	logger, err := decorator.NewCallbackLogger(mockLogger, opt)
	if err != nil {
		t.Fatalf("unexpected error when creating CallbackLogger: %v", err)
	}

	logger.Warn(expectedCtx, "warn message", "warn", 80)

	if len(mockLogger.warnIn) != 1 {
		t.Fatalf("Expected Warn to be called once, got %d calls", len(mockLogger.warnIn))
	}

	if mockLogger.warnIn[0].ctx != expectedCtx {
		t.Errorf(
			"Expected Warn context to be %v, got %v",
			expectedCtx,
			mockLogger.warnIn[0].ctx,
		)
	}

	if mockLogger.warnIn[0].msg != "warn message" {
		t.Errorf(
			"Expected Warn message to be 'warn message', got '%s'",
			mockLogger.warnIn[0].msg,
		)
	}

	if len(mockLogger.warnIn[0].args) != 2 {
		t.Fatalf(
			"Expected Warn to be called with 2 arguments, got %d",
			len(mockLogger.warnIn[0].args),
		)
	}

	if mockLogger.warnIn[0].args[0] != "warn" {
		t.Errorf(
			"Expected first argument to be 'warn', got '%v'",
			mockLogger.warnIn[0].args[0],
		)
	}

	if mockLogger.warnIn[0].args[1] != 80 {
		t.Errorf(
			"Expected second argument to be 80, got '%v'",
			mockLogger.warnIn[0].args[1],
		)
	}

	if !execed {
		t.Errorf("expected warn callback to be executed")
	}
}

// TestCallbackLogger_Error checks that CallbackLogger correctly invokes
// error callback and passes parameters to the underlying logger.
func TestCallbackLogger_Error(t *testing.T) {
	t.Parallel()

	type ctxKey struct{}

	mockLogger := new(mockLogger)
	expectedCtx := context.WithValue(context.Background(), ctxKey{}, "test error")
	expectedErr := errors.New("test error")
	execed := false

	opt := decorator.WithErrorCallback(func(ctx context.Context) {
		if ctx != expectedCtx {
			t.Errorf("expected context to be %v, got %v", expectedCtx, ctx)
		}

		execed = true
	})

	logger, err := decorator.NewCallbackLogger(mockLogger, opt)
	if err != nil {
		t.Fatalf("unexpected error when creating CallbackLogger: %v", err)
	}

	logger.Error(expectedCtx, expectedErr, "error", 99)

	if len(mockLogger.errorIn) != 1 {
		t.Fatalf("Expected Error to be called once, got %d calls", len(mockLogger.errorIn))
	}

	if mockLogger.errorIn[0].ctx != expectedCtx {
		t.Errorf(
			"Expected Error context to be %v, got %v",
			expectedCtx,
			mockLogger.errorIn[0].ctx,
		)
	}

	if !errors.Is(mockLogger.errorIn[0].err, expectedErr) {
		t.Errorf(
			"Expected Error to be %v, got %v",
			expectedErr,
			mockLogger.errorIn[0].err,
		)
	}

	if len(mockLogger.errorIn[0].args) != 2 {
		t.Fatalf(
			"Expected Error to be called with 2 arguments, got %d",
			len(mockLogger.errorIn[0].args),
		)
	}

	if mockLogger.errorIn[0].args[0] != "error" {
		t.Errorf(
			"Expected first argument to be 'error', got '%v'",
			mockLogger.errorIn[0].args[0],
		)
	}

	if mockLogger.errorIn[0].args[1] != 99 {
		t.Errorf(
			"Expected second argument to be 99, got '%v'",
			mockLogger.errorIn[0].args[1],
		)
	}

	if !execed {
		t.Errorf("expected error callback to be executed")
	}
}

// TestCallbackLogger_AllCallbacks checks that CallbackLogger correctly invokes
// all four callback functions when logging messages at different levels.
func TestCallbackLogger_AllCallbacks(t *testing.T) {
	t.Parallel()

	mockLogger := new(mockLogger)
	ctx := context.Background()
	execedDebug := false
	execedInfo := false
	execedWarn := false
	execedError := false

	opts := []decorator.CallbackOption{
		decorator.WithDebugCallback(func(ctx context.Context) {
			execedDebug = true
		}),
		decorator.WithInfoCallback(func(ctx context.Context) {
			execedInfo = true
		}),
		decorator.WithWarnCallback(func(ctx context.Context) {
			execedWarn = true
		}),
		decorator.WithErrorCallback(func(ctx context.Context) {
			execedError = true
		}),
	}

	logger, err := decorator.NewCallbackLogger(mockLogger, opts...)
	if err != nil {
		t.Fatalf("unexpected error when creating CallbackLogger: %v", err)
	}

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

// TestCallbackLogger_MultipleCallbacks checks that CallbackLogger correctly invokes
// callbacks multiple times when logging multiple messages.
func TestCallbackLogger_MultipleCallbacks(t *testing.T) {
	t.Parallel()

	mockLogger := new(mockLogger)
	ctx := context.Background()
	countDebug := 0
	countInfo := 0
	countWarn := 0
	countError := 0

	opts := []decorator.CallbackOption{
		decorator.WithDebugCallback(func(ctx context.Context) {
			countDebug++
		}),
		decorator.WithInfoCallback(func(ctx context.Context) {
			countInfo++
		}),
		decorator.WithWarnCallback(func(ctx context.Context) {
			countWarn++
		}),
		decorator.WithErrorCallback(func(ctx context.Context) {
			countError++
		}),
	}

	logger, err := decorator.NewCallbackLogger(mockLogger, opts...)
	if err != nil {
		t.Fatalf("unexpected error when creating CallbackLogger: %v", err)
	}

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

// TestCallbackLogger_CallbackAfterLog checks that callback is invoked after
// the message has been passed to the underlying logger.
func TestCallbackLogger_CallbackAfterLog(t *testing.T) {
	t.Parallel()

	mockLogger := new(mockLogger)
	ctx := context.Background()
	execedDebug := false
	execedInfo := false
	execedWarn := false
	execedError := false

	opts := []decorator.CallbackOption{
		decorator.WithDebugCallback(func(ctx context.Context) {
			if len(mockLogger.debugIn) != 1 {
				t.Errorf("Expected Debug to be called once, got %d calls", len(mockLogger.debugIn))
			}

			execedDebug = true
		}),
		decorator.WithInfoCallback(func(ctx context.Context) {
			if len(mockLogger.infoIn) != 1 {
				t.Errorf("Expected Info to be called once, got %d calls", len(mockLogger.infoIn))
			}

			execedInfo = true
		}),
		decorator.WithWarnCallback(func(ctx context.Context) {
			if len(mockLogger.warnIn) != 1 {
				t.Errorf("Expected Warn to be called once, got %d calls", len(mockLogger.warnIn))
			}

			execedWarn = true
		}),
		decorator.WithErrorCallback(func(ctx context.Context) {
			if len(mockLogger.errorIn) != 1 {
				t.Errorf("Expected Error to be called once, got %d calls", len(mockLogger.errorIn))
			}

			execedError = true
		}),
	}

	logger, err := decorator.NewCallbackLogger(mockLogger, opts...)
	if err != nil {
		t.Fatalf("unexpected error when creating CallbackLogger: %v", err)
	}

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
