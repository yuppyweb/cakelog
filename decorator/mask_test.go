package decorator_test

import (
	"context"
	"errors"
	"slices"
	"strings"
	"testing"

	"github.com/yuppyweb/cakelog"
	"github.com/yuppyweb/cakelog/decorator"
)

// mockMasker is a mock implementation of the Masker interface used for testing.
type mockMasker struct {
	msg           string
	maskMsg       string
	err           error
	maskErr       error
	args          []any
	maskArgs      []any
	returnNilArgs bool
}

// MaskMessage masks a log message and returns the masked version.
func (m *mockMasker) MaskMessage(msg string) string {
	m.msg = msg

	return m.maskMsg
}

// MaskError masks an error and returns the masked version.
func (m *mockMasker) MaskError(err error) error {
	m.err = err

	return m.maskErr
}

// MaskArguments masks log arguments and returns the masked version.
func (m *mockMasker) MaskArguments(args ...any) []any {
	m.args = slices.Clone(args)

	if m.returnNilArgs {
		return nil
	}

	if m.maskArgs == nil {
		return m.args
	}

	return slices.Clone(m.maskArgs)
}

var _ decorator.Masker = (*mockMasker)(nil)

// replaceArgMasker replaces the first argument it receives. It mutates the
// slice passed to MaskArguments so tests can assert that NewMask copied
// the caller's args slice.
type replaceArgMasker struct {
	replaced any
}

func (m replaceArgMasker) MaskMessage(msg string) string {
	return msg
}

func (m replaceArgMasker) MaskError(err error) error {
	return err
}

func (m replaceArgMasker) MaskArguments(args ...any) []any {
	if len(args) > 0 {
		args[0] = m.replaced
	}

	return args
}

var _ decorator.Masker = replaceArgMasker{}

func newMaskLogger(t *testing.T, log cakelog.Logger, maskers ...decorator.Masker) cakelog.Logger {
	t.Helper()

	logger, err := decorator.NewMask(log, maskers...)
	if err != nil {
		t.Fatalf("failed to create mask logger: %v", err)
	}

	return logger
}

// TestNewMask_NilLogger tests that NewMask returns a wrapped ErrNilLogger
// when the logger is nil.
func TestNewMask_NilLogger(t *testing.T) {
	t.Parallel()

	_, err := decorator.NewMask(nil)
	if err == nil {
		t.Fatal("expected error when providing a nil logger, got nil")
	}

	if !errors.Is(err, decorator.ErrNilLogger) {
		t.Errorf(
			"unexpected error when creating mask logger with nil logger:\nGot:  %v\nWant: %v",
			err,
			decorator.ErrNilLogger,
		)
	}

	if !strings.Contains(err.Error(), "mask logger:") {
		t.Errorf(
			"error message does not contain expected text:\nGot:  %s\nWant to contain: %s",
			err.Error(),
			"mask logger:",
		)
	}
}

// TestNewMask_TypedNilLogger tests that NewMask returns a wrapped
// ErrNilLogger when a typed nil logger is provided.
func TestNewMask_TypedNilLogger(t *testing.T) {
	t.Parallel()

	var typedNil *mockLogger

	_, err := decorator.NewMask(typedNil)
	if err == nil {
		t.Fatal("expected error when providing a typed nil logger, got nil")
	}

	if !errors.Is(err, decorator.ErrNilLogger) {
		t.Errorf(
			"unexpected error when creating mask logger with typed nil logger:\nGot:  %v\nWant: %v",
			err,
			decorator.ErrNilLogger,
		)
	}

	if !strings.Contains(err.Error(), "mask logger:") {
		t.Errorf(
			"error message does not contain expected text:\nGot:  %s\nWant to contain: %s",
			err.Error(),
			"mask logger:",
		)
	}
}

// TestNewMask_EmptyMaskers tests that NewMask returns a wrapped
// ErrEmptyMaskers when no maskers are provided.
func TestNewMask_EmptyMaskers(t *testing.T) {
	t.Parallel()

	_, err := decorator.NewMask(new(mockLogger))
	if err == nil {
		t.Fatal("expected error when providing no maskers, got nil")
	}

	if !errors.Is(err, decorator.ErrEmptyMaskers) {
		t.Errorf(
			"unexpected error when creating mask logger with no maskers:\nGot:  %v\nWant: %v",
			err,
			decorator.ErrEmptyMaskers,
		)
	}

	if !strings.Contains(err.Error(), "mask logger:") {
		t.Errorf(
			"error message does not contain expected text:\nGot:  %s\nWant to contain: %s",
			err.Error(),
			"mask logger:",
		)
	}
}

// TestNewMask_NilMasker tests that NewMask returns a wrapped ErrNilMasker
// when the first masker is nil.
func TestNewMask_NilMasker(t *testing.T) {
	t.Parallel()

	_, err := decorator.NewMask(new(mockLogger), nil)
	if err == nil {
		t.Fatal("expected error when providing a nil masker, got nil")
	}

	if !errors.Is(err, decorator.ErrNilMasker) {
		t.Errorf(
			"unexpected error when creating mask logger with nil masker:\nGot:  %v\nWant: %v",
			err,
			decorator.ErrNilMasker,
		)
	}

	if !strings.Contains(err.Error(), "mask logger: 1:") {
		t.Errorf(
			"error message does not contain expected text:\nGot:  %s\nWant to contain: %s",
			err.Error(),
			"mask logger: 1:",
		)
	}
}

