package decorator_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/yuppyweb/cakelog"
	"github.com/yuppyweb/cakelog/decorator"
)

// hookLogger invokes optional hooks on each log call.
type hookLogger struct {
	debugFn func()
	infoFn  func()
	warnFn  func()
	errorFn func()
}

// Debug runs debugFn when it is set.
func (hl *hookLogger) Debug(_ context.Context, _ string, _ ...any) {
	if hl.debugFn != nil {
		hl.debugFn()
	}
}

// Info runs infoFn when it is set.
func (hl *hookLogger) Info(_ context.Context, _ string, _ ...any) {
	if hl.infoFn != nil {
		hl.infoFn()
	}
}

// Warn runs warnFn when it is set.
func (hl *hookLogger) Warn(_ context.Context, _ string, _ ...any) {
	if hl.warnFn != nil {
		hl.warnFn()
	}
}

// Error runs errorFn when it is set.
func (hl *hookLogger) Error(_ context.Context, _ error, _ ...any) {
	if hl.errorFn != nil {
		hl.errorFn()
	}
}

// Ensure hookLogger implements the cakelog.Logger interface.
var _ cakelog.Logger = (*hookLogger)(nil)

func newMuxLogger(t *testing.T, logs ...cakelog.Logger) cakelog.Logger {
	t.Helper()

	logger, err := decorator.NewMux(logs...)
	if err != nil {
		t.Fatalf("failed to create mux logger: %v", err)
	}

	return logger
}

func assertLoggedMsg(
	t *testing.T,
	got []mockMsgArgs,
	ctx context.Context,
	msg string,
	args ...any,
) {
	t.Helper()

	if len(got) != 1 {
		t.Fatalf("expected 1 log call, got %d", len(got))
	}

	if got[0].ctx != ctx {
		t.Errorf("expected context %v, got %v", ctx, got[0].ctx)
	}

	if got[0].msg != msg {
		t.Errorf("expected message %q, got %q", msg, got[0].msg)
	}

	if len(got[0].args) != len(args) {
		t.Fatalf("expected %d arguments, got %d", len(args), len(got[0].args))
	}

	for i, arg := range args {
		if got[0].args[i] != arg {
			t.Errorf("expected argument %d to be %v, got %v", i, arg, got[0].args[i])
		}
	}
}

func assertLoggedErr(
	t *testing.T,
	got []mockErrArgs,
	ctx context.Context,
	wantErr error,
	args ...any,
) {
	t.Helper()

	if len(got) != 1 {
		t.Fatalf("expected 1 log call, got %d", len(got))
	}

	if got[0].ctx != ctx {
		t.Errorf("expected context %v, got %v", ctx, got[0].ctx)
	}

	if wantErr == nil {
		if got[0].err != nil {
			t.Errorf("expected nil error, got %v", got[0].err)
		}
	} else if !errors.Is(got[0].err, wantErr) {
		t.Errorf("expected error %v, got %v", wantErr, got[0].err)
	}

	if len(got[0].args) != len(args) {
		t.Fatalf("expected %d arguments, got %d", len(args), len(got[0].args))
	}

	for i, arg := range args {
		if got[0].args[i] != arg {
			t.Errorf("expected argument %d to be %v, got %v", i, arg, got[0].args[i])
		}
	}
}

// TestNewMux_NoLogger tests that NewMux returns a wrapped
// ErrNoLoggerProvided when no loggers are provided.
func TestNewMux_NoLogger(t *testing.T) {
	t.Parallel()

	_, err := decorator.NewMux()
	if err == nil {
		t.Fatal("expected error when providing no loggers, got nil")
	}

	if !errors.Is(err, decorator.ErrNoLoggerProvided) {
		t.Errorf(
			"unexpected error when creating mux logger with no loggers:\nGot:  %v\nWant: %v",
			err,
			decorator.ErrNoLoggerProvided,
		)
	}

	if !strings.Contains(err.Error(), "mux logger:") {
		t.Errorf(
			"error message does not contain expected text:\nGot:  %s\nWant to contain: %s",
			err.Error(),
			"mux logger:",
		)
	}
}

// TestNewMux_NilLogger tests that NewMux returns a wrapped ErrNilLogger
// when the first logger is nil.
func TestNewMux_NilLogger(t *testing.T) {
	t.Parallel()

	_, err := decorator.NewMux(nil)
	if err == nil {
		t.Fatal("expected error when providing a nil logger, got nil")
	}

	if !errors.Is(err, decorator.ErrNilLogger) {
		t.Errorf(
			"unexpected error when creating mux logger with nil logger:\nGot:  %v\nWant: %v",
			err,
			decorator.ErrNilLogger,
		)
	}

	if !strings.Contains(err.Error(), "mux logger: 1:") {
		t.Errorf(
			"error message does not contain expected text:\nGot:  %s\nWant to contain: %s",
			err.Error(),
			"mux logger: 1:",
		)
	}
}

