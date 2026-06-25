package decorator_test

import (
	"context"
	"errors"
	"testing"

	"github.com/yuppyweb/cakelog/decorator"
)

// TestWithContextValue tests that WithContextValue correctly stores and retrieves values
// in context chains, ensuring that each context level maintains its own values and
// can access values from parent contexts.
func TestWithContextValue(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	ctx1 := decorator.WithContextValue(ctx, "key1", "value1")
	ctx2 := decorator.WithContextValue(ctx1, "key2", "value2")

	if val := decorator.ContextValue(ctx, "key1"); val != nil {
		t.Errorf("Expected base context to not contain key 'key1'")
	}

	if val := decorator.ContextValue(ctx, "key2"); val != nil {
		t.Errorf("Expected base context to not contain key 'key2'")
	}

	if val := decorator.ContextValue(ctx1, "key1"); val != "value1" {
		t.Errorf("Expected context value for 'key1' to be 'value1', got '%v'", val)
	}

	if val := decorator.ContextValue(ctx1, "key2"); val != nil {
		t.Errorf("Expected context to not contain key 'key2'")
	}

	if val := decorator.ContextValue(ctx2, "key1"); val != "value1" {
		t.Errorf("Expected context value for 'key1' to be 'value1', got '%v'", val)
	}

	if val := decorator.ContextValue(ctx2, "key2"); val != "value2" {
		t.Errorf("Expected context value for 'key2' to be 'value2', got '%v'", val)
	}
}

// TestContextValue_EmptyContextKey tests that ContextValue returns nil when querying
// with an empty string key.
func TestContextValue_EmptyContextKey(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	if val := decorator.ContextValue(ctx, ""); val != nil {
		t.Errorf("Expected context to not contain empty string key")
	}
}

// TestContextValue_NonExistentKey tests that ContextValue correctly retrieves existing keys
// and returns nil for keys that do not exist in the context.
func TestContextValue_NonExistentKey(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctx = decorator.WithContextValue(ctx, "existingKey", "value")

	if val := decorator.ContextValue(ctx, "existingKey"); val != "value" {
		t.Errorf("Expected context value for 'existingKey' to be 'value', got '%v'", val)
	}

	if val := decorator.ContextValue(ctx, "nonExistentKey"); val != nil {
		t.Errorf("Expected context to not contain key 'nonExistentKey'")
	}
}

// TestNewContextLogger_NilLogger tests that NewContextLogger returns an error of type
// ErrNilLogger when provided with a nil logger.
func TestNewContextLogger_NilLogger(t *testing.T) {
	t.Parallel()

	_, err := decorator.NewContextLogger(nil)
	if err == nil {
		t.Fatal("expected error when providing nil logger, got nil")
	}

	if !errors.Is(err, decorator.ErrNilLogger) {
		t.Errorf(
			"unexpected error when creating ContextLogger with nil logger:\nGot:  %v\nWant: %v",
			err,
			decorator.ErrNilLogger,
		)
	}
}