// TestNewMask_TypedNilMasker tests that NewMask returns a wrapped
// ErrNilMasker when a typed nil masker is provided.
func TestNewMask_TypedNilMasker(t *testing.T) {
	t.Parallel()

	var typedNil *mockMasker

	_, err := decorator.NewMask(new(mockLogger), typedNil)
	if err == nil {
		t.Fatal("expected error when providing a typed nil masker, got nil")
	}

	if !errors.Is(err, decorator.ErrNilMasker) {
		t.Errorf(
			"unexpected error when creating mask logger with typed nil masker:\nGot:  %v\nWant: %v",
			err,
			decorator.ErrNilMasker,
		)
	}

	if !strings.Contains(err.Error(), "mask logger: 1:") {
		t.Errorf(
			"error message does not contain expected text:\nGot:  %s\nWant to contain: %s",
			err.Error(),
			"mask logger: 1:",
		)
	}
}

// TestNewMask_NilMaskerPosition tests that NewMask reports a 1-based
// position for a nil masker that is not the first in the list.
func TestNewMask_NilMaskerPosition(t *testing.T) {
	t.Parallel()

	_, err := decorator.NewMask(new(mockLogger), new(mockMasker), nil, new(mockMasker))
	if err == nil {
		t.Fatal("expected error when providing a nil masker, got nil")
	}

	if !errors.Is(err, decorator.ErrNilMasker) {
		t.Errorf(
			"unexpected error when creating mask logger with nil masker:\nGot:  %v\nWant: %v",
			err,
			decorator.ErrNilMasker,
		)
	}

	if !strings.Contains(err.Error(), "mask logger: 2:") {
		t.Errorf(
			"error message does not contain expected text:\nGot:  %s\nWant to contain: %s",
			err.Error(),
			"mask logger: 2:",
		)
	}
}

// TestNewMask_CopiesMaskers tests that NewMask copies the maskers slice so
// later mutations of the caller's slice do not affect the logger.
func TestNewMask_CopiesMaskers(t *testing.T) {
	t.Parallel()

	mockLog := new(mockLogger)
	mockMasker := new(mockMasker)
	mockMasker.maskMsg = "masked message"

	maskers := []decorator.Masker{mockMasker}

	logger := newMaskLogger(t, mockLog, maskers...)

	maskers[0] = nil

	logger.Info(context.Background(), "message")

	if len(mockLog.infoIn) != 1 {
		t.Fatalf("Expected Info to be called once, got %d calls", len(mockLog.infoIn))
	}

	if mockLog.infoIn[0].msg != "masked message" {
		t.Errorf(
			"Expected Info message to be 'masked message', got '%s'",
			mockLog.infoIn[0].msg,
		)
	}
}

// TestMaskLogger_CopiesArgs tests that maskers receive a copy of the args
// slice, so mutations of the caller's slice or of the slice passed to a
// masker do not affect each other. Nested maps and slices are not copied.
func TestMaskLogger_CopiesArgs(t *testing.T) {
	t.Parallel()

	base := new(mockLogger)
	logger := newMaskLogger(t, base, replaceArgMasker{replaced: "masked"})

	args := []any{"original"}
	logger.Info(context.Background(), "message", args...)

	if args[0] != "original" {
		t.Errorf("expected caller args to stay original, got %v", args[0])
	}

	args[0] = "mutated"

	if len(base.infoIn) != 1 || len(base.infoIn[0].args) != 1 {
		t.Fatalf("expected 1 info call with 1 argument, got %+v", base.infoIn)
	}

	if base.infoIn[0].args[0] != "masked" {
		t.Errorf("expected logger to receive masked args, got %v", base.infoIn[0].args[0])
	}
}

// TestMaskLogger_Debug tests the Debug method of MaskLogger with no arguments.
func TestMaskLogger_Debug(t *testing.T) {
	t.Parallel()

	type ctxKey struct{}

	mockLogger := new(mockLogger)
	expectedCtx := context.WithValue(context.Background(), ctxKey{}, "test debug 1")

	mockMasker := new(mockMasker)
	mockMasker.maskMsg = "masked debug message 1"

	logger, err := decorator.NewMask(mockLogger, mockMasker)
	if err != nil {
		t.Fatalf("failed to create mask logger: %v", err)
	}

	logger.Debug(expectedCtx, "debug message 1")

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

	if mockMasker.msg != "debug message 1" {
		t.Errorf(
			"Expected MaskMessage to be called with 'debug message 1', got '%s'",
			mockMasker.msg,
		)
	}

	if mockLogger.debugIn[0].msg != "masked debug message 1" {
		t.Errorf(
			"Expected Debug message to be 'masked debug message 1', got '%s'",
			mockLogger.debugIn[0].msg,
		)
	}

	if len(mockLogger.debugIn[0].args) != 0 {
		t.Fatalf(
			"Expected Debug to be called with 0 arguments, got %d",
			len(mockLogger.debugIn[0].args),
		)
	}
}

// TestMaskLogger_Debug_WithOneArg tests the Debug method of MaskLogger with one argument.
func TestMaskLogger_Debug_WithOneArg(t *testing.T) {
	t.Parallel()

	type ctxKey struct{}

	mockLogger := new(mockLogger)
	expectedCtx := context.WithValue(context.Background(), ctxKey{}, "test debug 2")

	mockMasker := new(mockMasker)
	mockMasker.maskMsg = "masked debug message 2"
	mockMasker.maskArgs = []any{"masked debug arg 2"}

	logger, err := decorator.NewMask(mockLogger, mockMasker)
	if err != nil {
		t.Fatalf("failed to create mask logger: %v", err)
	}

	logger.Debug(expectedCtx, "debug message 2", "debug 2")

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

	if mockMasker.msg != "debug message 2" {
		t.Errorf(
			"Expected MaskMessage to be called with 'debug message 2', got '%s'",
			mockMasker.msg,
		)
	}

	if mockLogger.debugIn[0].msg != "masked debug message 2" {
		t.Errorf(
			"Expected Debug message to be 'masked debug message 2', got '%s'",
			mockLogger.debugIn[0].msg,
		)
	}

	if len(mockLogger.debugIn[0].args) != 1 {
		t.Fatalf(
			"Expected Debug to be called with 1 argument, got %d",
			len(mockLogger.debugIn[0].args),
		)
	}

	if mockMasker.args[0] != "debug 2" {
		t.Errorf(
			"Expected MaskArguments to be called with 'debug 2', got '%v'",
			mockMasker.args[0],
		)
	}

	if mockLogger.debugIn[0].args[0] != "masked debug arg 2" {
		t.Errorf(
			"Expected first argument to be 'masked debug arg 2', got '%v'",
			mockLogger.debugIn[0].args[0],
		)
	}
}

