package decorator_test

import (
	"context"
	"errors"
	"testing"

	"github.com/yuppyweb/cakelog/decorator"
)

func TestContextLogger_Debug(t *testing.T) {
	t.Parallel()

	mockLogger := new(mockLogger)

	logger := decorator.NewContextLogger(mockLogger)

	ctx := logger.PutContext(context.Background(), "userID", 123)

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

	ctxValue, ok := mockLogger.debugIn[0].args[2].(map[any]any)
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

func TestContextLogger_Info(t *testing.T) {
	t.Parallel()

	mockLogger := new(mockLogger)

	logger := decorator.NewContextLogger(mockLogger)

	ctx := logger.PutContext(context.Background(), "requestID", "abc-123")

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

	ctxValue, ok := mockLogger.infoIn[0].args[2].(map[any]any)
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

func TestContextLogger_Warn(t *testing.T) {
	t.Parallel()

	mockLogger := new(mockLogger)

	logger := decorator.NewContextLogger(mockLogger)

	ctx := logger.PutContext(context.Background(), "sessionID", "xyz-789")

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

	ctxValue, ok := mockLogger.warnIn[0].args[2].(map[any]any)
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

func TestContextLogger_Error(t *testing.T) {
	t.Parallel()

	mockLogger := new(mockLogger)

	logger := decorator.NewContextLogger(mockLogger)

	ctx := logger.PutContext(context.Background(), "transactionID", "txn-456")

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

	ctxValue, ok := mockLogger.errorIn[0].args[2].(map[any]any)
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

func TestContextLogger_WithEmptyContext(t *testing.T) {
	t.Parallel()

	mockLogger := new(mockLogger)
	expectedErr := errors.New("error message")

	logger := decorator.NewContextLogger(mockLogger)

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

func TestContextLogger_ImmutableContext(t *testing.T) {
	t.Parallel()

	mockLogger := new(mockLogger)

	logger := decorator.NewContextLogger(mockLogger)

	ctx0 := context.Background()

	ctx1 := logger.PutContext(ctx0, "key1", "value1")
	ctx2 := logger.PutContext(ctx1, "key2", "value2")
	ctx3 := logger.PutContext(ctx2, "key3", "value3")
	ctx4 := logger.PutContext(ctx3, "key4", "value4")

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

	ctxValue1, ok := mockLogger.debugIn[0].args[0].(map[any]any)
	if !ok {
		t.Fatalf(
			"Expected Debug argument to be a map[string]any, got '%T'",
			mockLogger.debugIn[0].args[0],
		)
	}

	if len(ctxValue1) != 1 || ctxValue1["key1"] != "value1" {
		t.Errorf(
			"Expected Debug context to be map[any]any{\"key1\": \"value1\"}, got '%v'",
			ctxValue1,
		)
	}

	ctxValue2, ok := mockLogger.infoIn[0].args[0].(map[any]any)
	if !ok {
		t.Fatalf(
			"Expected Info argument to be a map[string]any, got '%T'",
			mockLogger.infoIn[0].args[0],
		)
	}

	if len(ctxValue2) != 2 || ctxValue2["key1"] != "value1" || ctxValue2["key2"] != "value2" {
		t.Errorf(
			"Expected Info context to be map[any]any{\"key1\": \"value1\", \"key2\": \"value2\"}, got '%v'",
			ctxValue2,
		)
	}

	ctxValue3, ok := mockLogger.warnIn[0].args[0].(map[any]any)
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
			"Expected Warn context to be map[any]any{\"key1\": \"value1\", "+
				"\"key2\": \"value2\", \"key3\": \"value3\"}, got '%v'",
			ctxValue3,
		)
	}

	ctxValue4, ok := mockLogger.errorIn[0].args[0].(map[any]any)
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
			"Expected Error context to be map[any]any{\"key1\": \"value1\", "+
				"\"key2\": \"value2\", \"key3\": \"value3\", \"key4\": \"value4\"}, got '%v'",
			ctxValue4,
		)
	}
}