// TestNewMux_TypedNilLogger tests that NewMux returns a wrapped ErrNilLogger
// when a typed nil logger is provided.
func TestNewMux_TypedNilLogger(t *testing.T) {
	t.Parallel()

	var typedNil *mockLogger

	_, err := decorator.NewMux(typedNil)
	if err == nil {
		t.Fatal("expected error when providing a typed nil logger, got nil")
	}

	if !errors.Is(err, decorator.ErrNilLogger) {
		t.Errorf(
			"unexpected error when creating mux logger with typed nil logger:\nGot:  %v\nWant: %v",
			err,
			decorator.ErrNilLogger,
		)
	}

	if !strings.Contains(err.Error(), "mux logger: 1:") {
		t.Errorf(
			"error message does not contain expected text:\nGot:  %s\nWant to contain: %s",
			err.Error(),
			"mux logger: 1:",
		)
	}
}

// TestNewMux_NilLoggerPosition tests that NewMux reports a 1-based position
// for a nil logger that is not the first in the list.
func TestNewMux_NilLoggerPosition(t *testing.T) {
	t.Parallel()

	_, err := decorator.NewMux(new(mockLogger), nil, new(mockLogger))
	if err == nil {
		t.Fatal("expected error when providing a nil logger, got nil")
	}

	if !errors.Is(err, decorator.ErrNilLogger) {
		t.Errorf(
			"unexpected error when creating mux logger with nil logger:\nGot:  %v\nWant: %v",
			err,
			decorator.ErrNilLogger,
		)
	}

	if !strings.Contains(err.Error(), "mux logger: 2:") {
		t.Errorf(
			"error message does not contain expected text:\nGot:  %s\nWant to contain: %s",
			err.Error(),
			"mux logger: 2:",
		)
	}
}

// TestNewMux_CopiesLoggers tests that NewMux copies the loggers slice so
// later mutations of the caller's slice do not affect the mux.
func TestNewMux_CopiesLoggers(t *testing.T) {
	t.Parallel()

	original := new(mockLogger)
	replacement := new(mockLogger)
	logs := []cakelog.Logger{original}

	logger, err := decorator.NewMux(logs...)
	if err != nil {
		t.Fatalf("failed to create mux logger: %v", err)
	}

	logs[0] = replacement

	logger.Info(context.Background(), "message")

	if len(original.infoIn) != 1 {
		t.Fatalf("expected original logger Info to be called once, got %d", len(original.infoIn))
	}

	if len(replacement.infoIn) != 0 {
		t.Errorf(
			"expected replacement logger not to receive Info, got %d calls",
			len(replacement.infoIn),
		)
	}
}

// TestMuxLogger_ForwardsToAll tests that Debug, Info, Warn, and Error
// forward the same context, message, error, and arguments to every logger.
func TestMuxLogger_ForwardsToAll(t *testing.T) {
	t.Parallel()

	type ctxKey struct{}

	first := new(mockLogger)
	second := new(mockLogger)
	logger := newMuxLogger(t, first, second)

	expectedCtx := context.WithValue(context.Background(), ctxKey{}, "mux")
	expectedErr := errors.New("error message")

	logger.Debug(expectedCtx, "debug message", "debug", 42)
	logger.Info(expectedCtx, "info message", "info", 75)
	logger.Warn(expectedCtx, "warn message", "warn", 80)
	logger.Error(expectedCtx, expectedErr, "error", 99)

	for _, log := range []*mockLogger{first, second} {
		assertLoggedMsg(t, log.debugIn, expectedCtx, "debug message", "debug", 42)
		assertLoggedMsg(t, log.infoIn, expectedCtx, "info message", "info", 75)
		assertLoggedMsg(t, log.warnIn, expectedCtx, "warn message", "warn", 80)
		assertLoggedErr(t, log.errorIn, expectedCtx, expectedErr, "error", 99)
	}
}

// TestMuxLogger_EmptyArgs tests that all methods forward calls with no
// extra arguments, including a nil error.
func TestMuxLogger_EmptyArgs(t *testing.T) {
	t.Parallel()

	first := new(mockLogger)
	second := new(mockLogger)
	logger := newMuxLogger(t, first, second)

	ctx := context.Background()

	logger.Debug(ctx, "debug message")
	logger.Info(ctx, "info message")
	logger.Warn(ctx, "warn message")
	logger.Error(ctx, nil)

	for _, log := range []*mockLogger{first, second} {
		assertLoggedMsg(t, log.debugIn, ctx, "debug message")
		assertLoggedMsg(t, log.infoIn, ctx, "info message")
		assertLoggedMsg(t, log.warnIn, ctx, "warn message")
		assertLoggedErr(t, log.errorIn, ctx, nil)
	}
}