// TestMaskLogger_Debug_WithMultipleArgs tests the Debug method of MaskLogger with multiple arguments.
func TestMaskLogger_Debug_WithMultipleArgs(t *testing.T) {
	t.Parallel()

	type ctxKey struct{}

	mockLogger := new(mockLogger)
	expectedCtx := context.WithValue(context.Background(), ctxKey{}, "test debug 3")

	mockMasker := new(mockMasker)
	mockMasker.maskMsg = "masked debug message 3"
	mockMasker.maskArgs = []any{"masked debug arg 3-1", 24}

	logger, err := decorator.NewMask(mockLogger, mockMasker)
	if err != nil {
		t.Fatalf("failed to create mask logger: %v", err)
	}

	logger.Debug(expectedCtx, "debug message 3", "debug 3-1", 42)

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

	if mockMasker.msg != "debug message 3" {
		t.Errorf(
			"Expected MaskMessage to be called with 'debug message 3', got '%s'",
			mockMasker.msg,
		)
	}

	if mockLogger.debugIn[0].msg != "masked debug message 3" {
		t.Errorf(
			"Expected Debug message to be 'masked debug message 3', got '%s'",
			mockLogger.debugIn[0].msg,
		)
	}

	if len(mockLogger.debugIn[0].args) != 2 {
		t.Fatalf(
			"Expected Debug to be called with 2 arguments, got %d",
			len(mockLogger.debugIn[0].args),
		)
	}

	if mockMasker.args[0] != "debug 3-1" {
		t.Errorf(
			"Expected first argument to be 'debug 3-1', got '%v'",
			mockMasker.args[0],
		)
	}

	if mockMasker.args[1] != 42 {
		t.Errorf(
			"Expected second argument to be 42, got '%v'",
			mockMasker.args[1],
		)
	}

	if mockLogger.debugIn[0].args[0] != "masked debug arg 3-1" {
		t.Errorf(
			"Expected first argument to be 'masked debug arg 3-1', got '%v'",
			mockLogger.debugIn[0].args[0],
		)
	}

	if mockLogger.debugIn[0].args[1] != 24 {
		t.Errorf(
			"Expected second argument to be 24, got '%v'",
			mockLogger.debugIn[0].args[1],
		)
	}
}

// TestMaskLogger_Debug_WithMultipleMaskers tests the Debug method of MaskLogger with multiple maskers.
func TestMaskLogger_Debug_WithMultipleMaskers(t *testing.T) {
	t.Parallel()

	type ctxKey struct{}

	mockLogger := new(mockLogger)
	expectedCtx := context.WithValue(context.Background(), ctxKey{}, "test debug 4")

	mockMasker1 := new(mockMasker)
	mockMasker1.maskMsg = "masked debug message 4-1"
	mockMasker1.maskArgs = []any{"masked debug arg 4-1"}

	mockMasker2 := new(mockMasker)
	mockMasker2.maskMsg = "masked debug message 4-2"
	mockMasker2.maskArgs = []any{"masked debug arg 4-2"}

	logger, err := decorator.NewMask(mockLogger, mockMasker1, mockMasker2)
	if err != nil {
		t.Fatalf("failed to create mask logger: %v", err)
	}

	logger.Debug(expectedCtx, "debug message 4", "debug 4-1")

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

	if mockMasker1.msg != "debug message 4" {
		t.Errorf(
			"Expected first MaskMessage to be called with 'debug message 4', got '%s'",
			mockMasker1.msg,
		)
	}

	if mockMasker2.msg != "masked debug message 4-1" {
		t.Errorf(
			"Expected second MaskMessage to be called with 'masked debug message 4-1', got '%s'",
			mockMasker2.msg,
		)
	}

	if mockLogger.debugIn[0].msg != "masked debug message 4-2" {
		t.Errorf(
			"Expected Debug message to be 'masked debug message 4-2', got '%s'",
			mockLogger.debugIn[0].msg,
		)
	}

	if mockMasker1.args[0] != "debug 4-1" {
		t.Errorf(
			"Expected first MaskArguments to be called with 'debug 4-1', got '%v'",
			mockMasker1.args[0],
		)
	}

	if mockMasker2.args[0] != "masked debug arg 4-1" {
		t.Errorf(
			"Expected second MaskArguments to be called with 'masked debug arg 4-1', got '%v'",
			mockMasker2.args[0],
		)
	}

	if len(mockLogger.debugIn[0].args) != 1 {
		t.Fatalf(
			"Expected Debug to be called with 1 argument, got %d",
			len(mockLogger.debugIn[0].args),
		)
	}

	if mockLogger.debugIn[0].args[0] != "masked debug arg 4-2" {
		t.Errorf(
			"Expected first argument to be 'masked debug arg 4-2', got '%v'",
			mockLogger.debugIn[0].args[0],
		)
	}
}

// TestMaskLogger_Info tests the Info method of MaskLogger with no arguments.
func TestMaskLogger_Info(t *testing.T) {
	t.Parallel()

	type ctxKey struct{}

	mockLogger := new(mockLogger)
	expectedCtx := context.WithValue(context.Background(), ctxKey{}, "test info 1")

	mockMasker := new(mockMasker)
	mockMasker.maskMsg = "masked info message 1"

	logger, err := decorator.NewMask(mockLogger, mockMasker)
	if err != nil {
		t.Fatalf("failed to create mask logger: %v", err)
	}

	logger.Info(expectedCtx, "info message 1")

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

	if mockMasker.msg != "info message 1" {
		t.Errorf(
			"Expected MaskMessage to be called with 'info message 1', got '%s'",
			mockMasker.msg,
		)
	}

	if mockLogger.infoIn[0].msg != "masked info message 1" {
		t.Errorf(
			"Expected Info message to be 'masked info message 1', got '%s'",
			mockLogger.infoIn[0].msg,
		)
	}

	if len(mockLogger.infoIn[0].args) != 0 {
		t.Fatalf(
			"Expected Info to be called with 0 arguments, got %d",
			len(mockLogger.infoIn[0].args),
		)
	}
}

