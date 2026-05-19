package adapter_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/yuppyweb/cakelog/adapter"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// mockZapCore is a custom zapcore.Core implementation used for testing the ZapLogger adapter.
// It captures log entries and fields for verification in tests.
type mockZapCore struct {
	entry  zapcore.Entry
	fields []zapcore.Field
}

// Enabled always returns true, indicating that all log levels are enabled for this mock core.
func (c *mockZapCore) Enabled(zapcore.Level) bool {
	return true
}

// With appends the provided fields to the mock core's fields and returns the core itself for chaining.
func (c *mockZapCore) With(fields []zapcore.Field) zapcore.Core {
	c.fields = append(c.fields, fields...)

	return c
}

// Check adds the log entry to the checked entry if the log level is enabled,
// allowing the ZapLogger to write log messages through this mock core.
func (c *mockZapCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	return ce.AddCore(ent, c)
}

// Write captures the log entry and fields when a log message is written through the ZapLogger,
// allowing tests to verify that the correct log level, message, and fields are being used.
func (c *mockZapCore) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	c.entry = ent
	c.fields = fields

	return nil
}

// Sync is a no-op for the mock core, as it does not perform any actual I/O operations.
func (c *mockZapCore) Sync() error {
	return nil
}

// TestZapLogger_DebugWithDefaultArgsKey verifies that ZapLogger correctly logs debug messages with default args key.
func TestZapLogger_DebugWithDefaultArgsKey(t *testing.T) {
	t.Parallel()

	mockCore := new(mockZapCore)
	values := []any{"debug", 42}

	logger, err := adapter.NewZapLogger(zap.New(mockCore))
	if err != nil {
		t.Fatalf("failed to create ZapLogger: %v", err)
	}

	logger.Debug(context.Background(), "debug message", values...)

	if mockCore.entry.Level != zap.DebugLevel {
		t.Errorf("unexpected log level: got %v, want %v", mockCore.entry.Level, zap.DebugLevel)
	}

	expectedFields := []zapcore.Field{
		zap.Any("context", values),
	}

	if len(mockCore.fields) != len(expectedFields) {
		t.Fatalf(
			"unexpected number of fields: got %d, want %d",
			len(mockCore.fields),
			len(expectedFields),
		)
	}

	for idx, field := range expectedFields {
		if field.Key != mockCore.fields[idx].Key {
			t.Errorf(
				"unexpected field key at index %d: got %q, want %q",
				idx,
				mockCore.fields[idx].Key,
				field.Key,
			)
		}
	}

	fieldValue, ok := mockCore.fields[0].Interface.([]any)
	if !ok {
		t.Errorf("unexpected field type: got %T, want []any", mockCore.fields[0].Interface)
	}

	for idx, val := range values {
		if fieldValue[idx] != val {
			t.Errorf(
				"unexpected field value at index %d: got %v, want %v",
				idx,
				fieldValue[idx],
				val,
			)
		}
	}
}

// TestZapLogger_InfoWithDefaultArgsKey verifies that ZapLogger correctly logs info messages with default args key.
func TestZapLogger_InfoWithDefaultArgsKey(t *testing.T) {
	t.Parallel()

	mockCore := new(mockZapCore)
	values := []any{"info", 65}

	logger, err := adapter.NewZapLogger(zap.New(mockCore))
	if err != nil {
		t.Fatalf("failed to create ZapLogger: %v", err)
	}

	logger.Info(context.Background(), "info message", values...)

	if mockCore.entry.Level != zap.InfoLevel {
		t.Errorf("unexpected log level: got %v, want %v", mockCore.entry.Level, zap.InfoLevel)
	}

	expectedFields := []zapcore.Field{
		zap.Any("context", values),
	}

	if len(mockCore.fields) != len(expectedFields) {
		t.Fatalf(
			"unexpected number of fields: got %d, want %d",
			len(mockCore.fields),
			len(expectedFields),
		)
	}

	for idx, field := range expectedFields {
		if field.Key != mockCore.fields[idx].Key {
			t.Errorf(
				"unexpected field key at index %d: got %q, want %q",
				idx,
				mockCore.fields[idx].Key,
				field.Key,
			)
		}
	}

	fieldValue, ok := mockCore.fields[0].Interface.([]any)
	if !ok {
		t.Errorf("unexpected field type: got %T, want []any", mockCore.fields[0].Interface)
	}

	for idx, val := range values {
		if fieldValue[idx] != val {
			t.Errorf(
				"unexpected field value at index %d: got %v, want %v",
				idx,
				fieldValue[idx],
				val,
			)
		}
	}
}