// TestMuxLogger_DuplicateLogger tests that passing the same logger more
// than once forwards each call to it on every occurrence.
func TestMuxLogger_DuplicateLogger(t *testing.T) {
	t.Parallel()

	same := new(mockLogger)
	logger := newMuxLogger(t, same, same)

	logger.Info(context.Background(), "message", "key", "value")

	if len(same.infoIn) != 2 {
		t.Fatalf("expected Info to be called twice, got %d", len(same.infoIn))
	}

	if same.infoIn[0].msg != "message" || same.infoIn[1].msg != "message" {
		t.Errorf(
			"expected both Info calls to use 'message', got %q and %q",
			same.infoIn[0].msg,
			same.infoIn[1].msg,
		)
	}
}

// TestMuxLogger_ForwardsInOrder tests that each log call is forwarded to
// the underlying loggers in the order they were provided.
func TestMuxLogger_ForwardsInOrder(t *testing.T) {
	t.Parallel()

	var order []string

	first := &hookLogger{
		debugFn: func() {
			order = append(order, "first:debug")
		},
		infoFn: func() {
			order = append(order, "first:info")
		},
		warnFn: func() {
			order = append(order, "first:warn")
		},
		errorFn: func() {
			order = append(order, "first:error")
		},
	}
	second := &hookLogger{
		debugFn: func() {
			order = append(order, "second:debug")
		},
		infoFn: func() {
			order = append(order, "second:info")
		},
		warnFn: func() {
			order = append(order, "second:warn")
		},
		errorFn: func() {
			order = append(order, "second:error")
		},
	}

	logger := newMuxLogger(t, first, second)

	logger.Debug(context.Background(), "debug")
	logger.Info(context.Background(), "info")
	logger.Warn(context.Background(), "warn")
	logger.Error(context.Background(), errors.New("error"))

	want := []string{
		"first:debug",
		"second:debug",
		"first:info",
		"second:info",
		"first:warn",
		"second:warn",
		"first:error",
		"second:error",
	}

	if len(order) != len(want) {
		t.Fatalf("expected %d calls, got %d: %v", len(want), len(order), order)
	}

	for i, name := range want {
		if order[i] != name {
			t.Errorf("expected call %d to be %q, got %q", i, name, order[i])
		}
	}
}

// TestMuxLogger_CopiesArgs tests that each logger receives its own copy of
// the args slice, so later mutations of the caller's slice or of one
// logger's args do not affect the other logger.
func TestMuxLogger_CopiesArgs(t *testing.T) {
	t.Parallel()

	first := new(mockLogger)
	second := new(mockLogger)
	logger := newMuxLogger(t, first, second)

	args := []any{"key", "value"}
	logger.Info(context.Background(), "message", args...)

	args[1] = "mutated"

	if len(first.infoIn) != 1 || len(second.infoIn) != 1 {
		t.Fatalf(
			"expected both loggers to receive one Info call, got %d and %d",
			len(first.infoIn),
			len(second.infoIn),
		)
	}

	if first.infoIn[0].args[1] != "value" {
		t.Errorf(
			"expected first logger to keep original args, got %v",
			first.infoIn[0].args[1],
		)
	}

	if second.infoIn[0].args[1] != "value" {
		t.Errorf(
			"expected second logger to keep original args, got %v",
			second.infoIn[0].args[1],
		)
	}

	first.infoIn[0].args[1] = "from-first"

	if second.infoIn[0].args[1] != "value" {
		t.Errorf(
			"expected second logger args to be isolated, got %v",
			second.infoIn[0].args[1],
		)
	}
}

// TestMuxLogger_PanicStopsRemaining tests that a panic in one logger is
// not recovered and stops remaining loggers.
func TestMuxLogger_PanicStopsRemaining(t *testing.T) {
	t.Parallel()

	second := new(mockLogger)
	first := &hookLogger{
		infoFn: func() {
			panic("mux logger panic")
		},
	}

	logger := newMuxLogger(t, first, second)

	defer func() {
		got := recover()
		if got == nil {
			return
		}

		if got != "mux logger panic" {
			t.Errorf("expected panic %q, got %v", "mux logger panic", got)
		}

		if len(second.infoIn) != 0 {
			t.Errorf(
				"expected second logger not to be called after panic, got %d calls",
				len(second.infoIn),
			)
		}
	}()

	logger.Info(context.Background(), "message")

	t.Fatal("expected panic from first logger")
}