// TestMaskLogger_Info_WithOneArg tests the Info method of MaskLogger with one argument.
func TestMaskLogger_Info_WithOneArg(t *testing.T) {
	t.Parallel()

	type ctxKey struct{}

	mockLogger := new(mockLogger)
	expectedCtx := context.WithValue(context.Background(), ctxKey{}, "test info 2")

	mockMasker := new(mockMasker)
	mockMasker.maskMsg = "masked info message 2"
	mockMasker.maskArgs = []any{"masked info arg 2"}

	logger, err := decorator.NewMask(mockLogger, mockMasker)
	if err != nil {
		t.Fatalf("failed to create mask logger: %v", err)
	}

	logger.Info(expectedCtx, "info message 2", "info 2")

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

	if mockMasker.msg != "info message 2" {
		t.Errorf(
			"Expected MaskMessage to be called with 'info message 2', got '%s'",
			mockMasker.msg,
		)
	}

	if mockLogger.infoIn[0].msg != "masked info message 2" {
		t.Errorf(
			"Expected Info message to be 'masked info message 2', got '%s'",
			mockLogger.infoIn[0].msg,
		)
	}

	if len(mockLogger.infoIn[0].args) != 1 {
		t.Fatalf(
			"Expected Info to be called with 1 argument, got %d",
			len(mockLogger.infoIn[0].args),
		)
	}

	if mockMasker.args[0] != "info 2" {
		t.Errorf(
			"Expected MaskArguments to be called with 'info 2', got '%v'",
			mockMasker.args[0],
		)
	}

	if mockLogger.infoIn[0].args[0] != "masked info arg 2" {
		t.Errorf(
			"Expected first argument to be 'masked info arg 2', got '%v'",
			mockLogger.infoIn[0].args[0],
		)
	}
}

// TestMaskLogger_Info_WithMultipleArgs tests the Info method of MaskLogger with multiple arguments.
func TestMaskLogger_Info_WithMultipleArgs(t *testing.T) {
	t.Parallel()

	type ctxKey struct{}

	mockLogger := new(mockLogger)
	expectedCtx := context.WithValue(context.Background(), ctxKey{}, "test info 3")

	mockMasker := new(mockMasker)
	mockMasker.maskMsg = "masked info message 3"
	mockMasker.maskArgs = []any{"masked info arg 3-1", 24}

	logger, err := decorator.NewMask(mockLogger, mockMasker)
	if err != nil {
		t.Fatalf("failed to create mask logger: %v", err)
	}

	logger.Info(expectedCtx, "info message 3", "info 3-1", 42)

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

	if mockMasker.msg != "info message 3" {
		t.Errorf(
			"Expected MaskMessage to be called with 'info message 3', got '%s'",
			mockMasker.msg,
		)
	}

	if mockLogger.infoIn[0].msg != "masked info message 3" {
		t.Errorf(
			"Expected Info message to be 'masked info message 3', got '%s'",
			mockLogger.infoIn[0].msg,
		)
	}

	if len(mockLogger.infoIn[0].args) != 2 {
		t.Fatalf(
			"Expected Info to be called with 2 arguments, got %d",
			len(mockLogger.infoIn[0].args),
		)
	}

	if mockMasker.args[0] != "info 3-1" {
		t.Errorf(
			"Expected first argument to be 'info 3-1', got '%v'",
			mockMasker.args[0],
		)
	}

	if mockMasker.args[1] != 42 {
		t.Errorf(
			"Expected second argument to be 42, got '%v'",
			mockMasker.args[1],
		)
	}

	if mockLogger.infoIn[0].args[0] != "masked info arg 3-1" {
		t.Errorf(
			"Expected first argument to be 'masked info arg 3-1', got '%v'",
			mockLogger.infoIn[0].args[0],
		)
	}

	if mockLogger.infoIn[0].args[1] != 24 {
		t.Errorf(
			"Expected second argument to be 24, got '%v'",
			mockLogger.infoIn[0].args[1],
		)
	}
}

// TestMaskLogger_Info_WithMultipleMaskers tests the Info method of MaskLogger with multiple maskers.
func TestMaskLogger_Info_WithMultipleMaskers(t *testing.T) {
	t.Parallel()

	type ctxKey struct{}

	mockLogger := new(mockLogger)
	expectedCtx := context.WithValue(context.Background(), ctxKey{}, "test info 4")

	mockMasker1 := new(mockMasker)
	mockMasker1.maskMsg = "masked info message 4-1"
	mockMasker1.maskArgs = []any{"masked info arg 4-1"}

	mockMasker2 := new(mockMasker)
	mockMasker2.maskMsg = "masked info message 4-2"
	mockMasker2.maskArgs = []any{"masked info arg 4-2"}

	logger, err := decorator.NewMask(mockLogger, mockMasker1, mockMasker2)
	if err != nil {
		t.Fatalf("failed to create mask logger: %v", err)
	}

	logger.Info(expectedCtx, "info message 4", "info 4-1")

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

	if mockMasker1.msg != "info message 4" {
		t.Errorf(
			"Expected first MaskMessage to be called with 'info message 4', got '%s'",
			mockMasker1.msg,
		)
	}

	if mockMasker2.msg != "masked info message 4-1" {
		t.Errorf(
			"Expected second MaskMessage to be called with 'masked info message 4-1', got '%s'",
			mockMasker2.msg,
		)
	}

	if mockLogger.infoIn[0].msg != "masked info message 4-2" {
		t.Errorf(
			"Expected Info message to be 'masked info message 4-2', got '%s'",
			mockLogger.infoIn[0].msg,
		)
	}

	if mockMasker1.args[0] != "info 4-1" {
		t.Errorf(
			"Expected first MaskArguments to be called with 'info 4-1', got '%v'",
			mockMasker1.args[0],
		)
	}

	if mockMasker2.args[0] != "masked info arg 4-1" {
		t.Errorf(
			"Expected second MaskArguments to be called with 'masked info arg 4-1', got '%v'",
			mockMasker2.args[0],
		)
	}

	if len(mockLogger.infoIn[0].args) != 1 {
		t.Fatalf(
			"Expected Info to be called with 1 argument, got %d",
			len(mockLogger.infoIn[0].args),
		)
	}

	if mockLogger.infoIn[0].args[0] != "masked info arg 4-2" {
		t.Errorf(
			"Expected first argument to be 'masked info arg 4-2', got '%v'",
			mockLogger.infoIn[0].args[0],
		)
	}
}