// TestContextLogger_Debug tests that the ContextLogger's Debug method correctly passes
// the debug message and arguments to the underlying logger, and includes context values
// in the arguments.
func TestContextLogger_Debug(t *testing.T) {
	t.Parallel()

	mockLogger := new(mockLogger)

	logger, err := decorator.NewContextLogger(mockLogger)
	if err != nil {
		t.Fatalf("failed to create ContextLogger: %v", err)
	}

	ctx := decorator.WithContextValue(context.Background(), "userID", 123)

	logger.Debug(ctx, "debug message", "debug", 42)

	if len(mockLogger.debugIn) != 1 {
		t.Fatalf("Expected Debug to be called once, got %d calls", len(mockLogger.debugIn))
	}

	if mockLogger.debugIn[0].msg != "debug message" {
		t.Errorf(
			"Expected Debug message to be 'debug message', got '%s'",
			mockLogger.debugIn[0].msg,
		)
	}

	if len(mockLogger.debugIn[0].args) != 3 {
		t.Fatalf(
			"Expected Debug to be called with 3 arguments, got %d",
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

	ctxValue, ok := mockLogger.debugIn[0].args[2].(map[string]any)
	if !ok {
		t.Fatalf(
			"Expected third argument to be a map[string]any, got '%T'",
			mockLogger.debugIn[0].args[2],
		)
	}

	userID, ok := ctxValue["userID"]
	if !ok {
		t.Fatalf("Expected context to contain key 'userID'")
	}

	if userID != 123 {
		t.Errorf(
			"Expected context value for 'userID' to be 123, got '%v'",
			userID,
		)
	}
}

// TestContextLogger_Info tests that the ContextLogger's Info method correctly passes
// the info message and arguments to the underlying logger, and includes context values
// in the arguments.
func TestContextLogger_Info(t *testing.T) {
	t.Parallel()

	mockLogger := new(mockLogger)

	logger, err := decorator.NewContextLogger(mockLogger)
	if err != nil {
		t.Fatalf("failed to create ContextLogger: %v", err)
	}

	ctx := decorator.WithContextValue(context.Background(), "requestID", "abc-123")

	logger.Info(ctx, "info message", "info", 75)

	if len(mockLogger.infoIn) != 1 {
		t.Fatalf("Expected Info to be called once, got %d calls", len(mockLogger.infoIn))
	}

	if mockLogger.infoIn[0].msg != "info message" {
		t.Errorf(
			"Expected Info message to be 'info message', got '%s'",
			mockLogger.infoIn[0].msg,
		)
	}

	if len(mockLogger.infoIn[0].args) != 3 {
		t.Fatalf(
			"Expected Info to be called with 3 arguments, got %d",
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

	ctxValue, ok := mockLogger.infoIn[0].args[2].(map[string]any)
	if !ok {
		t.Fatalf(
			"Expected third argument to be a map[string]any, got '%T'",
			mockLogger.infoIn[0].args[2],
		)
	}

	requestID, ok := ctxValue["requestID"]
	if !ok {
		t.Fatalf("Expected context to contain key 'requestID'")
	}

	if requestID != "abc-123" {
		t.Errorf(
			"Expected context value for 'requestID' to be 'abc-123', got '%v'",
			requestID,
		)
	}
}

// TestContextLogger_Warn tests that the ContextLogger's Warn method correctly passes
// the warn message and arguments to the underlying logger, and includes context values
// in the arguments.
func TestContextLogger_Warn(t *testing.T) {
	t.Parallel()

	mockLogger := new(mockLogger)

	logger, err := decorator.NewContextLogger(mockLogger)
	if err != nil {
		t.Fatalf("failed to create ContextLogger: %v", err)
	}

	ctx := decorator.WithContextValue(context.Background(), "sessionID", "xyz-789")

	logger.Warn(ctx, "warn message", "warn", 88)

	if len(mockLogger.warnIn) != 1 {
		t.Fatalf("Expected Warn to be called once, got %d calls", len(mockLogger.warnIn))
	}

	if mockLogger.warnIn[0].msg != "warn message" {
		t.Errorf(
			"Expected Warn message to be 'warn message', got '%s'",
			mockLogger.warnIn[0].msg,
		)
	}

	if len(mockLogger.warnIn[0].args) != 3 {
		t.Fatalf(
			"Expected Warn to be called with 3 arguments, got %d",
			len(mockLogger.warnIn[0].args),
		)
	}

	if mockLogger.warnIn[0].args[0] != "warn" {
		t.Errorf(
			"Expected first argument to be 'warn', got '%v'",
			mockLogger.warnIn[0].args[0],
		)
	}

	if mockLogger.warnIn[0].args[1] != 88 {
		t.Errorf(
			"Expected second argument to be 88, got '%v'",
			mockLogger.warnIn[0].args[1],
		)
	}

	ctxValue, ok := mockLogger.warnIn[0].args[2].(map[string]any)
	if !ok {
		t.Fatalf(
			"Expected third argument to be a map[string]any, got '%T'",
			mockLogger.warnIn[0].args[2],
		)
	}

	sessionID, ok := ctxValue["sessionID"]
	if !ok {
		t.Fatalf("Expected context to contain key 'sessionID'")
	}

	if sessionID != "xyz-789" {
		t.Errorf(
			"Expected context value for 'sessionID' to be 'xyz-789', got '%v'",
			sessionID,
		)
	}
}

// TestContextLogger_Error tests that the ContextLogger's Error method correctly passes
// the error and arguments to the underlying logger, and includes context values
// in the arguments.
func TestContextLogger_Error(t *testing.T) {
	t.Parallel()

	mockLogger := new(mockLogger)

	logger, err := decorator.NewContextLogger(mockLogger)
	if err != nil {
		t.Fatalf("failed to create ContextLogger: %v", err)
	}

	ctx := decorator.WithContextValue(context.Background(), "transactionID", "txn-456")

	expectedErr := errors.New("error message")
	logger.Error(ctx, expectedErr, "error", 90)

	if len(mockLogger.errorIn) != 1 {
		t.Fatalf("Expected Error to be called once, got %d calls", len(mockLogger.errorIn))
	}

	if !errors.Is(mockLogger.errorIn[0].err, expectedErr) {
		t.Errorf(
			"Expected Error to be called with error '%v', got '%v'",
			expectedErr,
			mockLogger.errorIn[0].err,
		)
	}

	if len(mockLogger.errorIn[0].args) != 3 {
		t.Fatalf(
			"Expected Error to be called with 3 arguments, got %d",
			len(mockLogger.errorIn[0].args),
		)
	}

	if mockLogger.errorIn[0].args[0] != "error" {
		t.Errorf(
			"Expected first argument to be 'error', got '%v'",
			mockLogger.errorIn[0].args[0],
		)
	}

	if mockLogger.errorIn[0].args[1] != 90 {
		t.Errorf(
			"Expected second argument to be 90, got '%v'",
			mockLogger.errorIn[0].args[1],
		)
	}

	ctxValue, ok := mockLogger.errorIn[0].args[2].(map[string]any)
	if !ok {
		t.Fatalf(
			"Expected third argument to be a map[string]any, got '%T'",
			mockLogger.errorIn[0].args[2],
		)
	}

	transactionID, ok := ctxValue["transactionID"]
	if !ok {
		t.Fatalf("Expected context to contain key 'transactionID'")
	}

	if transactionID != "txn-456" {
		t.Errorf(
			"Expected context value for 'transactionID' to be 'txn-456', got '%v'",
			transactionID,
		)
	}
}

// TestContextLogger_WithEmptyContext tests that all ContextLogger methods work correctly
// when called with an empty context, passing arguments to the underlying logger without
// any context values.
func TestContextLogger_WithEmptyContext(t *testing.T) {
	t.Parallel()

	mockLogger := new(mockLogger)
	expectedErr := errors.New("error message")

	logger, err := decorator.NewContextLogger(mockLogger)
	if err != nil {
		t.Fatalf("failed to create ContextLogger: %v", err)
	}

	logger.Debug(context.Background(), "debug message")
	logger.Info(context.Background(), "info message")
	logger.Warn(context.Background(), "warn message")
	logger.Error(context.Background(), expectedErr)

	if len(mockLogger.debugIn) != 1 {
		t.Fatalf("Expected Debug to be called once, got %d calls", len(mockLogger.debugIn))
	}

	if len(mockLogger.infoIn) != 1 {
		t.Fatalf("Expected Info to be called once, got %d calls", len(mockLogger.infoIn))
	}

	if len(mockLogger.warnIn) != 1 {
		t.Fatalf("Expected Warn to be called once, got %d calls", len(mockLogger.warnIn))
	}

	if len(mockLogger.errorIn) != 1 {
		t.Fatalf("Expected Error to be called once, got %d calls", len(mockLogger.errorIn))
	}

	if len(mockLogger.debugIn[0].args) != 0 {
		t.Fatalf(
			"Expected Debug to be called with 0 arguments, got %d",
			len(mockLogger.debugIn[0].args),
		)
	}

	if len(mockLogger.infoIn[0].args) != 0 {
		t.Fatalf(
			"Expected Info to be called with 0 arguments, got %d",
			len(mockLogger.infoIn[0].args),
		)
	}

	if len(mockLogger.warnIn[0].args) != 0 {
		t.Fatalf(
			"Expected Warn to be called with 0 arguments, got %d",
			len(mockLogger.warnIn[0].args),
		)
	}

	if len(mockLogger.errorIn[0].args) != 0 {
		t.Fatalf(
			"Expected Error to be called with 0 arguments, got %d",
			len(mockLogger.errorIn[0].args),
		)
	}

	if mockLogger.debugIn[0].msg != "debug message" {
		t.Errorf(
			"Expected Debug message to be 'debug message', got '%s'",
			mockLogger.debugIn[0].msg,
		)
	}

	if mockLogger.infoIn[0].msg != "info message" {
		t.Errorf(
			"Expected Info message to be 'info message', got '%s'",
			mockLogger.infoIn[0].msg,
		)
	}

	if mockLogger.warnIn[0].msg != "warn message" {
		t.Errorf(
			"Expected Warn message to be 'warn message', got '%s'",
			mockLogger.warnIn[0].msg,
		)
	}

	if !errors.Is(mockLogger.errorIn[0].err, expectedErr) {
		t.Errorf(
			"Expected Error to be called with error 'error message', got '%v'",
			mockLogger.errorIn[0].err,
		)
	}
}

// TestContextLogger_ImmutableContext tests that the ContextLogger correctly handles
// immutable context chains, ensuring that each logging call receives the correct set of
// context values at each level of the context hierarchy.
func TestContextLogger_ImmutableContext(t *testing.T) {
	t.Parallel()

	mockLogger := new(mockLogger)

	logger, err := decorator.NewContextLogger(mockLogger)
	if err != nil {
		t.Fatalf("failed to create ContextLogger: %v", err)
	}

	ctx0 := context.Background()

	ctx1 := decorator.WithContextValue(ctx0, "key1", "value1")
	ctx2 := decorator.WithContextValue(ctx1, "key2", "value2")
	ctx3 := decorator.WithContextValue(ctx2, "key3", "value3")
	ctx4 := decorator.WithContextValue(ctx3, "key4", "value4")

	logger.Debug(ctx1, "")
	logger.Info(ctx2, "")
	logger.Warn(ctx3, "")
	logger.Error(ctx4, errors.New(""))

	if len(mockLogger.debugIn) != 1 {
		t.Fatalf("Expected Debug to be called once, got %d calls", len(mockLogger.debugIn))
	}

	if len(mockLogger.infoIn) != 1 {
		t.Fatalf("Expected Info to be called once, got %d calls", len(mockLogger.infoIn))
	}

	if len(mockLogger.warnIn) != 1 {
		t.Fatalf("Expected Warn to be called once, got %d calls", len(mockLogger.warnIn))
	}

	if len(mockLogger.errorIn) != 1 {
		t.Fatalf("Expected Error to be called once, got %d calls", len(mockLogger.errorIn))
	}

	if len(mockLogger.debugIn[0].args) != 1 {
		t.Fatalf(
			"Expected Debug to be called with 1 argument, got %d",
			len(mockLogger.debugIn[0].args),
		)
	}

	if len(mockLogger.infoIn[0].args) != 1 {
		t.Fatalf(
			"Expected Info to be called with 1 argument, got %d",
			len(mockLogger.infoIn[0].args),
		)
	}

	if len(mockLogger.warnIn[0].args) != 1 {
		t.Fatalf(
			"Expected Warn to be called with 1 argument, got %d",
			len(mockLogger.warnIn[0].args),
		)
	}

	if len(mockLogger.errorIn[0].args) != 1 {
		t.Fatalf(
			"Expected Error to be called with 1 argument, got %d",
			len(mockLogger.errorIn[0].args),
		)
	}

	ctxValue1, ok := mockLogger.debugIn[0].args[0].(map[string]any)
	if !ok {
		t.Fatalf(
			"Expected Debug argument to be a map[string]any, got '%T'",
			mockLogger.debugIn[0].args[0],
		)
	}

	if len(ctxValue1) != 1 || ctxValue1["key1"] != "value1" {
		t.Errorf(
			"Expected Debug context to be map[string]any{\"key1\": \"value1\"}, got '%v'",
			ctxValue1,
		)
	}

	ctxValue2, ok := mockLogger.infoIn[0].args[0].(map[string]any)
	if !ok {
		t.Fatalf(
			"Expected Info argument to be a map[string]any, got '%T'",
			mockLogger.infoIn[0].args[0],
		)
	}

	if len(ctxValue2) != 2 || ctxValue2["key1"] != "value1" || ctxValue2["key2"] != "value2" {
		t.Errorf(
			"Expected Info context to be map[string]any{\"key1\": \"value1\", \"key2\": \"value2\"}, got '%v'",
			ctxValue2,
		)
	}

	ctxValue3, ok := mockLogger.warnIn[0].args[0].(map[string]any)
	if !ok {
		t.Fatalf(
			"Expected Warn argument to be a map[string]any, got '%T'",
			mockLogger.warnIn[0].args[0],
		)
	}

	if len(ctxValue3) != 3 ||
		ctxValue3["key1"] != "value1" ||
		ctxValue3["key2"] != "value2" ||
		ctxValue3["key3"] != "value3" {
		t.Errorf(
			"Expected Warn context to be map[string]any{\"key1\": \"value1\", "+
				"\"key2\": \"value2\", \"key3\": \"value3\"}, got '%v'",
			ctxValue3,
		)
	}

	ctxValue4, ok := mockLogger.errorIn[0].args[0].(map[string]any)
	if !ok {
		t.Fatalf(
			"Expected Error argument to be a map[string]any, got '%T'",
			mockLogger.errorIn[0].args[0],
		)
	}

	if len(ctxValue4) != 4 ||
		ctxValue4["key1"] != "value1" ||
		ctxValue4["key2"] != "value2" ||
		ctxValue4["key3"] != "value3" ||
		ctxValue4["key4"] != "value4" {
		t.Errorf(
			"Expected Error context to be map[string]any{\"key1\": \"value1\", "+
				"\"key2\": \"value2\", \"key3\": \"value3\", \"key4\": \"value4\"}, got '%v'",
			ctxValue4,
		)
	}
}