// TestZapLogger_WarnWithDefaultArgsKey verifies that ZapLogger correctly logs warn messages with default args key.
func TestZapLogger_WarnWithDefaultArgsKey(t *testing.T) {
	t.Parallel()

	mockCore := new(mockZapCore)
	values := []any{"warn", 99}

	logger, err := adapter.NewZapLogger(zap.New(mockCore))
	if err != nil {
		t.Fatalf("failed to create ZapLogger: %v", err)
	}

	logger.Warn(context.Background(), "warn message", values...)

	if mockCore.entry.Level != zap.WarnLevel {
		t.Errorf("unexpected log level: got %v, want %v", mockCore.entry.Level, zap.WarnLevel)
	}

	expectedFields := []zapcore.Field{
		zap.Any("context", values),
	}

	if len(mockCore.fields) != len(expectedFields) {
		t.Fatalf(
			"unexpected number of fields: got %d, want %d",
			len(mockCore.fields),
			len(expectedFields),
		)
	}

	for idx, field := range expectedFields {
		if field.Key != mockCore.fields[idx].Key {
			t.Errorf(
				"unexpected field key at index %d: got %q, want %q",
				idx,
				mockCore.fields[idx].Key,
				field.Key,
			)
		}
	}

	fieldValue, ok := mockCore.fields[0].Interface.([]any)
	if !ok {
		t.Errorf("unexpected field type: got %T, want []any", mockCore.fields[0].Interface)
	}

	for idx, val := range values {
		if fieldValue[idx] != val {
			t.Errorf(
				"unexpected field value at index %d: got %v, want %v",
				idx,
				fieldValue[idx],
				val,
			)
		}
	}
}

// TestZapLogger_ErrorWithDefaultArgsKey verifies that ZapLogger correctly logs error messages with default args key.
func TestZapLogger_ErrorWithDefaultArgsKey(t *testing.T) {
	t.Parallel()

	mockCore := new(mockZapCore)
	values := []any{"error", 123}
	expectedErr := errors.New("error message")

	logger, err := adapter.NewZapLogger(zap.New(mockCore))
	if err != nil {
		t.Fatalf("failed to create ZapLogger: %v", err)
	}

	logger.Error(context.Background(), expectedErr, values...)

	if mockCore.entry.Level != zap.ErrorLevel {
		t.Errorf("unexpected log level: got %v, want %v", mockCore.entry.Level, zap.ErrorLevel)
	}

	expectedFields := []zapcore.Field{
		zap.Any("context", values),
	}

	if len(mockCore.fields) != len(expectedFields) {
		t.Fatalf(
			"unexpected number of fields: got %d, want %d",
			len(mockCore.fields),
			len(expectedFields),
		)
	}

	for idx, field := range expectedFields {
		if field.Key != mockCore.fields[idx].Key {
			t.Errorf(
				"unexpected field key at index %d: got %q, want %q",
				idx,
				mockCore.fields[idx].Key,
				field.Key,
			)
		}
	}
}