// TestMaskLogger_Warn tests the Warn method of MaskLogger with no arguments.
func TestMaskLogger_Warn(t *testing.T) {
	t.Parallel()

	type ctxKey struct{}

	mockLogger := new(mockLogger)
	expectedCtx := context.WithValue(context.Background(), ctxKey{}, "test warn 1")

	mockMasker := new(mockMasker)
	mockMasker.maskMsg = "masked warn message 1"

	logger, err := decorator.NewMask(mockLogger, mockMasker)
	if err != nil {
		t.Fatalf("failed to create mask logger: %v", err)
	}

	logger.Warn(expectedCtx, "warn message 1")

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

	if mockMasker.msg != "warn message 1" {
		t.Errorf(
			"Expected MaskMessage to be called with 'warn message 1', got '%s'",
			mockMasker.msg,
		)
	}

	if mockLogger.warnIn[0].msg != "masked warn message 1" {
		t.Errorf(
			"Expected Warn message to be 'masked warn message 1', got '%s'",
			mockLogger.warnIn[0].msg,
		)
	}

	if len(mockLogger.warnIn[0].args) != 0 {
		t.Fatalf(
			"Expected Warn to be called with 0 arguments, got %d",
			len(mockLogger.warnIn[0].args),
		)
	}
}

// TestMaskLogger_Warn_WithOneArg tests the Warn method of MaskLogger with one argument.
func TestMaskLogger_Warn_WithOneArg(t *testing.T) {
	t.Parallel()

	type ctxKey struct{}

	mockLogger := new(mockLogger)
	expectedCtx := context.WithValue(context.Background(), ctxKey{}, "test warn 2")

	mockMasker := new(mockMasker)
	mockMasker.maskMsg = "masked warn message 2"
	mockMasker.maskArgs = []any{"masked warn arg 2"}

	logger, err := decorator.NewMask(mockLogger, mockMasker)
	if err != nil {
		t.Fatalf("failed to create mask logger: %v", err)
	}

	logger.Warn(expectedCtx, "warn message 2", "warn 2")

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

	if mockMasker.msg != "warn message 2" {
		t.Errorf(
			"Expected MaskMessage to be called with 'warn message 2', got '%s'",
			mockMasker.msg,
		)
	}

	if mockLogger.warnIn[0].msg != "masked warn message 2" {
		t.Errorf(
			"Expected Warn message to be 'masked warn message 2', got '%s'",
			mockLogger.warnIn[0].msg,
		)
	}

	if len(mockLogger.warnIn[0].args) != 1 {
		t.Fatalf(
			"Expected Warn to be called with 1 argument, got %d",
			len(mockLogger.warnIn[0].args),
		)
	}

	if mockMasker.args[0] != "warn 2" {
		t.Errorf(
			"Expected MaskArguments to be called with 'warn 2', got '%v'",
			mockMasker.args[0],
		)
	}

	if mockLogger.warnIn[0].args[0] != "masked warn arg 2" {
		t.Errorf(
			"Expected first argument to be 'masked warn arg 2', got '%v'",
			mockLogger.warnIn[0].args[0],
		)
	}
}

// TestMaskLogger_Warn_WithMultipleArgs tests the Warn method of MaskLogger with multiple arguments.
func TestMaskLogger_Warn_WithMultipleArgs(t *testing.T) {
	t.Parallel()

	type ctxKey struct{}

	mockLogger := new(mockLogger)
	expectedCtx := context.WithValue(context.Background(), ctxKey{}, "test warn 3")

	mockMasker := new(mockMasker)
	mockMasker.maskMsg = "masked warn message 3"
	mockMasker.maskArgs = []any{"masked warn arg 3-1", 24}

	logger, err := decorator.NewMask(mockLogger, mockMasker)
	if err != nil {
		t.Fatalf("failed to create mask logger: %v", err)
	}

	logger.Warn(expectedCtx, "warn message 3", "warn 3-1", 42)

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

	if mockMasker.msg != "warn message 3" {
		t.Errorf(
			"Expected MaskMessage to be called with 'warn message 3', got '%s'",
			mockMasker.msg,
		)
	}

	if mockLogger.warnIn[0].msg != "masked warn message 3" {
		t.Errorf(
			"Expected Warn message to be 'masked warn message 3', got '%s'",
			mockLogger.warnIn[0].msg,
		)
	}

	if len(mockLogger.warnIn[0].args) != 2 {
		t.Fatalf(
			"Expected Warn to be called with 2 arguments, got %d",
			len(mockLogger.warnIn[0].args),
		)
	}

	if mockMasker.args[0] != "warn 3-1" {
		t.Errorf(
			"Expected first argument to be 'warn 3-1', got '%v'",
			mockMasker.args[0],
		)
	}

	if mockMasker.args[1] != 42 {
		t.Errorf(
			"Expected second argument to be 42, got '%v'",
			mockMasker.args[1],
		)
	}

	if mockLogger.warnIn[0].args[0] != "masked warn arg 3-1" {
		t.Errorf(
			"Expected first argument to be 'masked warn arg 3-1', got '%v'",
			mockLogger.warnIn[0].args[0],
		)
	}

	if mockLogger.warnIn[0].args[1] != 24 {
		t.Errorf(
			"Expected second argument to be 24, got '%v'",
			mockLogger.warnIn[0].args[1],
		)
	}
}

