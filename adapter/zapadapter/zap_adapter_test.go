package zapadapter_test

import (
	"context"
	"errors"
	"testing"

	"github.com/yuppyweb/cakelog"
	"github.com/yuppyweb/cakelog/adapter/zapadapter"
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

// TestZapLogger_Debug verifies that ZapLogger correctly logs debug messages.
func TestZapLogger_Debug(t *testing.T) {
	t.Parallel()

	mockCore := new(mockZapCore)

	logger, err := zapadapter.New(zap.New(mockCore))
	if err != nil {
		t.Fatalf("failed to create ZapLogger: %v", err)
	}

	logger.Debug(context.Background(), "debug message", "debug", 42)

	if mockCore.entry.Level != zap.DebugLevel {
		t.Errorf("unexpected log level: got %v, want %v", mockCore.entry.Level, zap.DebugLevel)
	}

	if len(mockCore.fields) != 1 {
		t.Fatalf("unexpected number of fields: got %d, want 1", len(mockCore.fields))
	}

	assertZapField(t, mockCore.fields, "debug", 42)
}

// TestZapLogger_Info verifies that ZapLogger correctly logs info messages.
func TestZapLogger_Info(t *testing.T) {
	t.Parallel()

	mockCore := new(mockZapCore)

	logger, err := zapadapter.New(zap.New(mockCore))
	if err != nil {
		t.Fatalf("failed to create ZapLogger: %v", err)
	}

	logger.Info(context.Background(), "info message", "info", 65)

	if mockCore.entry.Level != zap.InfoLevel {
		t.Errorf("unexpected log level: got %v, want %v", mockCore.entry.Level, zap.InfoLevel)
	}

	if len(mockCore.fields) != 1 {
		t.Fatalf("unexpected number of fields: got %d, want 1", len(mockCore.fields))
	}

	assertZapField(t, mockCore.fields, "info", 65)
}

// TestZapLogger_Warn verifies that ZapLogger correctly logs warn messages.
func TestZapLogger_Warn(t *testing.T) {
	t.Parallel()

	mockCore := new(mockZapCore)

	logger, err := zapadapter.New(zap.New(mockCore))
	if err != nil {
		t.Fatalf("failed to create ZapLogger: %v", err)
	}

	logger.Warn(context.Background(), "warn message", "warn", 99)

	if mockCore.entry.Level != zap.WarnLevel {
		t.Errorf("unexpected log level: got %v, want %v", mockCore.entry.Level, zap.WarnLevel)
	}

	if len(mockCore.fields) != 1 {
		t.Fatalf("unexpected number of fields: got %d, want 1", len(mockCore.fields))
	}

	assertZapField(t, mockCore.fields, "warn", 99)
}

// TestZapLogger_Error verifies that ZapLogger correctly logs error messages.
func TestZapLogger_Error(t *testing.T) {
	t.Parallel()

	mockCore := new(mockZapCore)
	expectedErr := errors.New("error message")

	logger, err := zapadapter.New(zap.New(mockCore))
	if err != nil {
		t.Fatalf("failed to create ZapLogger: %v", err)
	}

	logger.Error(context.Background(), expectedErr, "error", 123)

	if mockCore.entry.Level != zap.ErrorLevel {
		t.Errorf("unexpected log level: got %v, want %v", mockCore.entry.Level, zap.ErrorLevel)
	}

	if len(mockCore.fields) != 1 {
		t.Fatalf("unexpected number of fields: got %d, want 1", len(mockCore.fields))
	}

	assertZapField(t, mockCore.fields, "error", 123)
}

// TestZapLogger_ErrorNil verifies that ZapLogger logs a nil error with an empty message.
func TestZapLogger_ErrorNil(t *testing.T) {
	t.Parallel()

	mockCore := new(mockZapCore)

	logger, err := zapadapter.New(zap.New(mockCore))
	if err != nil {
		t.Fatalf("failed to create ZapLogger: %v", err)
	}

	logger.Error(context.Background(), nil, "error", 123)

	if mockCore.entry.Level != zap.ErrorLevel {
		t.Errorf("unexpected log level: got %v, want %v", mockCore.entry.Level, zap.ErrorLevel)
	}

	if mockCore.entry.Message != "" {
		t.Errorf("unexpected log message: got %q, want empty", mockCore.entry.Message)
	}

	if len(mockCore.fields) != 1 {
		t.Fatalf("unexpected number of fields: got %d, want 1", len(mockCore.fields))
	}

	assertZapField(t, mockCore.fields, "error", 123)
}

// TestZapLogger_FieldsMap verifies that ZapLogger expands a map into individual fields.
func TestZapLogger_FieldsMap(t *testing.T) {
	t.Parallel()

	mockCore := new(mockZapCore)

	logger, err := zapadapter.New(zap.New(mockCore))
	if err != nil {
		t.Fatalf("failed to create ZapLogger: %v", err)
	}

	logger.Info(context.Background(), "mapped", map[string]any{"method": "GET", "status": 200})

	if len(mockCore.fields) != 2 {
		t.Fatalf("unexpected number of fields: got %d, want 2", len(mockCore.fields))
	}

	assertZapField(t, mockCore.fields, "method", "GET")
	assertZapField(t, mockCore.fields, "status", 200)
}

// TestZapLogger_NoArgs verifies that ZapLogger logs messages without extra arguments.
func TestZapLogger_NoArgs(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name  string
		call  func(cakelog.Logger)
		level zapcore.Level
		msg   string
	}{
		{
			name: "debug",
			call: func(logger cakelog.Logger) {
				logger.Debug(context.Background(), "debug message")
			},
			level: zap.DebugLevel,
			msg:   "debug message",
		},
		{
			name: "info",
			call: func(logger cakelog.Logger) {
				logger.Info(context.Background(), "info message")
			},
			level: zap.InfoLevel,
			msg:   "info message",
		},
		{
			name: "warn",
			call: func(logger cakelog.Logger) {
				logger.Warn(context.Background(), "warn message")
			},
			level: zap.WarnLevel,
			msg:   "warn message",
		},
		{
			name: "error",
			call: func(logger cakelog.Logger) {
				logger.Error(context.Background(), errors.New("error message"))
			},
			level: zap.ErrorLevel,
			msg:   "error message",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockCore := new(mockZapCore)

			logger, err := zapadapter.New(zap.New(mockCore))
			if err != nil {
				t.Fatalf("failed to create ZapLogger: %v", err)
			}

			tc.call(logger)

			if mockCore.entry.Level != tc.level {
				t.Errorf("unexpected log level: got %v, want %v", mockCore.entry.Level, tc.level)
			}

			if mockCore.entry.Message != tc.msg {
				t.Errorf("unexpected log message: got %q, want %q", mockCore.entry.Message, tc.msg)
			}

			if len(mockCore.fields) != 0 {
				t.Errorf("expected no fields, got %v", mockCore.fields)
			}
		})
	}
}

// TestNewZapLogger_WithNilLogger verifies that NewZapLogger returns an error when provided with a nil logger.
func TestNewZapLogger_WithNilLogger(t *testing.T) {
	t.Parallel()

	_, err := zapadapter.New(nil)
	if err == nil {
		t.Fatal("expected error when creating ZapLogger with nil logger, but got nil")
	}

	if !errors.Is(err, zapadapter.ErrNilZapLogger) {
		t.Errorf(
			"unexpected error when creating ZapLogger with nil logger:\nGot:  %v\nWant: %v",
			err,
			zapadapter.ErrNilZapLogger,
		)
	}
}

func assertZapField(t *testing.T, fields []zapcore.Field, key string, want any) {
	t.Helper()

	wantField := zap.Any(key, want)

	for _, field := range fields {
		if field.Key != key {
			continue
		}

		if !field.Equals(wantField) {
			t.Errorf("unexpected field %q: got %#v, want %#v", key, field, wantField)
		}

		return
	}

	t.Errorf("expected field %q, got %v", key, fields)
}