// TestZapLogger_DebugWithCustomArgsKey verifies that ZapLogger correctly logs debug messages with custom args key.
func TestZapLogger_DebugWithCustomArgsKey(t *testing.T) {
	t.Parallel()

	mockCore := new(mockZapCore)
	values := []any{"debug", 42}
	customKey := "custom_debug"

	logger, err := adapter.NewZapLogger(zap.New(mockCore), adapter.WithArgsKey(customKey))
	if err != nil {
		t.Fatalf("failed to create ZapLogger: %v", err)
	}

	logger.Debug(context.Background(), "debug message", values...)

	if mockCore.entry.Level != zap.DebugLevel {
		t.Errorf("unexpected log level: got %v, want %v", mockCore.entry.Level, zap.DebugLevel)
	}

	expectedFields := []zapcore.Field{
		zap.Any(customKey, values),
	}

	if len(mockCore.fields) != len(expectedFields) {
		t.Fatalf(
			"unexpected number of fields: got %d, want %d",
			len(mockCore.fields),
			len(expectedFields),
		)
	}

	for idx, field := range expectedFields {
		if field.Key != mockCore.fields[idx].Key {
			t.Errorf(
				"unexpected field key at index %d: got %q, want %q",
				idx,
				mockCore.fields[idx].Key,
				field.Key,
			)
		}
	}

	fieldValue, ok := mockCore.fields[0].Interface.([]any)
	if !ok {
		t.Errorf("unexpected field type: got %T, want []any", mockCore.fields[0].Interface)
	}

	for idx, val := range values {
		if fieldValue[idx] != val {
			t.Errorf(
				"unexpected field value at index %d: got %v, want %v",
				idx,
				fieldValue[idx],
				val,
			)
		}
	}
}

// TestZapLogger_InfoWithCustomArgsKey verifies that ZapLogger correctly logs info messages with custom args key.
func TestZapLogger_InfoWithCustomArgsKey(t *testing.T) {
	t.Parallel()

	mockCore := new(mockZapCore)
	values := []any{"info", 65}
	customKey := "custom_info"

	logger, err := adapter.NewZapLogger(zap.New(mockCore), adapter.WithArgsKey(customKey))
	if err != nil {
		t.Fatalf("failed to create ZapLogger: %v", err)
	}

	logger.Info(context.Background(), "info message", values...)

	if mockCore.entry.Level != zap.InfoLevel {
		t.Errorf("unexpected log level: got %v, want %v", mockCore.entry.Level, zap.InfoLevel)
	}

	expectedFields := []zapcore.Field{
		zap.Any(customKey, values),
	}

	if len(mockCore.fields) != len(expectedFields) {
		t.Fatalf(
			"unexpected number of fields: got %d, want %d",
			len(mockCore.fields),
			len(expectedFields),
		)
	}

	for idx, field := range expectedFields {
		if field.Key != mockCore.fields[idx].Key {
			t.Errorf(
				"unexpected field key at index %d: got %q, want %q",
				idx,
				mockCore.fields[idx].Key,
				field.Key,
			)
		}
	}

	fieldValue, ok := mockCore.fields[0].Interface.([]any)
	if !ok {
		t.Errorf("unexpected field type: got %T, want []any", mockCore.fields[0].Interface)
	}

	for idx, val := range values {
		if fieldValue[idx] != val {
			t.Errorf(
				"unexpected field value at index %d: got %v, want %v",
				idx,
				fieldValue[idx],
				val,
			)
		}
	}
}

// TestZapLogger_WarnWithCustomArgsKey verifies that ZapLogger correctly logs warn messages with custom args key.
func TestZapLogger_WarnWithCustomArgsKey(t *testing.T) {
	t.Parallel()

	mockCore := new(mockZapCore)
	values := []any{"warn", 99}
	customKey := "custom_warn"

	logger, err := adapter.NewZapLogger(zap.New(mockCore), adapter.WithArgsKey(customKey))
	if err != nil {
		t.Fatalf("failed to create ZapLogger: %v", err)
	}

	logger.Warn(context.Background(), "warn message", values...)

	if mockCore.entry.Level != zap.WarnLevel {
		t.Errorf("unexpected log level: got %v, want %v", mockCore.entry.Level, zap.WarnLevel)
	}

	expectedFields := []zapcore.Field{
		zap.Any(customKey, values),
	}

	if len(mockCore.fields) != len(expectedFields) {
		t.Fatalf(
			"unexpected number of fields: got %d, want %d",
			len(mockCore.fields),
			len(expectedFields),
		)
	}

	for idx, field := range expectedFields {
		if field.Key != mockCore.fields[idx].Key {
			t.Errorf(
				"unexpected field key at index %d: got %q, want %q",
				idx,
				mockCore.fields[idx].Key,
				field.Key,
			)
		}
	}

	fieldValue, ok := mockCore.fields[0].Interface.([]any)
	if !ok {
		t.Errorf("unexpected field type: got %T, want []any", mockCore.fields[0].Interface)
	}

	for idx, val := range values {
		if fieldValue[idx] != val {
			t.Errorf(
				"unexpected field value at index %d: got %v, want %v",
				idx,
				fieldValue[idx],
				val,
			)
		}
	}
}