// TestMaskLogger_Warn_WithMultipleMaskers tests the Warn method of MaskLogger with multiple maskers.
func TestMaskLogger_Warn_WithMultipleMaskers(t *testing.T) {
	t.Parallel()

	type ctxKey struct{}

	mockLogger := new(mockLogger)
	expectedCtx := context.WithValue(context.Background(), ctxKey{}, "test warn 4")

	mockMasker1 := new(mockMasker)
	mockMasker1.maskMsg = "masked warn message 4-1"
	mockMasker1.maskArgs = []any{"masked warn arg 4-1"}

	mockMasker2 := new(mockMasker)
	mockMasker2.maskMsg = "masked warn message 4-2"
	mockMasker2.maskArgs = []any{"masked warn arg 4-2"}

	logger, err := decorator.NewMask(mockLogger, mockMasker1, mockMasker2)
	if err != nil {
		t.Fatalf("failed to create mask logger: %v", err)
	}

	logger.Warn(expectedCtx, "warn message 4", "warn 4-1")

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

	if mockMasker1.msg != "warn message 4" {
		t.Errorf(
			"Expected first MaskMessage to be called with 'warn message 4', got '%s'",
			mockMasker1.msg,
		)
	}

	if mockMasker2.msg != "masked warn message 4-1" {
		t.Errorf(
			"Expected second MaskMessage to be called with 'masked warn message 4-1', got '%s'",
			mockMasker2.msg,
		)
	}

	if mockLogger.warnIn[0].msg != "masked warn message 4-2" {
		t.Errorf(
			"Expected Warn message to be 'masked warn message 4-2', got '%s'",
			mockLogger.warnIn[0].msg,
		)
	}

	if mockMasker1.args[0] != "warn 4-1" {
		t.Errorf(
			"Expected first MaskArguments to be called with 'warn 4-1', got '%v'",
			mockMasker1.args[0],
		)
	}

	if mockMasker2.args[0] != "masked warn arg 4-1" {
		t.Errorf(
			"Expected second MaskArguments to be called with 'masked warn arg 4-1', got '%v'",
			mockMasker2.args[0],
		)
	}

	if len(mockLogger.warnIn[0].args) != 1 {
		t.Fatalf(
			"Expected Warn to be called with 1 argument, got %d",
			len(mockLogger.warnIn[0].args),
		)
	}

	if mockLogger.warnIn[0].args[0] != "masked warn arg 4-2" {
		t.Errorf(
			"Expected first argument to be 'masked warn arg 4-2', got '%v'",
			mockLogger.warnIn[0].args[0],
		)
	}
}

// TestMaskLogger_Error tests the Error method of MaskLogger with no arguments.
func TestMaskLogger_Error(t *testing.T) {
	t.Parallel()

	type ctxKey struct{}

	mockLogger := new(mockLogger)
	expectedCtx := context.WithValue(context.Background(), ctxKey{}, "test error 1")
	expectedErr := errors.New("error message 1")

	mockMasker := new(mockMasker)
	mockMasker.maskErr = errors.New("masked error message 1")

	logger, err := decorator.NewMask(mockLogger, mockMasker)
	if err != nil {
		t.Fatalf("failed to create mask logger: %v", err)
	}

	logger.Error(expectedCtx, expectedErr)

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

	if !errors.Is(mockMasker.err, expectedErr) {
		t.Errorf(
			"Expected MaskError to be called with %v, got %v",
			expectedErr,
			mockMasker.err,
		)
	}

	if !errors.Is(mockLogger.errorIn[0].err, mockMasker.maskErr) {
		t.Errorf(
			"Expected Error to be called with %v, got %v",
			mockMasker.maskErr,
			mockLogger.errorIn[0].err,
		)
	}

	if len(mockLogger.errorIn[0].args) != 0 {
		t.Fatalf(
			"Expected Error to be called with 0 arguments, got %d",
			len(mockLogger.errorIn[0].args),
		)
	}
}

// TestMaskLogger_Error_WithOneArg tests the Error method of MaskLogger with one argument.
func TestMaskLogger_Error_WithOneArg(t *testing.T) {
	t.Parallel()

	type ctxKey struct{}

	mockLogger := new(mockLogger)
	expectedCtx := context.WithValue(context.Background(), ctxKey{}, "test error 2")
	expectedErr := errors.New("error message 2")

	mockMasker := new(mockMasker)
	mockMasker.maskErr = errors.New("masked error message 2")
	mockMasker.maskArgs = []any{"masked error arg 2"}

	logger, err := decorator.NewMask(mockLogger, mockMasker)
	if err != nil {
		t.Fatalf("failed to create mask logger: %v", err)
	}

	logger.Error(expectedCtx, expectedErr, "error 2")

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

	if !errors.Is(mockMasker.err, expectedErr) {
		t.Errorf(
			"Expected MaskError to be called with %v, got %v",
			expectedErr,
			mockMasker.err,
		)
	}

	if !errors.Is(mockLogger.errorIn[0].err, mockMasker.maskErr) {
		t.Errorf(
			"Expected Error to be called with %v, got %v",
			mockMasker.maskErr,
			mockLogger.errorIn[0].err,
		)
	}

	if len(mockLogger.errorIn[0].args) != 1 {
		t.Fatalf(
			"Expected Error to be called with 1 argument, got %d",
			len(mockLogger.errorIn[0].args),
		)
	}

	if mockMasker.args[0] != "error 2" {
		t.Errorf(
			"Expected MaskArguments to be called with 'error 2', got '%v'",
			mockMasker.args[0],
		)
	}

	if mockLogger.errorIn[0].args[0] != "masked error arg 2" {
		t.Errorf(
			"Expected first argument to be 'masked error arg 2', got '%v'",
			mockLogger.errorIn[0].args[0],
		)
	}
}

