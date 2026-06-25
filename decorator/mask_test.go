package decorator_test

import (
	"context"
	"errors"
	"testing"

	"github.com/yuppyweb/cakelog/decorator"
)

// mockMasker is a mock implementation of the Masker interface used for testing.
type mockMasker struct {
	msg      string
	maskMsg  string
	err      error
	maskErr  error
	args     []any
	maskArgs []any
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

// MaskArgument masks a single argument and returns the masked version.
func (m *mockMasker) MaskArgument(arg any) any {
	idx := len(m.args)
	args := make([]any, idx+1)

	copy(args, m.args)
	args[idx] = arg

	m.args = args

	if idx < len(m.maskArgs) {
		return m.maskArgs[idx]
	}

	return nil
}

var _ decorator.Masker = (*mockMasker)(nil)

// TestNewMaskLogger_NilLogger tests that NewMaskLogger returns an error when given a nil logger.
func TestNewMaskLogger_NilLogger(t *testing.T) {
	t.Parallel()

	_, err := decorator.NewMaskLogger(nil)
	if err == nil {
		t.Fatal("expected error when creating MaskLogger with nil logger, got nil")
	}

	if !errors.Is(err, decorator.ErrNilLogger) {
		t.Errorf("expected error to be ErrNilLogger, got %v", err)
	}
}

// TestNewMaskLogger_EmptyMaskers tests that NewMaskLogger returns an error when given no maskers.
func TestNewMaskLogger_EmptyMaskers(t *testing.T) {
	t.Parallel()

	_, err := decorator.NewMaskLogger(new(mockLogger))
	if err == nil {
		t.Fatal("expected error when creating MaskLogger with empty maskers, got nil")
	}

	if !errors.Is(err, decorator.ErrEmptyMaskers) {
		t.Errorf("expected error to be ErrEmptyMaskers, got %v", err)
	}
}

// TestNewMaskLogger_NilMasker tests that NewMaskLogger returns an error when given a nil masker.
func TestNewMaskLogger_NilMasker(t *testing.T) {
	t.Parallel()

	_, err := decorator.NewMaskLogger(new(mockLogger), nil)
	if err == nil {
		t.Fatal("expected error when creating MaskLogger with nil masker, got nil")
	}

	if !errors.Is(err, decorator.ErrNilMasker) {
		t.Errorf("expected error to be ErrNilMasker, got %v", err)
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

	logger, err := decorator.NewMaskLogger(mockLogger, mockMasker)
	if err != nil {
		t.Fatalf("failed to create MaskLogger: %v", err)
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

	logger, err := decorator.NewMaskLogger(mockLogger, mockMasker)
	if err != nil {
		t.Fatalf("failed to create MaskLogger: %v", err)
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
			"Expected MaskArgument to be called with 'debug 2', got '%v'",
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

	logger, err := decorator.NewMaskLogger(mockLogger, mockMasker)
	if err != nil {
		t.Fatalf("failed to create MaskLogger: %v", err)
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

	logger, err := decorator.NewMaskLogger(mockLogger, mockMasker1, mockMasker2)
	if err != nil {
		t.Fatalf("failed to create MaskLogger: %v", err)
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
			"Expected first MaskArgument to be called with 'debug 4-1', got '%v'",
			mockMasker1.args[0],
		)
	}

	if mockMasker2.args[0] != "masked debug arg 4-1" {
		t.Errorf(
			"Expected second MaskArgument to be called with 'masked debug arg 4-1', got '%v'",
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

	logger, err := decorator.NewMaskLogger(mockLogger, mockMasker)
	if err != nil {
		t.Fatalf("failed to create MaskLogger: %v", err)
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

	logger, err := decorator.NewMaskLogger(mockLogger, mockMasker)
	if err != nil {
		t.Fatalf("failed to create MaskLogger: %v", err)
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
			"Expected MaskArgument to be called with 'info 2', got '%v'",
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

	logger, err := decorator.NewMaskLogger(mockLogger, mockMasker)
	if err != nil {
		t.Fatalf("failed to create MaskLogger: %v", err)
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

	logger, err := decorator.NewMaskLogger(mockLogger, mockMasker1, mockMasker2)
	if err != nil {
		t.Fatalf("failed to create MaskLogger: %v", err)
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
			"Expected first MaskArgument to be called with 'info 4-1', got '%v'",
			mockMasker1.args[0],
		)
	}

	if mockMasker2.args[0] != "masked info arg 4-1" {
		t.Errorf(
			"Expected second MaskArgument to be called with 'masked info arg 4-1', got '%v'",
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

	logger, err := decorator.NewMaskLogger(mockLogger, mockMasker)
	if err != nil {
		t.Fatalf("failed to create MaskLogger: %v", err)
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

	logger, err := decorator.NewMaskLogger(mockLogger, mockMasker)
	if err != nil {
		t.Fatalf("failed to create MaskLogger: %v", err)
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
			"Expected MaskArgument to be called with 'warn 2', got '%v'",
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

	logger, err := decorator.NewMaskLogger(mockLogger, mockMasker)
	if err != nil {
		t.Fatalf("failed to create MaskLogger: %v", err)
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

	logger, err := decorator.NewMaskLogger(mockLogger, mockMasker1, mockMasker2)
	if err != nil {
		t.Fatalf("failed to create MaskLogger: %v", err)
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
			"Expected first MaskArgument to be called with 'warn 4-1', got '%v'",
			mockMasker1.args[0],
		)
	}

	if mockMasker2.args[0] != "masked warn arg 4-1" {
		t.Errorf(
			"Expected second MaskArgument to be called with 'masked warn arg 4-1', got '%v'",
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

	logger, err := decorator.NewMaskLogger(mockLogger, mockMasker)
	if err != nil {
		t.Fatalf("failed to create MaskLogger: %v", err)
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

	logger, err := decorator.NewMaskLogger(mockLogger, mockMasker)
	if err != nil {
		t.Fatalf("failed to create MaskLogger: %v", err)
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
			"Expected MaskArgument to be called with 'error 2', got '%v'",
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

	logger, err := decorator.NewMaskLogger(mockLogger, mockMasker)
	if err != nil {
		t.Fatalf("failed to create MaskLogger: %v", err)
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

	logger, err := decorator.NewMaskLogger(mockLogger, mockMasker1, mockMasker2)
	if err != nil {
		t.Fatalf("failed to create MaskLogger: %v", err)
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
			"Expected first MaskArgument to be called with 'error 4-1', got '%v'",
			mockMasker1.args[0],
		)
	}

	if mockMasker2.args[0] != "masked error arg 4-1" {
		t.Errorf(
			"Expected second MaskArgument to be called with 'masked error arg 4-1', got '%v'",
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