// TestZapLogger_ErrorWithCustomArgsKey verifies that ZapLogger correctly logs error messages with custom args key.
func TestZapLogger_ErrorWithCustomArgsKey(t *testing.T) {
	t.Parallel()

	mockCore := new(mockZapCore)
	values := []any{"error", 123}
	customKey := "custom_error"
	expectedErr := errors.New("error message")

	logger, err := adapter.NewZapLogger(zap.New(mockCore), adapter.WithArgsKey(customKey))
	if err != nil {
		t.Fatalf("failed to create ZapLogger: %v", err)
	}

	logger.Error(context.Background(), expectedErr, values...)

	if mockCore.entry.Level != zap.ErrorLevel {
		t.Errorf("unexpected log level: got %v, want %v", mockCore.entry.Level, zap.ErrorLevel)
	}

	expectedFields := []zapcore.Field{
		zap.Any(customKey, values),
	}

	if len(mockCore.fields) != len(expectedFields) {
		t.Fatalf(
			"unexpected number of fields: got %d, want %d",
			len(mockCore.fields),
			len(expectedFields),
		)
	}

	for idx, field := range expectedFields {
		if field.Key != mockCore.fields[idx].Key {
			t.Errorf(
				"unexpected field key at index %d: got %q, want %q",
				idx,
				mockCore.fields[idx].Key,
				field.Key,
			)
		}
	}

	fieldValue, ok := mockCore.fields[0].Interface.([]any)
	if !ok {
		t.Errorf("unexpected field type: got %T, want []any", mockCore.fields[0].Interface)
	}

	for idx, val := range values {
		if fieldValue[idx] != val {
			t.Errorf(
				"unexpected field value at index %d: got %v, want %v",
				idx,
				fieldValue[idx],
				val,
			)
		}
	}
}

// TestNewZapLogger_WithNilLogger verifies that NewZapLogger returns an error when provided with a nil logger.
func TestNewZapLogger_WithNilLogger(t *testing.T) {
	t.Parallel()

	_, err := adapter.NewZapLogger(nil)
	if err == nil {
		t.Fatal("expected error when creating ZapLogger with nil logger, but got nil")
	}

	if !errors.Is(err, adapter.ErrNilZapLogger) {
		t.Errorf(
			"unexpected error when creating ZapLogger with nil logger:\nGot:  %v\nWant: %v",
			err,
			adapter.ErrNilZapLogger,
		)
	}
}

// TestNewZapLogger_WithNilOption verifies that NewZapLogger returns an error when provided with a nil option.
func TestNewZapLogger_WithNilOption(t *testing.T) {
	t.Parallel()

	_, err := adapter.NewZapLogger(zap.NewNop(), nil)
	if err == nil {
		t.Fatal("expected error when creating ZapLogger with nil option, but got nil")
	}

	if !errors.Is(err, adapter.ErrNilZapOption) {
		t.Errorf(
			"unexpected error when creating ZapLogger with nil option:\nGot:  %v\nWant: %v",
			err,
			adapter.ErrNilZapOption,
		)
	}
}

// TestNewZapLogger_WithInvalidArgsKey verifies that NewZapLogger returns an error when provided with an empty args key.
func TestNewZapLogger_WithInvalidArgsKey(t *testing.T) {
	t.Parallel()

	_, err := adapter.NewZapLogger(zap.NewNop(), adapter.WithArgsKey(""))
	if err == nil {
		t.Fatal("expected error when creating ZapLogger with empty args key, but got nil")
	}

	if !errors.Is(err, adapter.ErrEmptyArgsKey) {
		t.Errorf(
			"unexpected error when creating ZapLogger with empty args key:\nGot:  %v\nWant: %v",
			err,
			adapter.ErrEmptyArgsKey,
		)
	}

	if !strings.Contains(err.Error(), "failed to apply option:") {
		t.Errorf(
			"error message does not contain expected text:\nGot:  %s\nWant to contain: %s",
			err.Error(),
			"failed to apply option:",
		)
	}
}