// TestMaskLogger_Error_WithMultipleArgs tests the Error method of MaskLogger with multiple arguments.
func TestMaskLogger_Error_WithMultipleArgs(t *testing.T) {
	t.Parallel()

	type ctxKey struct{}

	mockLogger := new(mockLogger)
	expectedCtx := context.WithValue(context.Background(), ctxKey{}, "test error 3")
	expectedErr := errors.New("error message 3")

	mockMasker := new(mockMasker)
	mockMasker.maskErr = errors.New("masked error message 3")
	mockMasker.maskArgs = []any{"masked error arg 3-1", 24}

	logger, err := decorator.NewMask(mockLogger, mockMasker)
	if err != nil {
		t.Fatalf("failed to create mask logger: %v", err)
	}

	logger.Error(expectedCtx, expectedErr, "error 3-1", 42)

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

	if !errors.Is(mockMasker.err, expectedErr) {
		t.Errorf(
			"Expected MaskError to be called with %v, got %v",
			expectedErr,
			mockMasker.err,
		)
	}

	if !errors.Is(mockLogger.errorIn[0].err, mockMasker.maskErr) {
		t.Errorf(
			"Expected Error to be called with %v, got %v",
			mockMasker.maskErr,
			mockLogger.errorIn[0].err,
		)
	}

	if len(mockLogger.errorIn[0].args) != 2 {
		t.Fatalf(
			"Expected Error to be called with 2 arguments, got %d",
			len(mockLogger.errorIn[0].args),
		)
	}

	if mockMasker.args[0] != "error 3-1" {
		t.Errorf(
			"Expected first argument to be 'error 3-1', got '%v'",
			mockMasker.args[0],
		)
	}

	if mockMasker.args[1] != 42 {
		t.Errorf(
			"Expected second argument to be 42, got '%v'",
			mockMasker.args[1],
		)
	}

	if mockLogger.errorIn[0].args[0] != "masked error arg 3-1" {
		t.Errorf(
			"Expected first argument to be 'masked error arg 3-1', got '%v'",
			mockLogger.errorIn[0].args[0],
		)
	}

	if mockLogger.errorIn[0].args[1] != 24 {
		t.Errorf(
			"Expected second argument to be 24, got '%v'",
			mockLogger.errorIn[0].args[1],
		)
	}
}

// TestMaskLogger_Error_WithMultipleMaskers tests the Error method of MaskLogger with multiple maskers.
func TestMaskLogger_Error_WithMultipleMaskers(t *testing.T) {
	t.Parallel()

	type ctxKey struct{}

	mockLogger := new(mockLogger)
	expectedCtx := context.WithValue(context.Background(), ctxKey{}, "test error 4")
	expectedErr := errors.New("error message 4")

	mockMasker1 := new(mockMasker)
	mockMasker1.maskErr = errors.New("masked error message 4-1")
	mockMasker1.maskArgs = []any{"masked error arg 4-1"}

	mockMasker2 := new(mockMasker)
	mockMasker2.maskErr = errors.New("masked error message 4-2")
	mockMasker2.maskArgs = []any{"masked error arg 4-2"}

	logger, err := decorator.NewMask(mockLogger, mockMasker1, mockMasker2)
	if err != nil {
		t.Fatalf("failed to create mask logger: %v", err)
	}

	logger.Error(expectedCtx, expectedErr, "error 4-1")

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

	if !errors.Is(mockMasker1.err, expectedErr) {
		t.Errorf(
			"Expected first MaskError to be called with %v, got %v",
			expectedErr,
			mockMasker1.err,
		)
	}

	if !errors.Is(mockMasker2.err, mockMasker1.maskErr) {
		t.Errorf(
			"Expected second MaskError to be called with %v, got %v",
			mockMasker1.maskErr,
			mockMasker2.err,
		)
	}

	if !errors.Is(mockLogger.errorIn[0].err, mockMasker2.maskErr) {
		t.Errorf(
			"Expected Error to be called with %v, got %v",
			mockMasker2.maskErr,
			mockLogger.errorIn[0].err,
		)
	}

	if mockMasker1.args[0] != "error 4-1" {
		t.Errorf(
			"Expected first MaskArguments to be called with 'error 4-1', got '%v'",
			mockMasker1.args[0],
		)
	}

	if mockMasker2.args[0] != "masked error arg 4-1" {
		t.Errorf(
			"Expected second MaskArguments to be called with 'masked error arg 4-1', got '%v'",
			mockMasker2.args[0],
		)
	}

	if len(mockLogger.errorIn[0].args) != 1 {
		t.Fatalf(
			"Expected Error to be called with 1 argument, got %d",
			len(mockLogger.errorIn[0].args),
		)
	}

	if mockLogger.errorIn[0].args[0] != "masked error arg 4-2" {
		t.Errorf(
			"Expected first argument to be 'masked error arg 4-2', got '%v'",
			mockLogger.errorIn[0].args[0],
		)
	}
}

// TestMaskLogger_Error_NilError tests that Error forwards a nil error through MaskError.
func TestMaskLogger_Error_NilError(t *testing.T) {
	t.Parallel()

	type ctxKey struct{}

	mockLogger := new(mockLogger)
	expectedCtx := context.WithValue(context.Background(), ctxKey{}, "test error nil")

	mockMasker := new(mockMasker)

	logger, err := decorator.NewMask(mockLogger, mockMasker)
	if err != nil {
		t.Fatalf("failed to create mask logger: %v", err)
	}

	logger.Error(expectedCtx, nil)

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

	if mockMasker.err != nil {
		t.Errorf("Expected MaskError to be called with nil, got %v", mockMasker.err)
	}

	if mockLogger.errorIn[0].err != nil {
		t.Errorf("Expected Error to be called with nil, got %v", mockLogger.errorIn[0].err)
	}
}

// TestMaskLogger_Error_MaskErrorReturnsNil tests that a nil result from MaskError
// is forwarded to the underlying logger.
func TestMaskLogger_Error_MaskErrorReturnsNil(t *testing.T) {
	t.Parallel()

	type ctxKey struct{}

	mockLogger := new(mockLogger)
	expectedCtx := context.WithValue(context.Background(), ctxKey{}, "test error returns nil")
	expectedErr := errors.New("error message")

	mockMasker := new(mockMasker)

	logger, err := decorator.NewMask(mockLogger, mockMasker)
	if err != nil {
		t.Fatalf("failed to create mask logger: %v", err)
	}

	logger.Error(expectedCtx, expectedErr)

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

	if !errors.Is(mockMasker.err, expectedErr) {
		t.Errorf(
			"Expected MaskError to be called with %v, got %v",
			expectedErr,
			mockMasker.err,
		)
	}

	if mockLogger.errorIn[0].err != nil {
		t.Errorf("Expected Error to be called with nil, got %v", mockLogger.errorIn[0].err)
	}
}

// TestMaskLogger_Info_NilArgsPreserved tests that a nil call-site args slice
// is passed to the first masker as nil, not as an empty non-nil slice.
func TestMaskLogger_Info_NilArgsPreserved(t *testing.T) {
	t.Parallel()

	mockLogger := new(mockLogger)
	mockMasker := new(mockMasker)
	mockMasker.maskMsg = "masked message"

	logger, err := decorator.NewMask(mockLogger, mockMasker)
	if err != nil {
		t.Fatalf("failed to create mask logger: %v", err)
	}

	logger.Info(context.Background(), "message")

	if mockMasker.args != nil {
		t.Errorf("expected first MaskArguments to receive nil args, got %v", mockMasker.args)
	}

	if len(mockLogger.infoIn) != 1 {
		t.Fatalf("Expected Info to be called once, got %d calls", len(mockLogger.infoIn))
	}

	if mockLogger.infoIn[0].args != nil {
		t.Errorf("expected logger to receive nil args, got %v", mockLogger.infoIn[0].args)
	}
}

// TestMaskLogger_Info_MaskArgumentsReturnsNil tests that a nil result from
// MaskArguments is forwarded as zero arguments to the next masker and logger.
func TestMaskLogger_Info_MaskArgumentsReturnsNil(t *testing.T) {
	t.Parallel()

	type ctxKey struct{}

	mockLogger := new(mockLogger)
	expectedCtx := context.WithValue(context.Background(), ctxKey{}, "test args nil")

	mockMasker1 := new(mockMasker)
	mockMasker1.returnNilArgs = true

	mockMasker2 := new(mockMasker)

	logger, err := decorator.NewMask(mockLogger, mockMasker1, mockMasker2)
	if err != nil {
		t.Fatalf("failed to create mask logger: %v", err)
	}

	logger.Info(expectedCtx, "message", "arg1", "arg2")

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

	if len(mockMasker1.args) != 2 {
		t.Fatalf(
			"Expected first MaskArguments to be called with 2 arguments, got %d",
			len(mockMasker1.args),
		)
	}

	if mockMasker1.args[0] != "arg1" {
		t.Errorf(
			"Expected first MaskArguments to be called with 'arg1', got '%v'",
			mockMasker1.args[0],
		)
	}

	if mockMasker1.args[1] != "arg2" {
		t.Errorf(
			"Expected first MaskArguments to be called with 'arg2', got '%v'",
			mockMasker1.args[1],
		)
	}

	if len(mockMasker2.args) != 0 {
		t.Fatalf(
			"Expected second MaskArguments to be called with 0 arguments, got %d",
			len(mockMasker2.args),
		)
	}

	if len(mockLogger.infoIn[0].args) != 0 {
		t.Fatalf(
			"Expected Info to be called with 0 arguments, got %d",
			len(mockLogger.infoIn[0].args),
		)
	}
}

// TestMaskLogger_Info_MaskArgumentsChangesLength tests that MaskArguments may
// return a slice of a different length than it received.
func TestMaskLogger_Info_MaskArgumentsChangesLength(t *testing.T) {
	t.Parallel()

	type ctxKey struct{}

	mockLogger := new(mockLogger)
	expectedCtx := context.WithValue(context.Background(), ctxKey{}, "test args length")

	mockMasker := new(mockMasker)
	mockMasker.maskMsg = "masked message"
	mockMasker.maskArgs = []any{"only"}

	logger, err := decorator.NewMask(mockLogger, mockMasker)
	if err != nil {
		t.Fatalf("failed to create mask logger: %v", err)
	}

	logger.Info(expectedCtx, "message", "arg1", "arg2")

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

	if mockLogger.infoIn[0].msg != "masked message" {
		t.Errorf(
			"Expected Info message to be 'masked message', got '%s'",
			mockLogger.infoIn[0].msg,
		)
	}

	if len(mockMasker.args) != 2 {
		t.Fatalf(
			"Expected MaskArguments to be called with 2 arguments, got %d",
			len(mockMasker.args),
		)
	}

	if len(mockLogger.infoIn[0].args) != 1 {
		t.Fatalf(
			"Expected Info to be called with 1 argument, got %d",
			len(mockLogger.infoIn[0].args),
		)
	}

	if mockLogger.infoIn[0].args[0] != "only" {
		t.Errorf(
			"Expected first argument to be 'only', got '%v'",
			mockLogger.infoIn[0].args[0],
		)
	}
}
